package agents

import (
	"context"
	"time"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/anthropiclks/client"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/store/agentexecution"
)

const (
	semLogPackageContext = "agents::"

	AgentResponseRAW  = "raw"
	AgentResponseJSON = "json"
	AgentResponseXML  = "xml"
	AgentResponseNone = "none"
)

type AgentExecutionResponse struct {
	CustomId    string      `json:"custom-id" yaml:"custom-id" mapstructure:"custom-id"`
	Status      string      `json:"status" yaml:"status" mapstructure:"status"`
	Content     []byte      `yaml:"-" mapstructure:"-" json:"-"`
	Ct          string      `json:"ct" yaml:"ct" mapstructure:"ct"`
	XMLSections XMLSections `json:"xml_sections" yaml:"xml_sections" mapstructure:"xml_sections"`
}

type AgentResponse struct {
	Responses []AgentExecutionResponse `json:"responses,omitempty" yaml:"responses,omitempty" mapstructure:"responses,omitempty"`
	BatchId   string                   `json:"batch-id,omitempty" yaml:"batch-id,omitempty" mapstructure:"batch-id,omitempty"`
	Status    string                   `json:"status,omitempty" yaml:"status,omitempty" mapstructure:"status,omitempty"`
}

type BatchExecutionHint struct {
	BatchId           string
	BatchPollInterval time.Duration
	MaxIterations     int
}

type Agent interface {
	Name() string
	Execute(ctx context.Context, execs agentexecution.AgentExecutions, batchExecutionHint ...BatchExecutionHint) (*AgentResponse, error)
	OnBatchResult(br *client.BatchResult, agentExec *agentexecution.AgentExecution) (AgentExecutionResponse, error)
}

type AgentFactory func(domain, site string) Agent
