package agentregistry

import (
	"fmt"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/anthropiclks/client"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/store/agentexecution"
)

const semLogContextBasePromptRegistry = "agent-registry::"

type Agent interface {
	Name() string
	Execute(execs agentexecution.AgentExecutions) (string, error)
	OnBatchResult(br *client.BatchResult, agentExec *agentexecution.AgentExecution) error
}

type AgentFactory func() Agent

var theRegistry = map[string]AgentFactory{}

func AddAgent(nm string, a AgentFactory) error {
	const semLogContext = semLogContextBasePromptRegistry + "add-agent"

	theRegistry[nm] = a
	return nil
}

func GetAgent(n string) (Agent, error) {
	const semLogContext = semLogContextBasePromptRegistry + "get-agent"

	af, ok := theRegistry[n]
	if !ok {
		return nil, fmt.Errorf("no agent named %q found", n)
	}

	return af(), nil
}
