package client

import (
	"encoding/json"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
)

// TurnEvent is sent on the progress channel (see WithProgress) after each
// completed turn in RunAgent.  The event carries enough information for
// callers to log a per-turn table, collect metrics, or drive a UI.
type TurnEvent struct {
	Turn       int             // 1-based turn number
	MaxTurns   int             // limit passed to RunAgent
	StopReason string          // API stop_reason for this turn
	ToolCalls  []ToolCallInfo  // non-empty when the turn ended with tool uses
	Usage      anthropic.Usage // token counts for this turn (not cumulative)
	Thinking   int64           // thinking tokens (OutputTokensDetails.ThinkingTokens)
	Duration   time.Duration   // wall-clock time for this turn's API call
}

// ToolCallInfo is a summary of one tool call within a turn.
type ToolCallInfo struct {
	Name   string
	Input  json.RawMessage
	Result string // tool execution output (as returned to the model)
}
