// Package client provides a clean, options-based wrapper around the Anthropic
// Go SDK for single-turn and multi-turn agentic message execution.
//
// Two entry points:
//
//   - Execute — one API call, no tools, returns the model's response text.
//   - RunAgent — streaming multi-turn loop: executes tool calls until the
//     model stops or maxTurns is reached.
//
// Usage:
//
//	sdk := anthropic.NewClient(option.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")))
//	c   := client.New(sdk)
//
//	// Single turn
//	resp, err := c.Execute(ctx,
//	    client.WithSystem("You are a helpful assistant."),
//	    client.WithUserText("Summarise this document."),
//	    client.WithMaxTokens(4096),
//	)
//
//	// Agentic loop
//	exec := fs.New(baseDir)
//	ts   := tools.New(tools.WithEditor(exec), tools.WithLister(exec))
//	agent, err := c.RunAgent(ctx,
//	    client.WithSystem(systemPrompt),
//	    client.WithUserText(userMsg),
//	    client.WithToolSet(ts),
//	    client.WithMaxTurns(30),
//	)
package client

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/rs/zerolog/log"
)

// Client wraps an Anthropic SDK client with an options-based execution API.
// Construct one via New and reuse it across calls — the underlying SDK client
// is safe for concurrent use.
type Client struct {
	sdk anthropic.Client
}

// New creates a Client from an already-configured Anthropic SDK client.
// Callers control API key, retries, timeout, and HTTP transport via the SDK
// client they pass in.
func New(sdk anthropic.Client) *Client {
	return &Client{sdk: sdk}
}

// Response is returned by Execute.
type Response struct {
	Text       string
	StopReason string
	Usage      anthropic.Usage
}

// AgentResponse is returned by RunAgent.
type AgentResponse struct {
	Text  string
	Turns int
	Usage anthropic.Usage // accumulated across all turns
}

// Execute sends a single non-streaming message and returns the response.
// Tools are ignored here; use RunAgent for workflows that require tool execution.
func (c *Client) Execute(ctx context.Context, opts ...Option) (*Response, error) {
	cfg := apply(opts)
	if len(cfg.content) == 0 {
		return nil, errors.New("client.Execute: user content is required (use WithUserText or WithUserContent)")
	}

	params := buildParams(cfg, []anthropic.MessageParam{
		anthropic.NewUserMessage(cfg.content...),
	})

	msg, err := c.sdk.Messages.New(ctx, params)
	if err != nil {
		return nil, err
	}
	if msg.StopReason == "max_tokens" {
		return nil, fmt.Errorf("client.Execute: response hit max_tokens (%d) — increase WithMaxTokens", cfg.maxTokens)
	}

	return &Response{
		Text:       joinTextBlocks(msg.Content),
		StopReason: string(msg.StopReason),
		Usage:      msg.Usage,
	}, nil
}

// RunAgent runs a streaming multi-turn agentic loop.
//
// Each turn:
//  1. Sends the current message history to the API (streaming).
//  2. If the model returns tool calls, executes them via the ToolSet and appends
//     results as a new user turn.
//  3. Repeats until the model stops without tool calls or maxTurns is reached.
//
// The final text response (all text blocks from the last assistant turn,
// joined with newlines) is returned in AgentResponse.Text.
func (c *Client) RunAgent(ctx context.Context, opts ...Option) (*AgentResponse, error) {
	cfg := apply(opts)
	if len(cfg.content) == 0 {
		return nil, errors.New("client.RunAgent: user content is required (use WithUserText or WithUserContent)")
	}

	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(cfg.content...),
	}

	var totalUsage anthropic.Usage

	for turn := 1; turn <= cfg.maxTurns; turn++ {
		params := buildParams(cfg, messages)

		turnStart := time.Now()
		stream := c.sdk.Messages.NewStreaming(ctx, params)
		defer stream.Close()

		var resp anthropic.Message
		for stream.Next() {
			if err := resp.Accumulate(stream.Current()); err != nil {
				// A turn that hits max_tokens mid tool-call leaves the tool_use
				// input JSON truncated; Accumulate then fails to marshal it
				// ("unexpected end of JSON input"), which crashes here before the
				// post-stream max_tokens guard below can run. Translate it into the
				// same actionable message instead of the opaque marshal error.
				if strings.Contains(err.Error(), "unexpected end of JSON input") {
					return nil, fmt.Errorf(
						"client.RunAgent: response truncated mid tool-call on turn %d — likely hit max_tokens (%d); increase WithMaxTokens: %w",
						turn, cfg.maxTokens, err,
					)
				}
				return nil, fmt.Errorf("client.RunAgent turn %d: accumulate: %w", turn, err)
			}
		}
		if err := stream.Err(); err != nil {
			return nil, fmt.Errorf("client.RunAgent turn %d: stream: %w", turn, err)
		}
		turnDuration := time.Since(turnStart)

		if resp.StopReason == "max_tokens" {
			return nil, fmt.Errorf(
				"client.RunAgent: hit max_tokens (%d) on turn %d — increase WithMaxTokens",
				cfg.maxTokens, turn,
			)
		}

		totalUsage.InputTokens += resp.Usage.InputTokens
		totalUsage.OutputTokens += resp.Usage.OutputTokens
		totalUsage.CacheReadInputTokens += resp.Usage.CacheReadInputTokens
		totalUsage.CacheCreationInputTokens += resp.Usage.CacheCreationInputTokens

		var toolUses []anthropic.ToolUseBlock
		var textParts []string
		for _, block := range resp.Content {
			switch v := block.AsAny().(type) {
			case anthropic.TextBlock:
				textParts = append(textParts, v.Text)
			case anthropic.ToolUseBlock:
				toolUses = append(toolUses, v)
			}
		}

		// Execute all tool calls and collect results before emitting the
		// progress event, so the event can carry each tool's result. On the
		// final turn (no tool uses) both loops are no-ops.
		var toolResults []anthropic.ContentBlockParamUnion
		results := make([]string, len(toolUses))
		for i, tu := range toolUses {
			if cfg.toolSet != nil {
				results[i] = cfg.toolSet.Execute(tu.Name, tu.Input)
			} else {
				results[i] = fmt.Sprintf("ERROR: no tool executor configured for %q", tu.Name)
			}
			toolResults = append(toolResults, anthropic.NewToolResultBlock(tu.ID, results[i], false))
		}

		if cfg.progress != nil {
			calls := make([]ToolCallInfo, len(toolUses))
			for i, tu := range toolUses {
				calls[i] = ToolCallInfo{Name: tu.Name, Input: tu.Input, Result: results[i]}
			}
			select {
			case cfg.progress <- TurnEvent{
				Turn:       turn,
				MaxTurns:   cfg.maxTurns,
				StopReason: string(resp.StopReason),
				ToolCalls:  calls,
				Usage:      resp.Usage,
				Thinking:   resp.Usage.OutputTokensDetails.ThinkingTokens,
				Duration:   turnDuration,
			}:
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		if len(toolUses) == 0 {
			return &AgentResponse{
				Text:  strings.Join(textParts, "\n"),
				Turns: turn,
				Usage: totalUsage,
			}, nil
		}

		// Append the assistant turn, preserving thinking blocks — the API
		// requires them in the conversation history when thinking is active —
		// followed by the tool results as the next user turn.
		messages = append(messages, anthropic.NewAssistantMessage(assistantContent(resp.Content)...))
		messages = append(messages, anthropic.NewUserMessage(toolResults...))
	}

	return nil, fmt.Errorf(
		"client.RunAgent: did not finish within %d turns — model may be stuck in a tool-use cycle",
		cfg.maxTurns,
	)
}

// BatchRequest is one entry in a SubmitBatch call. CustomID must be unique
// within the batch — results come back keyed by it, possibly out of request
// order. Options are the same per-request options accepted by Execute
// (WithModel, WithSystem, WithUserText, WithMaxTokens, WithThinking, …), so
// each request in a batch can use a different model, prompt, or settings.
type BatchRequest struct {
	CustomID string
	Options  []Option
}

// SubmitBatch submits multiple message requests as a single Message Batch.
//
// The batch is processed asynchronously: this returns the created batch as soon
// as it is accepted, not the responses. Poll c.sdk.Messages.Batches.Get with
// the returned batch ID, and once ProcessingStatus is "ended", collect results
// via c.sdk.Messages.Batches.ResultsStreaming — each result is keyed by the
// CustomID supplied here.
//
// Tools are ignored: batch requests are single-turn (there is no loop to
// execute tool calls against), matching Execute's behaviour.
func (c *Client) SubmitBatch(ctx context.Context, reqs ...BatchRequest) (*anthropic.MessageBatch, error) {
	if len(reqs) == 0 {
		return nil, errors.New("client.SubmitBatch: at least one request is required")
	}

	apiReqs := make([]anthropic.MessageBatchNewParamsRequest, 0, len(reqs))
	seen := make(map[string]struct{}, len(reqs))
	for i, r := range reqs {
		if r.CustomID == "" {
			return nil, fmt.Errorf("client.SubmitBatch: request %d has empty CustomID", i)
		}
		if _, dup := seen[r.CustomID]; dup {
			return nil, fmt.Errorf("client.SubmitBatch: duplicate CustomID %q", r.CustomID)
		}
		seen[r.CustomID] = struct{}{}

		cfg := apply(r.Options)
		if len(cfg.content) == 0 {
			return nil, fmt.Errorf("client.SubmitBatch: request %q has no user content (use WithUserText or WithUserContent)", r.CustomID)
		}

		apiReqs = append(apiReqs, anthropic.MessageBatchNewParamsRequest{
			CustomID: r.CustomID,
			Params:   buildBatchParams(cfg),
		})
	}

	return c.sdk.Messages.Batches.New(ctx, anthropic.MessageBatchNewParams{Requests: apiReqs})
}

// BatchStatus classifies the outcome of a single batch request. It mirrors the
// SDK's result union variants in a form that is easy to switch on.
type BatchStatus string

const (
	BatchSucceeded BatchStatus = "succeeded"
	BatchErrored   BatchStatus = "errored"
	BatchCanceled  BatchStatus = "canceled"
	BatchExpired   BatchStatus = "expired"
	// BatchTruncated is a synthetic status (not an API variant): the request was
	// accepted and "succeeded", but generation hit max_tokens, so Text holds a
	// truncated response. Surfaced distinctly to match Execute/RunAgent, which
	// treat max_tokens as an error.
	BatchTruncated BatchStatus = "truncated"
)

// BatchResult is the convenience form of a single batch response, keyed to the
// CustomID supplied in the corresponding BatchRequest. It flattens the SDK's
// result union: on success Text/StopReason/Usage are populated (mirroring
// Response); otherwise Status says why and Err carries the reason.
type BatchResult struct {
	CustomID   string
	Status     BatchStatus
	Text       string          // populated when Status == BatchSucceeded
	StopReason string          // populated when Status == BatchSucceeded
	Usage      anthropic.Usage // populated when Status == BatchSucceeded
	Err        error           // populated when Status != BatchSucceeded
}

// Succeeded reports whether the request completed with a full response.
// A max_tokens truncation has Status BatchTruncated (not BatchSucceeded), so it
// reports false here — matching Execute/RunAgent, which treat max_tokens as an error.
func (r BatchResult) Succeeded() bool { return r.Status == BatchSucceeded }

// GetBatch fetches the current state of a submitted batch. Inspect the returned
// batch's ProcessingStatus — results are only available once it is
// anthropic.MessageBatchProcessingStatusEnded.
func (c *Client) GetBatch(ctx context.Context, batchID string) (*anthropic.MessageBatch, error) {
	if batchID == "" {
		return nil, errors.New("client.GetBatch: batchID is required")
	}
	return c.sdk.Messages.Batches.Get(ctx, batchID, anthropic.MessageBatchGetParams{})
}

// CollectBatchResults streams every result of a completed batch and returns
// them as BatchResults, one per request, keyed by CustomID.
//
// The batch must have finished processing (ProcessingStatus ==
// anthropic.MessageBatchProcessingStatusEnded) — check via GetBatch first;
// results are unavailable before then. The returned slice order follows the
// results stream, which may differ from submission order; match by CustomID.
//
// A non-nil error is returned only for a transport/stream failure. Per-request
// failures (errored/canceled/expired) are reported in each BatchResult's Status
// and Err, not as the returned error.
func (c *Client) CollectBatchResults(ctx context.Context, batchID string) ([]BatchResult, error) {
	if batchID == "" {
		return nil, errors.New("client.CollectBatchResults: batchID is required")
	}

	stream := c.sdk.Messages.Batches.ResultsStreaming(ctx, batchID, anthropic.MessageBatchResultsParams{})
	defer stream.Close()

	var out []BatchResult
	for stream.Next() {
		out = append(out, toBatchResult(stream.Current()))
	}
	if err := stream.Err(); err != nil {
		return out, fmt.Errorf("client.CollectBatchResults: stream: %w", err)
	}
	return out, nil
}

// toBatchResult flattens one SDK individual response into a BatchResult.
func toBatchResult(resp anthropic.MessageBatchIndividualResponse) BatchResult {
	r := BatchResult{CustomID: resp.CustomID}
	switch v := resp.Result.AsAny().(type) {
	case anthropic.MessageBatchSucceededResult:
		// The Batches API has no max_tokens result variant: a request that hits
		// the cap comes back here as "succeeded" with StopReason == "max_tokens"
		// and truncated content. Surface it as an error to match Execute/RunAgent,
		// while still populating Text/Usage so callers can inspect the partial output.
		r.Text = joinTextBlocks(v.Message.Content)
		r.StopReason = string(v.Message.StopReason)
		r.Usage = v.Message.Usage
		if v.Message.StopReason == anthropic.StopReasonMaxTokens {
			r.Status = BatchTruncated
			r.Err = fmt.Errorf("request %q hit max_tokens — response truncated; increase WithMaxTokens", resp.CustomID)
		} else {
			r.Status = BatchSucceeded
		}
	case anthropic.MessageBatchErroredResult:
		r.Status = BatchErrored
		r.Err = fmt.Errorf("request %q errored: %s", resp.CustomID, v.Error.Error.Message)
	case anthropic.MessageBatchCanceledResult:
		r.Status = BatchCanceled
		r.Err = fmt.Errorf("request %q was canceled", resp.CustomID)
	case anthropic.MessageBatchExpiredResult:
		r.Status = BatchExpired
		r.Err = fmt.Errorf("request %q expired before processing", resp.CustomID)
	default:
		r.Status = BatchErrored
		r.Err = fmt.Errorf("request %q returned an unrecognised result type", resp.CustomID)
	}
	return r
}

// buildBatchParams builds the per-request params for a batch entry, reusing
// buildParams so temperature/thinking/tools handling stays defined in one place.
// The batch request-params type mirrors MessageNewParams field-for-field, so we
// build the standard params and copy the shared fields across.
func buildBatchParams(cfg config) anthropic.MessageBatchNewParamsRequestParams {
	p := buildParams(cfg, []anthropic.MessageParam{
		anthropic.NewUserMessage(cfg.content...),
	})
	return anthropic.MessageBatchNewParamsRequestParams{
		Model:        p.Model,
		MaxTokens:    p.MaxTokens,
		Messages:     p.Messages,
		System:       p.System,
		Temperature:  p.Temperature,
		Thinking:     p.Thinking,
		OutputConfig: p.OutputConfig,
	}
}

// buildParams constructs MessageNewParams from cfg and the current message history.
// Called once per turn in RunAgent; called once in Execute.
func buildParams(cfg config, messages []anthropic.MessageParam) anthropic.MessageNewParams {
	params := anthropic.MessageNewParams{
		Model:     cfg.model,
		MaxTokens: cfg.maxTokens,
		Messages:  messages,
	}

	if cfg.system != "" {
		params.System = []anthropic.TextBlockParam{{Text: cfg.system}}
	}

	if cfg.toolSet != nil {
		params.Tools = cfg.toolSet.Params()
	}

	// Thinking is driven by cfg.thinkingEffort ("" = off). The wire form is chosen
	// from the model: adaptive (output_config.effort) on 4.6+, budget_tokens on
	// older models. When thinking is actually applied, temperature must be omitted
	// (API requirement); if a requested budget can't fit max_tokens it is dropped
	// (and logged), and temperature is then allowed to apply as normal.
	thinkingApplied := false
	if cfg.thinkingEffort != "" {
		if modelSupportsAdaptiveThinking(cfg.model) {
			// Adaptive form: output_config.effort (Claude 4.6 and later).
			params.Thinking = anthropic.ThinkingConfigParamUnion{
				OfAdaptive: &anthropic.ThinkingConfigAdaptiveParam{},
			}
			params.OutputConfig.Effort = cfg.thinkingEffort
			thinkingApplied = true
		} else {
			// Explicit form: budget_tokens (models prior to Claude 4.6). Use the
			// caller-supplied budget when >0, else the effort->budget fallback map.
			budget := cfg.thinkingBudget
			if budget <= 0 {
				budget = effortToBudget(cfg.thinkingEffort)
			}
			if b, ok := clampThinkingBudget(budget, cfg.maxTokens); ok {
				params.Thinking = anthropic.ThinkingConfigParamOfEnabled(int64(b))
				thinkingApplied = true
			} else {
				log.Warn().
					Str("model", string(cfg.model)).
					Str("effort", string(cfg.thinkingEffort)).
					Int64("max-tokens", cfg.maxTokens).
					Msg(semLogContextBuildParams + " thinking dropped: max_tokens too small to fit the budget_tokens minimum (1024) on a non-adaptive model")
			}
		}
	}

	if !thinkingApplied && cfg.temperature != nil && modelSupportsTemperature(cfg.model) {
		params.Temperature = anthropic.Float(*cfg.temperature)
	}

	// JSON structured output: constrain the response to the supplied JSON schema.
	// Set as a sub-field so it composes with output_config.effort above rather than
	// clobbering it.
	if len(cfg.outputSchema) > 0 {
		params.OutputConfig.Format = anthropic.JSONOutputFormatParam{Schema: cfg.outputSchema}
	}

	return params
}

const semLogContextBuildParams = "anthropiclks-client::buildParams"

// modelsWithAdaptiveThinking lists the model families that use the adaptive
// thinking form (output_config.effort). These are Claude 4.6 and later; on them
// the budget_tokens form is removed and returns a 400. Matched as substrings so
// aliases and dated snapshots both hit (mirrors modelsWithoutSampling). Anything
// not listed here (opus-4-5, sonnet-4-5, haiku-4-5, sonnet-3-7, …) uses
// budget_tokens — note capability is family-specific, not numeric (haiku-4-5 is
// budget-only despite the "4-5").
var modelsWithAdaptiveThinking = []string{
	"opus-4-6",
	"sonnet-4-6",
	"opus-4-7",
	"opus-4-8",
	"sonnet-5",
	"opus-5",
	"fable-5",
	"mythos-5",
}

// modelSupportsAdaptiveThinking reports whether the model uses the adaptive
// thinking form (output_config.effort) rather than budget_tokens.
func modelSupportsAdaptiveThinking(m anthropic.Model) bool {
	s := strings.ToLower(string(m))
	for _, id := range modelsWithAdaptiveThinking {
		if strings.Contains(s, id) {
			return true
		}
	}
	return false
}

// effortToBudget maps an effort level to an approximate budget_tokens value for
// the non-adaptive (pre-4.6) models. Anthropic publishes no official equivalence
// — effort is qualitative, budget_tokens is a hard count — so these are pragmatic
// defaults, tunable and only used as a fallback when the caller passes no explicit
// budget. Unknown/empty effort falls back to the medium value.
func effortToBudget(effort anthropic.OutputConfigEffort) int {
	switch effort {
	case anthropic.OutputConfigEffortLow:
		return 4096
	case anthropic.OutputConfigEffortHigh:
		return 16384
	case anthropic.OutputConfigEffortXhigh:
		return 24576
	case anthropic.OutputConfigEffortMax:
		return 32768
	default: // medium and anything unrecognized
		return 8192
	}
}

// clampThinkingBudget fits a budget_tokens value into the API's constraints:
// it must be >= 1024 and strictly < max_tokens. Returns ok=false when max_tokens
// is too small to satisfy both (i.e. <= 1024), so the caller drops thinking.
func clampThinkingBudget(budget int, maxTokens int64) (int, bool) {
	const minBudget = 1024
	if maxTokens <= minBudget {
		return 0, false
	}
	if budget < minBudget {
		budget = minBudget
	}
	if int64(budget) >= maxTokens {
		budget = int(maxTokens) - 1
	}
	return budget, true
}

// modelsWithoutSampling lists the model families that removed the sampling
// parameters (temperature/top_p/top_k): sending temperature to them returns a
// 400 "`temperature` is deprecated for this model." Matched as substrings so
// aliases and dated snapshots (e.g. claude-opus-4-7, claude-opus-4-8@...) both
// hit. Opus 4.6, Sonnet 4.6, Haiku 4.5, and older still accept temperature.
var modelsWithoutSampling = []string{
	"opus-4-7",
	"opus-4-8",
	"fable-5",
	"mythos-5",
}

// modelSupportsTemperature reports whether the given model still accepts the
// temperature parameter. When it does not, buildParams silently omits any
// temperature the caller set via WithTemperature, so a WithTemperature call is
// a no-op on those models rather than a request-breaking 400.
func modelSupportsTemperature(m anthropic.Model) bool {
	s := strings.ToLower(string(m))
	for _, id := range modelsWithoutSampling {
		if strings.Contains(s, id) {
			return false
		}
	}
	return true
}

// assistantContent converts a response content slice to the param union form
// expected when appending an assistant turn to the message history.
// Thinking and redacted-thinking blocks are preserved — the API requires them
// in the history when extended thinking is enabled.
func assistantContent(blocks []anthropic.ContentBlockUnion) []anthropic.ContentBlockParamUnion {
	out := make([]anthropic.ContentBlockParamUnion, 0, len(blocks))
	for _, block := range blocks {
		switch v := block.AsAny().(type) {
		case anthropic.TextBlock:
			out = append(out, anthropic.NewTextBlock(v.Text))
		case anthropic.ToolUseBlock:
			out = append(out, anthropic.NewToolUseBlock(v.ID, v.Input, v.Name))
		case anthropic.ThinkingBlock:
			out = append(out, anthropic.NewThinkingBlock(v.Signature, v.Thinking))
		case anthropic.RedactedThinkingBlock:
			out = append(out, anthropic.NewRedactedThinkingBlock(v.Data))
		}
	}
	return out
}

// joinTextBlocks concatenates the text of all TextBlock entries in a content slice.
func joinTextBlocks(blocks []anthropic.ContentBlockUnion) string {
	var parts []string
	for _, b := range blocks {
		if tb, ok := b.AsAny().(anthropic.TextBlock); ok {
			parts = append(parts, tb.Text)
		}
	}
	return strings.Join(parts, "\n")
}
