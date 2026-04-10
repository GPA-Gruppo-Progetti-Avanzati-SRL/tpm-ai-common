package anthropiclks

import (
	"time"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/prompts"
)

type Request struct {
	Vars map[string]prompts.Variable
}

func (r Request) TextVariables() map[string]string {
	var m map[string]string
	for k, v := range r.Vars {
		switch v.Ct {
		case prompts.TextVariable:
			if len(m) == 0 {
				m = make(map[string]string)
			}
			m[k] = string(v.Value)
		}
	}

	return m
}

type RequestParam func(r *Request)

func WithPromptInput(vt prompts.VariableType, n string, ndx int, v []byte) RequestParam {
	return func(o *Request) {
		if o.Vars == nil {
			o.Vars = make(map[string]prompts.Variable)
		}

		if ndx < 0 {
			ndx = 0
		}
		o.Vars[n] = prompts.Variable{Name: n, Index: ndx, Ct: vt, Value: v}
	}
}

type Response struct {
	Content map[string]prompts.MessagePart
}

// BatchRequest is a single named request within a batch submission.
type BatchRequest struct {
	CustomID string
	Params   []RequestParam
}

// BatchRequestCounts mirrors the per-state counters returned by the Batches API.
type BatchRequestCounts struct {
	Processing int64 `yaml:"processing,omitempty" mapstructure:"processing,omitempty" json:"processing,omitempty"`
	Succeeded  int64 `yaml:"succeeded,omitempty" mapstructure:"succeeded,omitempty" json:"succeeded,omitempty"`
	Errored    int64 `yaml:"errored,omitempty" mapstructure:"errored,omitempty" json:"errored,omitempty"`
	Canceled   int64 `yaml:"canceled,omitempty" mapstructure:"canceled,omitempty" json:"canceled,omitempty"`
	Expired    int64 `yaml:"expired,omitempty" mapstructure:"expired,omitempty" json:"expired,omitempty"`
}

const (
	BatchStatusInProgress = "in_progress"
	BatchStatusCanceling  = "canceling"
	BatchStatusEnded      = "ended"
)

// BatchStatus holds the lifecycle state of a submitted batch.
type BatchStatus struct {
	ID               string             `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	ProcessingStatus string             `yaml:"processing_status,omitempty" mapstructure:"processing_status,omitempty" json:"processing_status,omitempty"` // "in_progress" | "canceling" | "ended"
	RequestCounts    BatchRequestCounts `yaml:"request_counts,omitempty" mapstructure:"request_counts,omitempty" json:"request_counts,omitempty"`
	ExpiresAt        time.Time          `yaml:"expires_at,omitempty" mapstructure:"expires_at,omitempty" json:"expires_at,omitempty"`
}

// BatchResult is the outcome for a single request inside a completed batch.
type BatchResult struct {
	CustomID string
	Response *Response
	Err      error // non-nil for errored / canceled / expired items
}
