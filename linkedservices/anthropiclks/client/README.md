# client — Anthropic message client

A thin, options-based wrapper around the Anthropic Go SDK for single-turn and
multi-turn agentic execution.

## Design

### Two entry points

| Method     | When to use                                                                                    |
|------------|------------------------------------------------------------------------------------------------|
| `Execute`  | One API call, no tools. Returns the model's text response.                                     |
| `RunAgent` | Streaming multi-turn loop. Executes tool calls until the model stops or `maxTurns` is reached. |

Both methods accept the same functional-option set. Options that do not apply
to a given method (e.g. `WithToolSet`, `WithMaxTurns`, `WithProgress` on
`Execute`) are silently ignored.

### Construction

```go
sdk := anthropic.NewClient(option.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")))
c   := client.New(sdk)
```

`Client` wraps the SDK client and is safe for concurrent use. Create one and
reuse it across calls — each call is stateless.

### Options reference

| Option                         | Default             | Description                                                           |
|--------------------------------|---------------------|-----------------------------------------------------------------------|
| `WithModel(m)`                 | `claude-sonnet-4-6` | Anthropic model to use                                                |
| `WithMaxTokens(n)`             | `8192`              | Maximum output tokens                                                 |
| `WithTemperature(f)`           | omitted             | Sampling temperature (omitted automatically when thinking is enabled) |
| `WithSystem(s)`                | `""`                | System prompt                                                         |
| `WithUserText(s)`              | —                   | Single plain-text user message                                        |
| `WithUserContent(blocks...)`   | —                   | Multi-part user message (cached blocks, images, …)                    |
| `WithThinking(budget)`         | disabled            | Extended thinking — explicit `budget_tokens` form (pre-4.7 models)    |
| `WithAdaptiveThinking(effort)` | disabled            | Extended thinking — `output_config.effort` form (4-7 and later)       |
| `WithToolSet(ts)`              | `nil`               | Tool set for `RunAgent`                                               |
| `WithMaxTurns(n)`              | `50`                | Turn limit for `RunAgent`                                             |
| `WithProgress(ch)`             | `nil`               | Progress channel for `RunAgent` (see below)                           |

### Extended thinking

Two options cover the two thinking API forms. You choose based on what your
model supports — the client makes no assumptions about the model ID.

| Option                                            | API form                               | Use with                        |
|---------------------------------------------------|----------------------------------------|---------------------------------|
| `WithThinking(budget int)`                        | `{"type":"enabled","budget_tokens":N}` | Models prior to claude-opus-4-7 |
| `WithAdaptiveThinking(effort OutputConfigEffort)` | `output_config.effort`                 | claude-opus-4-7 and later       |

In both cases temperature is automatically omitted (API requirement). The two
options are mutually exclusive — the last one set on a call wins.

```go
// Adaptive form — claude-opus-4-7 and any future model that uses effort:
c.RunAgent(ctx,
    client.WithModel(anthropic.ModelClaudeOpus4_7),
    client.WithAdaptiveThinking(anthropic.OutputConfigEffortHigh),
    ...
)

// Explicit budget form — older models:
c.RunAgent(ctx,
    client.WithModel(anthropic.ModelClaudeOpus4_5),
    client.WithThinking(10_000),
    ...
)
```

Thinking blocks and redacted-thinking blocks are preserved verbatim in the
conversation history — the API enforces this cryptographically and rejects
requests that omit them.

---

## Usage

### Single-turn

```go
resp, err := c.Execute(ctx,
    client.WithSystem("You are a helpful assistant."),
    client.WithUserText("Summarise this document."),
    client.WithMaxTokens(4096),
)
// resp.Text       — model response
// resp.StopReason — "end_turn", etc.
// resp.Usage      — token counts
```

### Agentic loop

```go
exec := fs.New(baseDir)
ts   := tools.New(tools.WithEditor(exec), tools.WithLister(exec))

agent, err := c.RunAgent(ctx,
    client.WithSystem(systemPrompt),
    client.WithUserText(userMsg),
    client.WithToolSet(ts),
    client.WithMaxTurns(30),
)
// agent.Text  — final text from the last assistant turn
// agent.Turns — number of turns consumed
// agent.Usage — token usage accumulated across all turns
```

---

## Progress reporting

`WithProgress` attaches a `chan<- TurnEvent` that `RunAgent` writes to after
every completed turn. This lets callers log a status table, drive a TUI, or
collect metrics without modifying the core loop.

### TurnEvent fields

```go
type TurnEvent struct {
    Turn       int             // 1-based turn number
    MaxTurns   int             // limit passed to RunAgent
    StopReason string          // API stop_reason ("tool_use", "end_turn", …)
    ToolCalls  []ToolCallInfo  // populated when StopReason == "tool_use"
    Usage      anthropic.Usage // token counts for this turn only (not cumulative)
    Thinking   int64           // thinking tokens (OutputTokensDetails.ThinkingTokens)
    Duration   time.Duration   // wall-clock time for this turn's API call
}

type ToolCallInfo struct {
    Name  string
    Input json.RawMessage
}
```

### Blocking behaviour

The send is wrapped in a `select` that also watches `ctx.Done()`. This means:

- If the receiver is keeping up (buffered channel or goroutine), the loop
  continues immediately.
- If the context is cancelled while the send is pending, `RunAgent` returns
  `ctx.Err()` instead of blocking forever.

### Typical pattern — goroutine consumer

```go
progress := make(chan client.TurnEvent, 16) // buffer avoids stalling the loop
go func() {
    for ev := range progress {
        logTurnTable(ev) // your display/logging logic
    }
}()

resp, err := c.RunAgent(ctx,
    client.WithSystem(system),
    client.WithUserContent(content...),
    client.WithToolSet(ts),
    client.WithMaxTurns(40),
    client.WithProgress(progress),
)
close(progress) // signals the goroutine to exit
```

### Recreating the markdown progress table

```go
func logTurnTable(ev client.TurnEvent) {
    if ev.Turn == 1 {
        fmt.Println("| Turn | Status | In | Out | Time | Tools |")
        fmt.Println("|-----:|:-------|---:|----:|-----:|:------|")
    }

    var status string
    if len(ev.ToolCalls) > 0 {
        word := "tool calls"
        if len(ev.ToolCalls) == 1 {
            word = "tool call"
        }
        status = fmt.Sprintf("%d %s", len(ev.ToolCalls), word)
    } else {
        status = ev.StopReason
    }

    toolLabels := make([]string, len(ev.ToolCalls))
    for i, tc := range ev.ToolCalls {
        toolLabels[i] = "`" + tc.Name + "`"
    }

    fmt.Printf("| %2d/%d | %s | %d | %d | %.1fs | %s |\n",
        ev.Turn, ev.MaxTurns,
        status,
        ev.Usage.InputTokens,
        ev.Usage.OutputTokens,
        ev.Duration.Seconds(),
        strings.Join(toolLabels, " · "),
    )
}
```
