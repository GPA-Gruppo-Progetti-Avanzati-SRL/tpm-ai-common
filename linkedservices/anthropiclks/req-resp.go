package anthropiclks

import (
	"fmt"
	"strings"
	"time"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/prompts"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/rs/zerolog/log"
)

type Request struct {
	Vars        map[string]prompts.Variable
	MaxTokens   int64                  `yaml:"max-tokens,omitempty" mapstructure:"max-tokens,omitempty" json:"max-tokens,omitempty"`
	Temperature float64                `yaml:"temperature,omitempty" mapstructure:"temperature,omitempty" json:"temperature,omitempty"`
	Prompt      prompts.PromptTemplate `yaml:"prompt,omitempty" mapstructure:"prompt,omitempty" json:"prompt,omitempty"`
	Model       anthropic.Model        `yaml:"model,omitempty" mapstructure:"model,omitempty" json:"model,omitempty"`
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

func WithMaxTokens(n int64) RequestParam {
	return func(o *Request) {
		if n > 0 {
			o.MaxTokens = n
		}
	}
}

func WithTemperature(f float64) RequestParam {
	return func(o *Request) {
		if f < 0 {
			f = 0
		}

		o.Temperature = f
	}
}

func WithPrompt(p prompts.PromptTemplate) RequestParam {
	return func(o *Request) {
		o.Prompt = p
	}
}

func WithPromptName(n string) RequestParam {
	return func(o *Request) {
		pt, err := prompts.GetPrompt(n)
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to get prompt %s", n)
			return
		}
		o.Prompt = pt
	}
}

func WithModel(n anthropic.Model) RequestParam {
	return func(o *Request) {
		o.Model = n
	}
}

// var defCliOptions = ClientOptions{Model: anthropic.ModelClaudeOpus4_6, MaxTokens: 20000, Temperature: 0.2}

//func DefaultClientOptions(opts ClientOptions) RequestParam {
//	opts.Model = util.StringCoalesce(opts.Model, defCliOptions.Model)
//	opts.Temperature = util.Float64Coalesce(opts.Temperature, defCliOptions.Temperature)
//	opts.MaxTokens = util.Int64Coalesce(opts.MaxTokens, defCliOptions.MaxTokens)
//	return opts
//}

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

func (br *BatchRequest) CustomIdWellFormed() (string, string, string, bool) {
	return BatchRequestCustomIdWellFormed(br.CustomID)
}

func BatchRequestCustomId(parts ...string) string {
	return strings.Join(parts, "_")
}

func BatchRequestCustomIdWellFormed(customId string) (string, string, string, bool) {
	const semLogContext = semLogContextBase + "batch-request::custom-id-well-formed"
	if customId == "" {
		err := fmt.Errorf("custom-id is required")
		log.Error().Err(err).Msg(semLogContext)
		return "", "", "", false
	}

	parts := strings.Split(customId, "_")
	switch len(parts) {
	case 0:
		fallthrough
	case 1:
		fallthrough
	case 2:
		err := fmt.Errorf("custom-id must be in the format <part-1>:<part-2>:<part-3>")
		log.Error().Err(err).Msg(semLogContext)
		return "", "", "", false
	case 3:
		return parts[0], parts[1], parts[2], true
	default:
		return parts[0], parts[1], strings.Join(parts[2:], ":"), true
	}

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
