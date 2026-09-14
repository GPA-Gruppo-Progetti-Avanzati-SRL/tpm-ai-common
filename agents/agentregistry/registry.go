package agentregistry

import (
	"fmt"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/agents"
)

const semLogContextBasePromptRegistry = "agent-registry::"

var theRegistry = map[string]agents.AgentFactory{}

func AddAgent(nm string, a agents.AgentFactory) error {
	const semLogContext = semLogContextBasePromptRegistry + "add-agent"

	theRegistry[nm] = a
	return nil
}

func GetAgent(domain, site, n string) (agents.Agent, error) {
	const semLogContext = semLogContextBasePromptRegistry + "get-agent"

	af, ok := theRegistry[n]
	if !ok {
		return nil, fmt.Errorf("no agent named %q found", n)
	}

	return af(domain, site), nil
}
