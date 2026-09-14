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

func (a *Agent) Execute(_ context.Context, execs agentexecution.AgentExecutions, batchExecution agents.BatchExecutionHint) (*agents.AgentResponse, error) {
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

	if a.cfgs[ndx].promptDefinition.HasSchemaOutput() {
		return agents.AgentExecutionResponse{
			CustomId:    br.CustomID,
			Status:      string(br.Status),
			Content:     []byte(br.Text),
			Ct:          agents.AgentResponseJSON,
			XMLSections: nil,
		}, nil
	}

	xs, err := a.cfgs[ndx].promptDefinition.XMLSections.ExtractFromText(br.Text)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return agents.AgentExecutionResponse{CustomId: br.CustomID, Status: string(br.Status)}, err
	}

	return agents.AgentExecutionResponse{
		CustomId:    br.CustomID,
		Status:      string(br.Status),
		Content:     []byte(br.Text),
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

	var err error

	var reqs []client.BatchRequest
	for _, item := range a.cfgs {

		systemPrompt, err := item.promptDefinition.SystemPrompt()
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return "", nil, err
		}

		userPrompt, err := item.promptDefinition.UserPrompt(item.PromptVars, true, true)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return "", nil, err
		}

		reqs = append(reqs, client.BatchRequest{
			CustomID: item.CustomID,
			Options: []client.Option{
				client.WithModel(DefaultConfig.Model),
				client.WithMaxTokens(DefaultConfig.MaxTokens),
				client.WithTemperature(DefaultConfig.Temperature),
				client.WithSystem(string(systemPrompt)),
				client.WithUserText(string(userPrompt)),
			},
		})
	}

	batch, err := cli.SubmitBatch(context.Background(), reqs...)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return "", nil, err
	}

	return batch.ID, nil, nil
}
