package promptagent

import (
	"context"
	"errors"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/agents"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/agents/prompt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/anthropiclks"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/anthropiclks/client"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/store/agentexecution"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/rs/zerolog/log"
)

const (
	Name                 = "prompt-agent"
	semLogPackageContext = "prompt-agent::"
)

type Agent struct {
	cfgs []*Config
}

func NewAgentFactory(domain, site string) agents.Agent {
	return &Agent{}
}

func (a *Agent) Name() string {
	return Name
}

func (a *Agent) Execute(_ context.Context, execs agentexecution.AgentExecutions) ([]agents.AgentResponse, string, error) {
	const semLogContext = semLogPackageContext + "execute"
	var err error

	if len(execs) == 0 {
		err = errors.New(semLogContext + " no agent executions provided")
		log.Error().Err(err).Msg(semLogContext)
		return nil, "", err
	}

	for _, exec := range execs {
		cfg := NewConfig(exec.Params)
		if err = cfg.Validate(); err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, "", err
		}

		cfg.promptDefinition, err = prompt.ReadPromptDefinition(cfg.PromptDefinitionFile)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, "", err
		}

		a.cfgs = append(a.cfgs, cfg)
	}

	batchId, _, err := a.submitBatch(execs)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, "", err
	}

	return nil, batchId, nil
}

func (a *Agent) OnBatchResult(br *client.BatchResult, item *agentexecution.AgentExecution) ([]agents.AgentResponse, error) {
	const semLogContext = semLogPackageContext + "on-batch-result"
	var err error

	/*
		xmlSections := a.promptDefinition.XMLSections
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}

		err = xmlSections.ExtractFromText(br.Text)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}

		for _, sect := range xmlSections.Sections {

		}
	*/

	return nil, err
}

func (a *Agent) submitBatch(items agentexecution.AgentExecutions) (string, []agents.AgentResponse, error) {

	const semLogContext = semLogPackageContext + "submit-batch"

	var err error

	_, err = anthropiclks.NewClientNG()
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return "", nil, err
	}

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
			CustomID: util.NewUUID(),
			Options: []client.Option{
				client.WithModel(DefaultConfig.Model),
				client.WithMaxTokens(DefaultConfig.MaxTokens),
				client.WithTemperature(DefaultConfig.Temperature),
				client.WithSystem(string(systemPrompt)),
				client.WithUserText(string(userPrompt)),
			},
		})
	}

	/*
		batch, err := cli.SubmitBatch(context.Background(), reqs...)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return "", nil, err
		}

		return batch.ID, nil, nil
	*/

	return "", nil, nil

}
