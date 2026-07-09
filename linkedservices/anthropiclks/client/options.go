package client

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/anthropiclks/tools"
	anthropic "github.com/anthropics/anthropic-sdk-go"
)

// Option configures a single Execute or RunAgent call.
type Option func(*config)

type config struct {
	model            anthropic.Model
	maxTokens        int64
	temperature      *float64 // nil means omit (required when thinking is enabled)
	system           string
	content          []anthropic.ContentBlockParamUnion
	thinking         int                           // >0 = explicit budget_tokens form
	adaptiveThinking *anthropic.OutputConfigEffort // non-nil = output_config.effort form
	toolSet          *tools.ToolSet
	maxTurns         int
	progress         chan<- TurnEvent
}

func newConfig() config {
	return config{
		model:     anthropic.ModelClaudeSonnet4_6,
		maxTokens: 8192,
		maxTurns:  50,
	}
}

func apply(opts []Option) config {
	cfg := newConfig()
	for _, o := range opts {
		o(&cfg)
	}
	return cfg
}

// WithModel sets the model. Defaults to claude-sonnet-4-6.
func WithModel(m anthropic.Model) Option {
	return func(c *config) { c.model = m }
}

// WithMaxTokens sets the maximum output tokens. Defaults to 8192.
func WithMaxTokens(n int64) Option {
	return func(c *config) {
		if n > 0 {
			c.maxTokens = n
		}
	}
}

// WithTemperature sets the sampling temperature.
// Ignored when WithThinking is also set (API requirement: temperature must be
// absent when extended thinking is enabled), and ignored on models that removed
// the sampling parameters (Opus 4.7/4.8, Fable 5, Mythos 5), where sending it
// returns a 400 — see modelSupportsTemperature.
func WithTemperature(f float64) Option {
	return func(c *config) { c.temperature = &f }
}

// WithSystem sets the system prompt.
func WithSystem(s string) Option {
	return func(c *config) { c.system = s }
}

// WithUserText sets a single plain-text user message.
func WithUserText(text string) Option {
	return func(c *config) {
		c.content = []anthropic.ContentBlockParamUnion{anthropic.NewTextBlock(text)}
	}
}

// WithUserContent sets the user message content blocks directly, allowing
// cached blocks, images, or multi-part messages.
func WithUserContent(blocks ...anthropic.ContentBlockParamUnion) Option {
	return func(c *config) { c.content = blocks }
}

// WithThinking enables extended thinking using the explicit budget_tokens form.
// Use for models that support this form (e.g. claude-opus-4-5 and earlier).
// When set, temperature is automatically omitted from the request.
// Mutually exclusive with WithAdaptiveThinking — the last one set wins.
func WithThinking(budgetTokens int) Option {
	return func(c *config) {
		if budgetTokens > 0 {
			c.thinking = budgetTokens
			c.adaptiveThinking = nil
		}
	}
}

// WithAdaptiveThinking enables extended thinking using the output_config.effort
// form. Use for models that support this form (e.g. claude-opus-4-7 and later).
// When set, temperature is automatically omitted from the request.
// Mutually exclusive with WithThinking — the last one set wins.
func WithAdaptiveThinking(effort anthropic.OutputConfigEffort) Option {
	return func(c *config) {
		c.adaptiveThinking = &effort
		c.thinking = 0
	}
}

// WithToolSet attaches a ToolSet for RunAgent. Ignored by Execute.
func WithToolSet(ts *tools.ToolSet) Option {
	return func(c *config) { c.toolSet = ts }
}

// WithMaxTurns sets the maximum number of agentic turns for RunAgent.
// Defaults to 50. Ignored by Execute.
func WithMaxTurns(n int) Option {
	return func(c *config) {
		if n > 0 {
			c.maxTurns = n
		}
	}
}

// WithProgress sets a channel that receives a TurnEvent after every completed
// turn in RunAgent.  The send blocks until the receiver reads, so the channel
// should be adequately buffered or consumed in a separate goroutine to avoid
// stalling the agent loop.  Ignored by Execute.
func WithProgress(ch chan<- TurnEvent) Option {
	return func(c *config) { c.progress = ch }
}
