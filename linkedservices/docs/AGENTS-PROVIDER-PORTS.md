# AGENTS-PROVIDER-PORTS — shape of the provider-agnostic port layer

> Status: **Design note** (refines `AGENTS-TO-BE.md` §4) · 2026-09-14 · Scope: the
> concrete shape of the `linkedservices/llm` ports (`BulkGenerator` /
> `StreamGenerator`), how options are adapted across providers, and how the design
> maps onto `agents/promptagent`.
>
> No implementation exists yet — this is the design agreed before coding. It
> sharpens the "not much detail" left in `AGENTS-TO-BE.md` §4 (the two-layer
> abstraction) and resolves the open question there about how the async batch
> lifecycle and per-provider options are modelled.

## 1. Why this note

`AGENTS-TO-BE.md` §4 proposes two layers in `tpm-ai-common` — symmetric provider
wrappers and provider-agnostic ports — but sketches the ports as a single
`BulkGenerator.Generate([]Request) ([]Result, error)` plus a `Capabilities()`
placeholder. That signature hides the hardest fact about the providers: on
Anthropic, "Bulk" is the **async Batch API** (submit → poll → collect, resumable by
a durable batch id, fire-and-forget possible), while on Ollama/vLLM "Bulk" is
**synchronous bounded concurrency**. A single blocking method can't honestly carry
both. This note fixes the port shape so async is real when the provider offers it
and collapses to nothing when it doesn't.

## 2. Levels — three, not two

The port is **on top of** the existing provider wrapper, not inside it:

- **Layer 1 — provider wrapper** (`anthropiclks/client.Client`, and future
  `ollamalks/client`, `openailks/client`). Provider-specific through and through:
  Anthropic-shaped methods (`Execute`, `SubmitBatch`/`GetBatch`/`CollectBatchResults`)
  and Anthropic-semantic options (`options.go`). **Unchanged by this design** and
  unaware the ports exist.
- **Adapter** — a small type in each provider's corner (e.g. `anthropicBulk`) that
  *implements* the neutral port by calling its wrapper. Depends on Layer 1.
- **Layer 2 — provider-agnostic ports** (`linkedservices/llm`). Knows nothing about
  any provider. The only thing an agent imports.

So the dependency direction is: agent → port (`llm`) ← adapter → wrapper (`*lks/client`).
The wrapper never depends upward.

## 3. Port shape — a handle/job model, plus a blocking convenience

The Bulk port is three operations, not one, so an async provider can do real work
across the calls while a synchronous provider collapses them:

```go
type BulkGenerator interface {
    Submit(ctx context.Context, reqs []Request) (Handle, error) // start; return a handle
    Poll(ctx context.Context, h Handle) (Status, error)          // is it done?
    Collect(ctx context.Context, h Handle) ([]Result, error)     // fetch results
    Capabilities() Capabilities
}
```

- **Async provider (Anthropic):** `Submit` creates the batch and returns its
  **durable, serializable** id; `Poll` hits batch status; `Collect` streams results.
- **Synchronous provider (Ollama/vLLM):** `Submit` runs the bounded-concurrency pool
  and returns a handle that already holds the results; `Poll` always reports ready;
  `Collect` returns them. Handle is ephemeral.

### Where the synchronous work happens: eagerly in `Submit`

A synchronous provider has no background job, so the work must happen in one of the
three foreground calls. The design does it **eagerly in `Submit`**: `Submit` runs the
concurrency pool to completion and stashes the results (or error) in the `Handle`;
`Poll` then always reports ready; `Collect` is a **pure, idempotent read** of what
`Submit` already produced (it never re-runs the work). So results are available
immediately after `Submit` returns.

The reason to do the work in `Submit` rather than defer it into `Collect` is to keep
`Collect` a side-effect-free read on **both** kinds of provider — async `Collect`
fetches already-produced remote results, sync `Collect` returns already-produced
in-memory results — so caller code is identical and `Collect` is safe to call more
than once. The only difference the caller feels is *where the time goes*: on a sync
provider `Submit` blocks for the work and `Poll`/`Collect` are instant; on an async
provider `Submit` returns fast and the waiting is across `Poll`. The blocking
`Generate` convenience hides even that.

The one asymmetry this introduces: because a sync `Submit` runs eagerly, there is
nothing to cancel cheaply between `Submit` and `Collect` (unlike async, where the
work is remote and `Submit` was cheap). If cheap pre-`Collect` cancellation ever
mattered for a sync provider you would defer the work into `Collect` instead, at the
cost of `Collect` no longer being a pure read. For `promptagent` (which always
submits then collects) eager-in-`Submit` is the right default.

A blocking convenience `Generate(ctx, []Request) ([]Result, error)` is **not a
separate interface** — it is `Submit` + poll-until-ready + `Collect`, provided as a
free helper (or a default method) for callers that don't care about the async
lifecycle. There is **no singular `Generator`**: a single generation is `Bulk` with
a one-element slice, so a standalone non-streaming single-shot interface would only
restate a special case.

The `Handle` carries resumability: for a durable-batch provider it can be persisted
and rebuilt from a stored id (this is what powers `promptagent`'s `--resume-batch-id`
and fire-and-forget); for a synchronous provider it is ephemeral.

### Delivery axis is the only reason for a second interface

Cardinality (one vs. many) does **not** justify a separate type. The one axis that
does is **delivery shape**:

- `BulkGenerator` — one or many, sync or async, but always a **completed** result
  (`[]Request → []Result`). Async-vs-blocking lives *inside* it via the handle model.
- `StreamGenerator` — one request, **tokens over time** (`TokenStream`), which
  `[]Result` cannot express. Genuinely separate.

`promptagent` needs **only `BulkGenerator`** (both its modes — batch and online —
produce finished results, neither streams). `StreamGenerator` is `ana_fun`'s merge
step.

## 4. Options — a neutral **superset** Request, adapted downward

Two distinct option sets, with the adapter translating between them:

- **Neutral set** (`llm.Request`, in `linkedservices/llm`) — provider-agnostic,
  expressed in neutral terms.
- **Provider set** (`options.go` in each `*lks/client`) — rich, provider-shaped.
  **Unchanged.**

The neutral `Request` is a **superset** (union), not an intersection: it carries a
field for anything any provider might support, each provider **honours what it can
and ignores the rest** — the same "set it; the layer below honours it if it can"
behaviour `anthropiclks/client` already uses for temperature
(`modelSupportsTemperature` silently drops it on models that 400 on it). Crucially,
the superset fields stay **neutral types** (`Temperature *float64`,
`ReasoningEffort string`, `Schema map[string]any`), never Anthropic SDK types — so
provider types do **not** leak upward. Raising `options.go` as-is (Anthropic
`OutputConfigEffort`, `JSONOutputFormatParam`, model-quirk gating) into the neutral
layer is rejected: it would drag provider semantics up and force other adapters to
reverse-map Anthropic concepts. **Provider quirks stay down in the wrapper**
(`buildParams`, `modelSupportsTemperature`, `effortToBudget`), below the port.

### Field taxonomy — three buckets

Honor-or-ignore is only safe for one of these buckets:

1. **Essential inputs** — `Model`, `UserText`, `System`, `MaxTokens`. Not
   "honoured or ignored": they are the payload. Always used; dropping one is a
   broken request.
2. **Advisory knobs** — `Temperature`, `TopP`, `Seed`, `StopSequences`,
   `ReasoningEffort`, prompt caching. **Safe to drop**: ignoring shifts
   sampling/quality/cost, never the shape or validity of the result. The superset's
   silent honor-or-ignore applies here.
3. **Contract-changing capabilities** — silent omission makes the result
   structurally wrong. Must be honoured or fail loudly; never silently dropped.

The superset works **only** because bucket 3 is gated by `Capabilities()` (§5).

## 5. Capabilities() — gating the contract-changing fields

```go
// Capabilities advertises the contract-changing features a provider supports —
// the ones whose silent omission would make a result structurally wrong. Advisory
// knobs (temperature, top_p, seed, stop sequences, reasoning effort, prompt
// caching) are deliberately NOT here.
//
// Invariant: the zero value is "supports nothing", so a new adapter that forgets
// to set a flag fails safe (the agent errors instead of silently dropping).
type Capabilities struct {
    StructuredOutput bool // honours Request.Schema (JSON-schema-constrained output)
    ToolUse          bool // honours tool definitions / function calling
    Multimodal       bool // honours non-text input blocks (images, PDFs, documents)
    DurableBatch     bool // Bulk returns a resumable handle (fire-and-forget + resume by id)
}
```

The contract-changing set is **five features, not two**: structured output (schema),
tool use, multimodal input, durable/async batch, and streaming. Streaming is **not**
a struct flag — it is discovered by whether the value also satisfies
`StreamGenerator` (`sg, ok := gen.(StreamGenerator)`). So:

- the **interface set** (`BulkGenerator` / `StreamGenerator`) gates the *delivery
  shape*;
- `Capabilities()` gates the contract-changing *Request fields* (the four flags).

Two enforcement points, complementary:

- **Adapter is the safety net.** Each adapter's `Submit` honours a contract-changing
  field or returns an error — never silently drops. Example: `req.Schema != nil &&
  !caps.StructuredOutput` → error. A caller cannot accidentally get unconstrained
  text tagged as JSON.
- **`Capabilities()` is for pre-flight selection.** The worker/agent checks it to
  pick a provider, or to fail early with a clear message, before submitting.

The zero-value invariant is the quiet workhorse: a half-finished Ollama adapter that
has not wired structured output reports `StructuredOutput: false`, so a schema
request errors loudly instead of returning unconstrained text.

## 6. Applied to `promptagent`

Current `promptagent` reaches for `anthropiclks/client` directly and hand-branches
online vs. batch. Under this design:

- The agent holds a **`BulkGenerator`** port (only that). `buildOptions` stops
  returning `[]client.Option` and instead builds a neutral `llm.Request{Model,
  System, UserText, MaxTokens, Temperature, Schema}`. The Anthropic adapter lowers
  that into the wrapper's `With*` options; all Anthropic-specific handling
  (temperature gating, thinking form, `schema → output_config.format`) stays inside
  the wrapper, untouched.
- **Anthropic exposes two `BulkGenerator` implementations** — a synchronous one
  (`Submit` = the N `Execute` calls, handle instantly ready) and an async-batch one
  (`Submit` = real Batch submit, durable handle). The existing `--mode online|batch`
  flag is exactly the worker choosing which implementation to inject; the agent does
  not know which it got.
- **`BatchExecutionHint` maps onto the handle lifecycle:** empty `BatchId` =
  `Submit`; a supplied `BatchId` = rebuild a `Handle` and `Collect` (resume); no
  hint / online = blocking `Generate` on the synchronous implementation;
  `poll-interval`/`max-iterations` are the `Poll` cadence.
- **`--resume-batch-id` / fire-and-forget become a capability check:** guard on
  `gen.Capabilities().DurableBatch` before allowing resume, instead of assuming the
  Batch API is present.
- `buildExecutionResponse` and `writeResponses` **do not move** — they already
  operate on provider-neutral text/JSON, above the port.

## 7. Open questions

- Exact `Request` / `Result` field lists (which advisory knobs to include in the
  superset now vs. later).
- `Status` shape returned by `Poll` (ready/working/failed, plus per-request partials?).
- `Handle` serialization contract for `DurableBatch` providers (what the caller
  stores to resume — just the id, or an opaque blob).
- Whether an umbrella `Provider` interface embeds both ports, or agents always hold
  the narrow port(s) they need (§4 of `AGENTS-TO-BE.md` leans to the latter).
- Package name (`llm`) and adapter placement (in `llm` vs. beside each wrapper).

## Related

- `AGENTS-TO-BE.md` — the ADR this refines (§4 Layer 1/Layer 2, the roadmap, the
  EINO decision).
- `AGENTS-IO.md` — the I/O companion.
