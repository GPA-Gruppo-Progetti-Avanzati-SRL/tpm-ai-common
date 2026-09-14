package promptagent

import (
	"context"
	"errors"
	"time"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/agents"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/agents/prompt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/anthropiclks"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/anthropiclks/client"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/store/agentexecution"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/rs/zerolog/log"
)

const (
	Name                 = "prompt-agent"
	semLogPackageContext = "prompt-agent::"
)

type Agent struct {
	cfgs AgentParams
}

func NewAgentFactory(domain, site string) agents.Agent {
	return &Agent{}
}

func (a *Agent) Name() string {
	return Name
}

func (a *Agent) Execute(_ context.Context, execs agentexecution.AgentExecutions, batchExecution ...agents.BatchExecutionHint) (*agents.AgentResponse, error) {
	const semLogContext = semLogPackageContext + "execute"
	var err error

	if len(execs) == 0 {
		err = errors.New(semLogContext + " no agent executions provided")
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	for _, exec := range execs {
		cfg := NewConfig(exec.Params)
		if err = cfg.Validate(); err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}

		cfg.CustomID = exec.CustomID
		cfg.promptDefinition, err = prompt.ReadPromptDefinition(cfg.PromptDefinitionFile)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}

		a.cfgs = append(a.cfgs, *cfg)
	}

	cli, err := anthropiclks.NewClientNG()
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	// No BatchExecutionHint means run synchronously (online): one immediate call
	// per execution. A hint selects the asynchronous Batch API path.
	if len(batchExecution) == 0 {
		return a.executeOnline(cli)
	}

	return a.executeBatch(cli, execs, batchExecution[0])
}

// executeOnline runs one synchronous LLM call per configured execution. A single
// execution that fails (transport error, or a max_tokens truncation surfaced as
// an error by client.Execute) is recorded with Status "error" / Ct "none" and
// does not abort the run; the overall AgentResponse.Status becomes "error" if
// any execution failed.
func (a *Agent) executeOnline(cli *client.Client) (*agents.AgentResponse, error) {
	const semLogContext = semLogPackageContext + "execute-online"

	agentResponse := &agents.AgentResponse{Status: "done"}
	for i := range a.cfgs {
		cfg := &a.cfgs[i]

		opts, err := a.buildOptions(cfg)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}

		out, err := cli.Execute(context.Background(), opts...)
		if err != nil {
			log.Error().Err(err).Str("custom-id", cfg.CustomID).Msg(semLogContext + " online execution failed")
			agentResponse.Responses = append(agentResponse.Responses, agents.AgentExecutionResponse{
				CustomId: cfg.CustomID,
				Status:   "error",
				Ct:       agents.AgentResponseNone,
			})
			agentResponse.Status = "error"
			continue
		}

		resp, err := a.buildExecutionResponse(cfg, cfg.CustomID, "succeeded", out.Text)
		if err != nil {
			log.Error().Err(err).Str("custom-id", cfg.CustomID).Msg(semLogContext)
			agentResponse.Responses = append(agentResponse.Responses, agents.AgentExecutionResponse{
				CustomId: cfg.CustomID,
				Status:   "error",
				Ct:       agents.AgentResponseNone,
			})
			agentResponse.Status = "error"
			continue
		}

		agentResponse.Responses = append(agentResponse.Responses, resp)
	}

	return agentResponse, nil
}

// executeBatch runs the executions through the asynchronous Batch API: submit (or
// resume via hint.BatchId), poll until the batch ends, then collect and map each
// result.
func (a *Agent) executeBatch(cli *client.Client, execs agentexecution.AgentExecutions, batchExecution agents.BatchExecutionHint) (*agents.AgentResponse, error) {
	const semLogContext = semLogPackageContext + "execute-batch"
	var err error

	var batchId string
	if batchExecution.BatchId == "" {
		batchId, _, err = a.submitBatch(cli, execs)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}
	} else {
		batchId = batchExecution.BatchId
	}

	batchReady, err := a.pollBatch(cli, batchId, batchExecution)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	agentResponse := &agents.AgentResponse{BatchId: batchId}
	if batchReady {
		log.Info().Msg(semLogContext + " collecting batch results")
		brs, err := cli.CollectBatchResults(context.Background(), batchId)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}

		agentResponse.Status = "done"
		for _, br := range brs {
			ndx := execs.FirstByCustomId(br.CustomID)
			if ndx == -1 {
				log.Error().Str("custom-id", br.CustomID).Msg(semLogContext + " custom id not matched")
			} else {
				resp, err := a.OnBatchResult(&br, &execs[ndx])
				if err != nil {
					log.Error().Err(err).Msg(semLogContext)
					return nil, err
				}

				agentResponse.Responses = append(agentResponse.Responses, resp)
			}
		}
	} else {
		err = errors.New(" batch results not ready")
		log.Error().Err(err).Msg(semLogContext)
		agentResponse.Status = "error"

	}
	return agentResponse, err
}

func (a *Agent) OnBatchResult(br *client.BatchResult, item *agentexecution.AgentExecution) (agents.AgentExecutionResponse, error) {
	const semLogContext = semLogPackageContext + "on-batch-result"
	var err error

	ndx := a.cfgs.FirstByCustomId(br.CustomID)
	if ndx == -1 {
		log.Error().Str("custom-id", br.CustomID).Msg(semLogContext + " custom id not matched")
		return agents.AgentExecutionResponse{CustomId: br.CustomID}, err
	}

	if br.Status != client.BatchSucceeded {
		err = errors.New(" result not succeeded")
		log.Error().Err(err).Str("custom-id", br.CustomID).Msg(semLogContext)
		return agents.AgentExecutionResponse{CustomId: br.CustomID, Status: string(br.Status), Ct: agents.AgentResponseNone}, err
	}

	return a.buildExecutionResponse(&a.cfgs[ndx], br.CustomID, string(br.Status), br.Text)
}

// buildExecutionResponse maps a successful model text output onto an
// AgentExecutionResponse using the prompt definition's declared output shape: a
// schema-output prompt yields JSON content, otherwise the XML-delimited sections
// are extracted from the text. Shared by the batch and online paths.
func (a *Agent) buildExecutionResponse(cfg *Config, customID, status, text string) (agents.AgentExecutionResponse, error) {
	const semLogContext = semLogPackageContext + "build-execution-response"

	if cfg.promptDefinition.HasSchemaOutput() {
		return agents.AgentExecutionResponse{
			CustomId:    customID,
			Status:      status,
			Content:     []byte(text),
			Ct:          agents.AgentResponseJSON,
			XMLSections: nil,
		}, nil
	}

	xs, err := cfg.promptDefinition.XMLSections.ExtractFromText(text)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return agents.AgentExecutionResponse{CustomId: customID, Status: status}, err
	}

	return agents.AgentExecutionResponse{
		CustomId:    customID,
		Status:      status,
		Content:     []byte(text),
		Ct:          agents.AgentResponseXML,
		XMLSections: xs,
	}, nil
}

func (a *Agent) pollBatch(cli *client.Client, batchId string, batchExecutionParams agents.BatchExecutionHint) (bool, error) {

	const semLogContext = semLogPackageContext + "poll-batch"
	log.Info().Dur("interval", batchExecutionParams.BatchPollInterval).Int("max-iterations", batchExecutionParams.MaxIterations).Msg(semLogContext)

	var numIterations int
	for {
		b, err := cli.GetBatch(context.Background(), batchId)
		if err != nil {
			return false, err
		}

		log.Info().Msgf("%s status=%s processing=%d", semLogContext, b.ProcessingStatus, b.RequestCounts.Processing)
		if b.ProcessingStatus == anthropic.MessageBatchProcessingStatusEnded {
			return true, nil
		}

		numIterations++
		if batchExecutionParams.MaxIterations > 0 && numIterations >= batchExecutionParams.MaxIterations {
			err = errors.New(" max iterations exceeded")
			log.Error().Err(err).Msg(semLogContext)
			return false, err
		}

		time.Sleep(batchExecutionParams.BatchPollInterval)
	}

}

func (a *Agent) submitBatch(cli *client.Client, items agentexecution.AgentExecutions) (string, []agents.AgentResponse, error) {

	const semLogContext = semLogPackageContext + "submit-batch"

	var reqs []client.BatchRequest
	for i := range a.cfgs {
		cfg := &a.cfgs[i]

		opts, err := a.buildOptions(cfg)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return "", nil, err
		}

		reqs = append(reqs, client.BatchRequest{
			CustomID: cfg.CustomID,
			Options:  opts,
		})
	}

	batch, err := cli.SubmitBatch(context.Background(), reqs...)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return "", nil, err
	}

	return batch.ID, nil, nil
}

// buildOptions assembles the client call options for a single execution from its
// config and prompt definition (system + rendered user prompt). Shared by the
// batch and online paths so both send identical request parameters.
func (a *Agent) buildOptions(cfg *Config) ([]client.Option, error) {
	const semLogContext = semLogPackageContext + "build-options"

	systemPrompt, err := cfg.promptDefinition.SystemPrompt()
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	userPrompt, err := cfg.promptDefinition.UserPrompt(cfg.PromptVars, true, true)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	return []client.Option{
		client.WithModel(cfg.Model),
		client.WithMaxTokens(cfg.MaxTokens),
		client.WithTemperature(cfg.Temperature),
		client.WithSystem(string(systemPrompt)),
		client.WithUserText(string(userPrompt)),
	}, nil
}
