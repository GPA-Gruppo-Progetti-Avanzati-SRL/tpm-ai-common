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

		if cfg.progress != nil {
			calls := make([]ToolCallInfo, len(toolUses))
			for i, tu := range toolUses {
				calls[i] = ToolCallInfo{Name: tu.Name, Input: tu.Input}
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
		// requires them in the conversation history when thinking is active.
		messages = append(messages, anthropic.NewAssistantMessage(assistantContent(resp.Content)...))

		// Execute all tool calls and collect results.
		var toolResults []anthropic.ContentBlockParamUnion
		for _, tu := range toolUses {
			var result string
			if cfg.toolSet != nil {
				result = cfg.toolSet.Execute(tu.Name, tu.Input)
			} else {
				result = fmt.Sprintf("ERROR: no tool executor configured for %q", tu.Name)
			}
			toolResults = append(toolResults, anthropic.NewToolResultBlock(tu.ID, result, false))
		}
		messages = append(messages, anthropic.NewUserMessage(toolResults...))
	}

	return nil, fmt.Errorf(
		"client.RunAgent: did not finish within %d turns — model may be stuck in a tool-use cycle",
		cfg.maxTurns,
	)
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

	if cfg.adaptiveThinking != nil {
		// Adaptive form: output_config.effort (e.g. claude-opus-4-7 and later).
		// Temperature must be omitted when thinking is enabled (API requirement).
		params.Thinking = anthropic.ThinkingConfigParamUnion{
			OfAdaptive: &anthropic.ThinkingConfigAdaptiveParam{},
		}
		params.OutputConfig = anthropic.OutputConfigParam{
			Effort: *cfg.adaptiveThinking,
		}
	} else if cfg.thinking > 0 {
		// Explicit form: budget_tokens (models prior to claude-opus-4-7).
		// Temperature must be omitted when thinking is enabled (API requirement).
		params.Thinking = anthropic.ThinkingConfigParamOfEnabled(int64(cfg.thinking))
	} else if cfg.temperature != nil {
		params.Temperature = anthropic.Float(*cfg.temperature)
	}

	return params
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
