# ReAct Gap Analysis — `anthropiclks/client`

Analysis of `client.go` (plus `options.go`, `event.go`, and the `tools` package it
depends on) against the ReAct agent architecture described in *"Implementing a basic
ReAct agent"* (chapter 4). The goal is to identify which architectural components the
book prescribes are present, missing, or only partially realized — and to sketch the
two abstractions that would close the most important gaps.

> Scope: findings + proposed sketch only. No production code is modified.

---

## 0. Continuation / resume state

**Status as of 2026-09-09:** analysis + sketch complete. **No production code changed.**
The only file written this session is this document (`client/GAP-analysis.md`) plus the
untracked scratch dir `linkedservices/anthropiclks/loop/` (pre-existing, unrelated).

**How to resume:** read §2 (gaps) for *what* and §4 (sketch) for *how*. Pick up at §3
priority order. Nothing here has been committed to code, so there is no half-done edit to
reconcile — implementation starts from a clean tree.

### Files reviewed (the evidence base)
| File | Role in the analysis |
|---|---|
| `client/client.go` | The implementation under review — `RunAgent` is the ReAct loop. |
| `client/options.go` | `config` + `With*` options; where new options (`WithOutputSchema`, context injection) would attach. |
| `client/event.go` | `TurnEvent`/`ToolCallInfo` — progress channel only; not a persisted trace. |
| `tools/tools.go` | `ToolSet` — the hardcoded 2-tool `switch` to be replaced by a registry. |

### Decision log (settled)
- Close the two foundational gaps first: `ExecutionContext` and a `Tool` interface. MCP
  and structured output are explicitly deferred (see §4.5) and depend on those two.
- Preserve current behaviour: existing text-editor + `list_files` tools become registered
  `Tool` implementations; `TurnEvent` progress channel stays *in addition to* the new
  persisted event trail.
- Follow the package's existing style: options-based config, reuse `anthropic.*` SDK
  types, small focused types.

### Open questions (decide before/while implementing)
1. **`ExecContext` location.** Sketch defines a narrow `ExecContext` interface inside the
   `tools` package to avoid a `client` ↔ `tools` import cycle. Alternative: move the
   context type into a shared subpackage. → *undecided.*
2. **`FinalResult` type.** Sketched as `string`. Structured output (§4.5) will likely need
   `any` / an interface to hold a typed result (the book's `str | BaseModel`). Decide
   whether to make it `any` now or refactor later. → *undecided.*
3. **Context ownership.** Should `RunAgent` always create the `ExecutionContext`, or accept
   an optional caller-supplied one (book's `run(user_input, context=None)` pattern) to allow
   multi-turn continuation / injected trace IDs? → *leaning: accept optional, create if nil.*

### Next steps (in order)
1. Add `ExecutionContext` + `Event`/`ContentItem` types (new file, e.g. `client/context.go`).
2. Add `Tool` + `ExecContext` interfaces and `FunctionTool` in `tools`; convert `ToolSet`
   to a registry; re-register the two existing tools to prove behaviour is preserved.
3. Thread the context through `RunAgent`; record events; add `Context` to `AgentResponse`.
4. (Deferred) MCP loader → `[]Tool`. (Deferred) `WithOutputSchema` + `final_answer` tool.

---

## 1. What already matches the book

| Book component | Where in code | Notes |
|---|---|---|
| Think–act loop (§4.1, §4.6) | `RunAgent` in `client.go` | Streams a turn, extracts `ToolUseBlock`s, executes them, appends results as the next user turn, repeats. This is the tool-calling ReAct cycle. |
| Graceful tool error handling (§4.6.4 `act()`) | `ToolSet.Execute` in `tools/tools.go` | Returns `"ERROR: ..."` strings back into the conversation instead of crashing, so the model can adapt. Matches the book's try/except-as-observation pattern. |
| Loop guard (`max_steps`) | `maxTurns` in `options.go` / `RunAgent` | Clear "stuck in a tool-use cycle" error when the limit is hit. |
| `max_tokens` handling | `RunAgent`, `toBatchResult` | More robust than the book — includes a mid-tool-call truncation guard. |

---

## 2. Gaps, mapped to the book's components

### 2.1 No `ExecutionContext` (§4.3) — the central gap

The book's core thesis is consolidating execution state into **one container** rather
than scattering it across local variables. `RunAgent` does exactly the scattering the
book warns against: `messages`, `turn`, and `totalUsage` are loop-local and discarded
when the function returns. Specifically absent:

- **`execution_id`** — no run identifier to correlate logs/errors across a session.
- **`state`** — no scratch-pad key/value store that tools or later steps can read/write.
- **`final_result`** — no explicit completion flag; termination is *inferred* from
  "zero tool uses this turn."
- **Persisted, inspectable event history** — `TurnEvent`/`ToolCallInfo` (`event.go`) is
  *emitted* on a progress channel and then gone. There is no equivalent of the book's
  stored `events []Event` or `AgentResult.context` that lets you inspect the full trace
  after the run (no `display_trace` equivalent). `AgentResponse` returns only
  `Text/Turns/Usage`.

### 2.2 No unified tool abstraction (§4.4)

The book's point is a `BaseTool` interface so any local function or MCP tool plugs in
uniformly. Here, `ToolSet` is a hardcoded `switch` over exactly two tools (the built-in
text editor and `list_files`). Consequences:

- **Not extensible** — adding a tool means editing `tools.go`; you cannot register an
  arbitrary local function. There is no `FunctionTool`/`@tool` equivalent and no schema
  generation from a signature.
- **No context propagation to tools** — `Execute(name, rawInput)` receives only the tool
  name and raw JSON. The book's key design decision is that every tool's `execute`
  receives the `ExecutionContext`, so tools can read shared state or check permissions.
  That channel does not exist here (and cannot, since there is no context to pass).

### 2.3 No MCP tool integration (§4.4.4)

There is no `load_mcp_tools` equivalent. Tools are local-only. The book treats "local
functions and MCP tools behind one interface" as a headline feature.

### 2.4 No structured output (§4.7)

`RunAgent` returns free-form `Text` only. There is no `output_type`, no dynamically
created `final_answer` tool, and no `tool_choice: "required"` to force a typed result.
(`buildParams` sets `OutputConfig.Effort`, but that is thinking-effort, unrelated to an
output schema.) So there is no type-safe / guaranteed-format result path.

### 2.5 Minor: no `run` / `step` / `think` / `act` separation (§4.6)

Everything is inlined into one `RunAgent` function. Functionally fine, but the book's
modular decomposition is what makes `_prepare_llm_request` the "context engineering"
seam — the place to summarize history, drop stale tool results, or inject memory later.
Here there is no such seam.

---

## 3. Priority

The two changes that unlock the most (and that later chapters build on) are:

1. **Introduce an `ExecutionContext`** holding `executionID`, an appended event log, a
   `state map[string]any`, and `finalResult` — then thread it through the loop and
   return it on `AgentResponse` for inspection.
2. **Define a `Tool` interface** (`Name` / `Definition` / `Execute(ctx, input)`) so
   `ToolSet` dispatches over registered tools generically and passes context in — which
   also makes MCP tools and a `final_answer` structured-output tool droppable later.

---

## 4. Proposed sketch

> Design notes: these follow the package's existing conventions — options-based
> configuration, small focused types, SDK types (`anthropic.*`) reused rather than
> re-wrapped. The sketch is illustrative, not compiled; field/method names are chosen to
> read naturally against `RunAgent`.

### 4.1 `ExecutionContext`

```go
package client

import (
	"encoding/json"
	"sync"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
)

// ContentKind classifies a recorded content item — the Go analogue of the book's
// discriminated ContentItem union (Message | ToolCall | ToolResult).
type ContentKind string

const (
	KindMessage    ContentKind = "message"
	KindToolCall   ContentKind = "tool_call"
	KindToolResult ContentKind = "tool_result"
)

// ContentItem is one recorded occurrence inside an Event.
type ContentItem struct {
	Kind ContentKind `json:"kind"`

	// KindMessage
	Role string `json:"role,omitempty"` // "user" | "assistant" | "system"
	Text string `json:"text,omitempty"`

	// KindToolCall / KindToolResult
	ToolCallID string          `json:"tool_call_id,omitempty"`
	ToolName   string          `json:"tool_name,omitempty"`
	Input      json.RawMessage `json:"input,omitempty"`  // tool_call
	Result     string          `json:"result,omitempty"` // tool_result
	IsError    bool            `json:"is_error,omitempty"`
}

// Event wraps content with metadata for the audit trail (book §4.3, Event class).
type Event struct {
	ExecutionID string        `json:"execution_id"`
	Timestamp   time.Time     `json:"timestamp"`
	Author      string        `json:"author"` // "user" or the agent name
	Content     []ContentItem `json:"content"`
}

// ExecutionContext is the central store for a single RunAgent invocation
// (book §4.3). It replaces the scattered loop-local variables (messages, turn,
// totalUsage) with one container that can be threaded to tools and returned to
// the caller for inspection.
type ExecutionContext struct {
	ExecutionID string
	CurrentStep int
	Events      []Event
	Usage       anthropic.Usage
	FinalResult string

	mu    sync.Mutex     // guards State for context-aware tools that run concurrently
	State map[string]any // scratch-pad for tools / intermediate results
}

// NewExecutionContext seeds a context. The caller supplies the ID so it can be
// generated deterministically or injected from an outer trace/request ID.
func NewExecutionContext(executionID string) *ExecutionContext {
	return &ExecutionContext{
		ExecutionID: executionID,
		State:       make(map[string]any),
	}
}

func (c *ExecutionContext) AddEvent(e Event)       { c.Events = append(c.Events, e) }
func (c *ExecutionContext) IncrementStep()         { c.CurrentStep++ }
func (c *ExecutionContext) Done() bool             { return c.FinalResult != "" }

// Get / Set give tools race-free access to the shared scratch-pad.
func (c *ExecutionContext) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.State[key]
	return v, ok
}

func (c *ExecutionContext) Set(key string, val any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.State[key] = val
}
```

### 4.2 `Tool` interface

```go
package tools

import (
	"context"
	"encoding/json"

	anthropic "github.com/anthropics/anthropic-sdk-go"
)

// ExecContext is the minimal view of the client's ExecutionContext that a tool is
// allowed to touch. Defined as an interface (in the tools package) to avoid an
// import cycle with the client package, and to keep the tool contract narrow —
// tools get scratch-pad access, not the whole loop.
type ExecContext interface {
	Get(key string) (any, bool)
	Set(key string, val any)
	ExecutionID() string
}

// Tool is the unified abstraction (book §4.4, BaseTool). Every tool — local
// function, built-in editor, or MCP-backed — implements it, so ToolSet can
// dispatch generically and pass context in.
type Tool interface {
	// Name is the identifier the model uses in tool_use blocks.
	Name() string

	// Param returns the SDK tool definition sent to the API.
	Param() anthropic.ToolUnionParam

	// Execute runs the tool. rawInput is the raw JSON from the ToolUseBlock.
	// The returned string is fed back to the model as the tool_result; on
	// failure it should be an "ERROR: ..." string (matching current behaviour),
	// with isError signalling the tool_result error flag.
	Execute(ctx context.Context, ec ExecContext, rawInput json.RawMessage) (result string, isError bool)
}
```

### 4.3 `FunctionTool` adapter (wraps a plain Go func — book §4.4.3)

```go
// FunctionTool adapts an ordinary handler into a Tool, mirroring the book's
// FunctionTool. Context-aware handlers receive ec; context-free ones ignore it.
type FunctionTool struct {
	name  string
	param anthropic.ToolUnionParam
	fn    func(ctx context.Context, ec ExecContext, in json.RawMessage) (string, bool)
}

func NewFunctionTool(
	name string,
	param anthropic.ToolUnionParam,
	fn func(ctx context.Context, ec ExecContext, in json.RawMessage) (string, bool),
) *FunctionTool {
	return &FunctionTool{name: name, param: param, fn: fn}
}

func (t *FunctionTool) Name() string                    { return t.name }
func (t *FunctionTool) Param() anthropic.ToolUnionParam { return t.param }
func (t *FunctionTool) Execute(ctx context.Context, ec ExecContext, in json.RawMessage) (string, bool) {
	return t.fn(ctx, ec, in)
}
```

### 4.4 How `ToolSet` and `RunAgent` change (illustrative)

`ToolSet` becomes a registry over `Tool`s rather than a hardcoded switch:

```go
type ToolSet struct {
	tools map[string]Tool // keyed by Name()
}

func (ts *ToolSet) Register(t Tool) { ts.tools[t.Name()] = t }

func (ts *ToolSet) Params() []anthropic.ToolUnionParam {
	out := make([]anthropic.ToolUnionParam, 0, len(ts.tools))
	for _, t := range ts.tools {
		out = append(out, t.Param())
	}
	return out
}

func (ts *ToolSet) Execute(ctx context.Context, ec ExecContext, name string, in json.RawMessage) (string, bool) {
	t, ok := ts.tools[name]
	if !ok {
		return "ERROR: unknown tool '" + name + "'", true
	}
	return t.Execute(ctx, ec, in)
}
```

The existing text-editor and `list_files` tools become `Tool` implementations
registered by default, so current behaviour is preserved. `RunAgent` then:

1. Creates (or accepts) an `*ExecutionContext` instead of loose locals.
2. Records each turn as an `Event` on `ctx.Events` (audit trail) **in addition to**
   emitting the existing `TurnEvent` on the progress channel.
3. Passes the context into `toolSet.Execute(...)` so tools can use `State`.
4. Returns the context on `AgentResponse` (e.g. a new `Context *ExecutionContext`
   field) so callers can inspect the full trace — the `AgentResult.context` of the book.

### 4.5 Follow-on (out of scope for this sketch)

Once the two abstractions above exist, the remaining book features drop in cleanly:

- **MCP tools** — a `LoadMCPTools(...)` constructor returning `[]Tool` that
  `ToolSet.Register`s, each wrapping an MCP call in a `FunctionTool`.
- **Structured output** — a `WithOutputSchema(...)` option that registers a synthetic
  `final_answer` `Tool`, with the loop treating a successful `final_answer` call as
  `FinalResult`. The book (§4.7) uses this "tool as output formatter" pattern; adopt it
  for `RunAgent` (it unifies typed output with loop termination).

  **The plain-language version first (read this before the details below).**
  It's easy to tangle the parts, so here is the whole mechanism as simple steps:

  1. **Registering the tool just puts it on the menu.** You create a `final_answer` tool
     whose input *schema* is the shape you want back. Registering it does **not** make the
     model use it — it only makes it available. (Why a tool at all? Because a tool's input
     is the *only* channel through which the model emits schema-validated JSON. The tool
     **is** the schema.)
  2. **Two ways to get the model to actually call it:** *tell it* (mention it in the prompt
     / tool description — a soft nudge the model may or may not follow), or *force it* (set
     `tool_choice` — a hard guarantee the API enforces).
  3. **Don't force it up front.** If you force `final_answer` on turn 1, the model answers
     before doing any real work. So during the run leave `tool_choice: auto` and let the
     model use its real tools normally.
  4. **The loop watches every turn for a `final_answer` call.** This detection is the
     load-bearing part — on `auto` turns the model may decide it's done and call
     `final_answer` on its own, and your code must catch that, validate it, and stop.
  5. **If the model finishes without calling it** (a text-only turn — its "I'm done"
     signal) **or hits max turns**, do **one more turn** with `tool_choice` set to force
     `final_answer`, plus a short user nudge like *"Now return your final answer by calling
     the final_answer tool."* The wording barely matters (the forcing compels the call);
     the user message is mainly there to make the request well-formed, since the history
     currently ends on an assistant turn.
  6. **Mental model to keep:** *the tool is the schema carrier; detection is mandatory;
     forcing is optional insurance* you apply only at conclusion time.

  **Now the details** — three adjustments to the book's version, because its plain approach
  breaks or misbehaves on this Anthropic/thinking-enabled codebase:

  1. **Force the *named* tool, not `"required"`.** `tool_choice` has four modes; the book
     uses two of them under OpenAI/LiteLLM names, which map to the Anthropic SDK as:
     `auto` → `{type:"auto"}` (model decides, may emit text); `required` → `{type:"any"}`
     (must call *some* tool, model picks which); the *named* form
     `{type:"function",function:{name:"X"}}` → `{type:"tool", name:"X"}` (must call *that*
     tool); `none` → `{type:"none"}`. The book sets `required`/`any` for the whole run —
     that only forces *some* tool, so with real tools present the model can keep calling
     `search`/editor instead of finishing, and it forbids plain text turns. Instead force
     the **specific** tool (`{type:"tool", name:"final_answer"}`), and **only when you want
     the run to conclude — not on every turn.** Forcing the named tool from turn 1 makes
     the model answer before doing any work; the natural flow is `auto` during the run, then
     switch to forcing `final_answer` at conclusion time (e.g. on the last allowed turn, or
     once the model produces a turn with no tool call).
  2. **Disable forcing while extended thinking is active.** On Anthropic, `tool_choice`
     must be `auto` when extended thinking is enabled — so forcing `final_answer` together
     with `WithThinking`/`WithAdaptiveThinking` is an **API error**. When thinking is on,
     leave `tool_choice: auto` and rely on detecting + validating the `final_answer` call
     rather than forcing it. (This is the most important caveat for this repo.)
  3. **Keep validate-and-retry on the tool result.** Tool-calling is trained behaviour,
     not constrained decoding — the model can still emit a schema-invalid object. Validate
     the `final_answer` input against the schema and feed a validation error back as the
     tool result so the model retries, exactly as the book relies on Pydantic to do.

  Note on the schema shape: making the output type a *parameter* (`output: <Schema>`) nests
  it one level (`{"output": {...}}`); acceptable, but keep it in mind when defining the tool.

  > **Note on completion detection — the book vs. this repo.**
  >
  > *In one line:* the book's "am I done?" check is fragile because of *how it stores
  > events*, but this repo's `RunAgent` doesn't inherit that fragility — so for us
  > `final_answer` is a nice-to-have (typed output + a cleaner "done" signal), not a
  > bug fix. The rest of this note explains why.
  >
  > *The book's `_is_final_response` and why it checks tool results.* In the book's
  > `step()` (§4.6.3) a single step appends **two** events to the context: (1) always the
  > assistant event (its `Message` text and/or `ToolCall`s), and (2) *only if* tools were
  > called, a second **tool-result event** holding the `ToolResult`s. `run()` then checks
  > `_is_final_response(events[-1])`. Mid-loop, `events[-1]` is therefore the
  > **tool-result event**, which contains `ToolResult`s but **no** `ToolCall`s. If the
  > predicate checked only `has_tool_calls`, that event would look "final" and the loop
  > would **halt right after executing a tool, before feeding the result back** to the
  > model. Hence it must also check `has_tool_results`; together they encode "final iff a
  > pure message event — no tool calls *and* no tool results." That definition detects
  > completion by the **absence** of tool activity, which is why it is fragile in the
  > abstract: any message-only turn (e.g. chatty "let me now calculate…" narration) looks
  > final.
  >
  > *Why the current Go `RunAgent` does **not** share that fragility.* Two reasons:
  > (a) **Structural** — `RunAgent` decides termination from the assistant turn's
  > `len(toolUses) == 0` **before** appending any tool-result message, and never
  > re-inspects a tool-result event. So the `has_tool_results` failure mode is absent by
  > construction. (b) **API semantics** — on Anthropic, when Claude uses a tool the
  > preamble text and the `tool_use` block arrive in the **same** assistant turn
  > (`stop_reason = "tool_use"`); a text-*only* turn means `stop_reason = "end_turn"`, i.e.
  > genuinely done. `resp.Accumulate` reassembles that single turn's blocks (text +
  > tool_use) together, *preserving* this — but the guarantee comes from the turn model
  > plus the direct tool-use check, not from accumulation itself. (Residual: a model could
  > still end its turn early with a weak text-only non-answer — but that is a model-quality
  > issue the book's heuristic shares, not the book's specific halt-before-observing bug.)
  >
  > *Implication for the port.* For `RunAgent`, `final_answer` is therefore **not** a
  > termination bugfix — its value is (a) typed output and (b) an *explicit, positive*
  > completion signal (a successful `final_answer` result) instead of inference-by-absence,
  > which is a robustness nicety, not a correctness patch. When no output schema is
  > configured there is no `final_answer` tool and the existing `len(toolUses) == 0` check
  > remains the termination signal for free-text runs — which is sound here.

- **`Execute` (single-shot, no tools) is a different case.** Do **not** run this tool-loop
  pattern there — for one-shot typed output prefer a native response-schema path rather
  than spinning the agent loop. The tool-based approach earns its keep only in the
  multi-tool `RunAgent` loop.

- **Where `final_answer` lands in `RunAgent` — three spots.** `final_answer` is a *special*
  tool: it does no external work; its job is to carry the typed payload **and** end the
  loop. So it cannot ride the normal execute-then-continue path. It lands in three places
  (line refs are to the current `client.go`):

  1. **Definition → `params.Tools` (via the ToolSet).** Same place every tool's schema
     lands today: `buildParams` → `cfg.toolSet.Params()` (~`client.go:416-418`). A
     `WithOutputSchema(schema)` option synthesizes a `Tool` whose `Param()` carries
     `InputSchema = <output schema>` and registers it, so the model sees `final_answer`
     alongside the real tools. This option also drives the `tool_choice` decision
     (adjustment #1: force the *named* tool at conclusion time, honouring the thinking
     caveat in #2).

  2. **Handling → a branch in the tool-use loop, not the normal execute path.** Today the
     loop collects `toolUses`, executes *all* via `cfg.toolSet.Execute`, and — because
     `len(toolUses) != 0` — appends results and **continues** (~`client.go:161-218`). For
     `final_answer` that "continue" is exactly wrong. Intercept it before then:

     ```go
     for _, tu := range toolUses {
         if cfg.outputToolName != "" && tu.Name == cfg.outputToolName {
             val, err := decodeAndValidate(cfg.outputSchema, tu.Input)
             if err != nil {
                 // adjustment #3: feed the error back, DON'T terminate — model retries
                 toolResults = append(toolResults,
                     anthropic.NewToolResultBlock(tu.ID, "ERROR: "+err.Error(), true))
                 continue
             }
             return &AgentResponse{           // valid → this IS the end
                 Output: val,                  // typed result (new field)
                 Text:   joinTextBlocks(resp.Content),
                 Turns:  turn,
                 Usage:  totalUsage,
             }, nil
         }
         // ...normal tool → cfg.toolSet.Execute as today
     }
     ```

     So `final_answer` becomes a **second termination condition**, sitting beside the
     existing `len(toolUses) == 0` check: the loop ends either on a text-only turn
     (free-text mode) *or* on a successful `final_answer` call (structured mode).

  3. **Result surfacing → a new field on `AgentResponse`.** The parsed value has nowhere to
     go today (`AgentResponse` is `Text/Turns/Usage`). Add `Output any` on `AgentResponse`
     (or set `FinalResult` on the `ExecutionContext` from §4.1). `Text` keeps any
     accompanying prose; `Output` carries the schema-validated object.

  **What it does *not* need:** no real entry in `ToolSet`'s dispatch `switch` / no
  side-effecting executor — its "execution" is just validate-and-capture, done inline in
  the branch above. Even if it were routed through a generic `Tool.Execute`, the
  **termination decision must live in `RunAgent`**: `ToolSet.Execute` returns a `string`
  and cannot stop the loop. Architectural takeaway — the *schema* lands in the ToolSet, but
  the *control-flow* (terminate vs. retry) lands in `RunAgent`.

---

## Appendix A — book reference (distilled)

> Source: *"Implementing a basic ReAct agent"* (chapter 4), Manning liveBook —
> `https://mng.bz/X7MG`. This is a **condensed, paraphrased** reference of the
> architecture the analysis above cites — component roles and data-structure *shapes*
> only, not the book's prose. Consult the source for full explanation and code.

### A.1 ReAct in one line
Instead of answering all at once, the agent alternates **reason → act → observe** until
it has enough to answer. Modern implementations replace fragile text parsing with
**tool-calling**: the reasoning moves inside the LLM's tool-selection, so it's still ReAct.

### A.2 Component roadmap (book table 4.1)
| Order | Component | Role |
|---|---|---|
| 1 | `ExecutionContext` | Central store for all execution state. |
| 2 | Tool abstraction | Unifies local functions + MCP tools under one interface that can receive context. |
| 3 | LLM communication layer | `LlmRequest` selects what to send; `LlmClient` calls the API; `LlmResponse` standardizes replies. |
| 4 | `Agent` | Orchestrator: creates context, runs the think–act loop. |

### A.3 Content, Event, ExecutionContext (book §4.3) — shapes
- **ContentItem** = discriminated union by `type`:
  - `Message{ role: system|user|assistant, content: str }`
  - `ToolCall{ tool_call_id, name, arguments: dict }`
  - `ToolResult{ tool_call_id, name, status: success|error, content: list }`
  - (`tool_call_id` links a result back to its originating call — matters for parallel calls.)
- **Event** = content + metadata: `{ id, execution_id, timestamp, author, content: [ContentItem] }`
  (`execution_id` groups a run; `author` = "user" or agent name.)
- **ExecutionContext** = `{ execution_id, events: [Event], current_step: int,
  state: dict, final_result: str|Model|None }` + `add_event()`, `increment_step()`.
  `current_step` prevents infinite loops; `state` is a tool-accessible scratch-pad.

### A.4 Tool abstraction (book §4.4) — shapes
- **BaseTool**: `name`, `description`, `tool_definition` (schema for the LLM), and abstract
  `async execute(context, **kwargs)`. **Every** tool receives `context` — that's the
  context-propagation hook; tools that don't need it ignore it.
- **FunctionTool**: wraps a plain function; inspects the signature for a `context` param
  and forwards it only if present; auto-generates `tool_definition` from type hints.
- **`@tool`** decorator: sugar for `FunctionTool(func)`, with optional name/description override.
- **`load_mcp_tools(connection)`**: connects to an MCP server, lists its tools, wraps each
  in a `FunctionTool` (passing the server's `inputSchema` as `tool_definition` rather than
  generating one). The agent treats local and MCP tools identically.

### A.5 LLM communication layer (book §4.5) — shapes
- **LlmRequest** = `{ instructions: [str], contents: [ContentItem], tools: [BaseTool],
  tool_choice: auto|required|None }`. Deliberately holds **no** ExecutionContext — the
  agent curates a subset into it. *This is where context engineering happens.*
- **LlmResponse** = `{ content: [ContentItem], error_message: str|None, usage_metadata }`.
  API errors are captured, not raised, so the loop can handle them.
- **LlmClient**: `generate(request) -> LlmResponse`; builds provider messages, extracts
  tool definitions, calls the API, parses the reply back into ContentItems.

### A.6 Agent loop (book §4.6) — five methods
- `_setup_tools()` — register tools.
- `run(user_input, context=None)` — create/reuse context, add the user Event, loop
  `step()` until `final_result` or `max_steps`, return `AgentResult{output, context}`.
- `step()` — one think–act cycle: `_prepare_llm_request()` → `think()` → record Event →
  if tool calls, `act()` and record results → `increment_step()`.
- `think(request)` — delegate to `LlmClient.generate`.
- `act(context, tool_calls)` — dispatch each call by name, pass `context`, wrap failures as
  `ToolResult(status="error")` (error becomes an observation the model adapts to).
- Completion: `_is_final_response()` (no tool calls/results in the event) +
  `_extract_final_result()`.

### A.7 Structured output (book §4.7) — mechanism
Reuse tool-calling to enforce a typed result. Given an `output_type` schema:
- `_setup_tools()` dynamically creates a `final_answer` tool whose input **is** the schema.
- `_prepare_llm_request()` sets `tool_choice="required"` to force a tool call.
- `_is_final_response()` / `_extract_final_result()` treat a successful `final_answer`
  call as completion and return the validated typed object (not free text).
