package client

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/anthropiclks/tools"
	anthropic "github.com/anthropics/anthropic-sdk-go"
)

// Option configures a single Execute or RunAgent call.
type Option func(*config)

type config struct {
	model          anthropic.Model
	maxTokens      int64
	temperature    *float64 // nil means omit (required when thinking is enabled)
	system         string
	content        []anthropic.ContentBlockParamUnion
	thinkingEffort anthropic.OutputConfigEffort // "" = no thinking; the model-independent driver
	thinkingBudget int                          // >0 = explicit budget_tokens override for non-adaptive models; <=0 = use fallback map
	toolSet        *tools.ToolSet
	maxTurns       int
	progress       chan<- TurnEvent
	outputSchema   map[string]any // non-empty = request JSON structured output against this JSON schema
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

// WithThinking enables extended thinking. effort is the model-independent driver,
// expressed as an output_config.effort level (low/medium/high/xhigh/max); the
// client picks the wire form from the model:
//
//   - On models that support adaptive thinking (Claude 4.6 and later), effort is
//     sent as output_config.effort. Any budget argument is ignored — the
//     budget_tokens form returns a 400 on those models.
//   - On older models (pre-4.6), effort is translated to the budget_tokens form.
//     Pass an optional budget to set that value explicitly; omit it (or pass <=0)
//     to use the built-in effort->budget fallback map. The budget is clamped to
//     [1024, max_tokens); if max_tokens is too small to fit even the 1024 minimum,
//     thinking is dropped (and logged) rather than sent as a request-breaking value.
//
// Passing an empty effort disables thinking. When thinking is enabled, temperature
// is automatically omitted from the request.
func WithThinking(effort anthropic.OutputConfigEffort, budget ...int) Option {
	return func(c *config) {
		c.thinkingEffort = effort
		c.thinkingBudget = 0
		if len(budget) > 0 && budget[0] > 0 {
			c.thinkingBudget = budget[0]
		}
	}
}

// WithOutputSchema requests JSON structured output constrained to the given JSON
// schema (sent as output_config.format = json_schema). The schema map must be a
// valid Anthropic structured-output schema (root object, additionalProperties:
// false, all properties required). Applies to both Execute and the batch path.
func WithOutputSchema(schema map[string]any) Option {
	return func(c *config) { c.outputSchema = schema }
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
