package agents

import (
	"context"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/anthropiclks/client"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/store/agentexecution"
)

const (
	semLogPackageContext = "agents::"

	AgentResponseRAW  = "raw"
	AgentResponseJSON = "json"
	AgentResponseXML  = "xml"
)

type AgentExecutionResponse struct {
	CustomId    string        `json:"custom-id" yaml:"custom-id" mapstructure:"custom-id"`
	Content     []byte        `json:"content" yaml:"content" mapstructure:"content"`
	Typ         string        `json:"type" yaml:"type" mapstructure:"type"`
	XMLSections []XMLSections `json:"sections" yaml:"sections" mapstructure:"sections"`
}

type AgentResponse struct {
	Content []AgentExecutionResponse `json:"responses" yaml:"responses" mapstructure:"responses"`
	batchId string
}

type Agent interface {
	Name() string
	Execute(ctx context.Context, execs agentexecution.AgentExecutions) ([]AgentResponse, string, error)
	OnBatchResult(br *client.BatchResult, agentExec *agentexecution.AgentExecution) ([]AgentResponse, error)
}

type AgentFactory func(domain, site string) Agent
