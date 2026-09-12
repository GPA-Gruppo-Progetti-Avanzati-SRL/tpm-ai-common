# AGENTS-IO — prompt contracts & structured output across providers (Anthropic · Ollama · vLLM)

> Status: **Proposed** (design note + decisions) · 2026-09-12 · Scope: the **I/O contract**
> between our prompts and the LLM providers — how we mark up **inputs** (XML tags) and how we
> get **machine-parseable outputs** (XML tags vs. JSON Schema enforcement) across **Anthropic**,
> **Ollama**, and **vLLM** (OpenAI-compatible).
>
> **Companion to [`AGENTS-TO-BE.md`](./AGENTS-TO-BE.md).** That doc decides *how we reach*
> providers (symmetric wrappers + `BulkGenerator`/`StreamGenerator` ports, batch first-class,
> no framework). This doc decides *what we send and receive*. It **extends** the §4 ports with a
> structured-output capability and does **not** revisit the batch/framework verdicts.

## TL;DR — decisions

1. **XML tags as *input* markers stay.** `<cobol_source>…</cobol_source>` is current best
   practice (Anthropic leans hardest on it; portable as a plain delimiter everywhere). Nothing
   deprecates it; structured output is an *output* feature and does not touch input markup.
2. **XML tags as *output* markers are a convention, not an enforced contract.** For output a
   program parses, prefer **native JSON-Schema enforcement** where the output is genuinely
   **data** (fields / enums / numbers). Keep XML output tags where the output is a **document**
   (markdown, Mermaid, long free text) — schema enforcement buys little there and is awkward.
3. **JSON Schema is the cross-provider lingua franca.** Anthropic, Ollama, and OpenAI/vLLM have
   all converged on "give me an object matching this JSON Schema." That makes a **single
   provider-agnostic structured-output port** feasible: caller supplies a schema, adapter maps
   it to each provider's native param, port **re-validates** the returned JSON client-side.
4. **Add an optional `Format` field to the §4 `Request`** (nil = today's free-form text). No new
   port interface; structured output is a property of a generation request and composes with
   `Generate` (bulk) and `Execute`.
5. **Enforce hard limits client-side.** Only `type` / `enum` / `required` are reliably enforced
   across all three. Length / numeric / `pattern` constraints are **best-effort** and differ by
   provider and version — keep them as prompt instructions *and* validate after parsing.
6. **Transport (concurrency-only on-prem scenario): champion-native.** Keep **Anthropic native**
   (fidelity, prompt caching, and the Batch API of `AGENTS-TO-BE`); use **LiteLLM** only as the
   swap layer / gateway in front of the **OpenAI-compatible** on-prem backends (Ollama, vLLM).
   Structured output maps cleanly either way, because it travels as a JSON Schema.
7. **EINO buys nothing for this.** Structured output is a per-request field, not orchestration;
   its only relevant offer (provider abstraction) is already covered by the §4 ports / LiteLLM.
   Verdict unchanged from `AGENTS-TO-BE` §5.

---

## 1. The two jobs XML tags do — keep them separate

Our prompts (e.g. `cob-game/linkedservices/prompts_repo/copy-summary-prompt.txt`,
`cob-game/linkedservices/prompts_repo/node-summary-prompt.txt`) use XML tags
for **two unrelated jobs**. Judging them together is the mistake; the JSON-Schema feature only
touches the second.

| Job | Example | What it is | Affected by JSON Schema? |
|---|---|---|---|
| **Input markers** | `<cobol_source>{{.COBOL_SOURCE}}</cobol_source>` | Delimit / label parts of the prompt we *send in* | **No** |
| **Output markers** | model must emit `<summary>…</summary>`, `<type-of-copy>…</type-of-copy>` | A convention for semi-structured output we *parse ourselves* | **Yes — this is what schema enforcement replaces** |

## 2. Input markers — current, portable, keep

XML input tags are **not deprecated** and remain the recommended way to structure Claude prompts;
Claude was tuned to attend to them. They also work as plain delimiters on Ollama/OpenAI models
(those docs often suggest markdown headers or `###`, but XML delimiters are fine). Structured
output changes nothing here. **Decision: keep all input-side XML markup as-is.**

## 3. Output: XML tags vs. JSON Schema — when each

The distinction is **enforcement**:

- **XML output tags** — a *prompting convention*. Reliable-ish, but **not guaranteed**: a tag can
  be dropped or malformed, or stray text can appear outside the tags. We carry the parsing +
  validation. (This is exactly why the prompts spend several lines on "no code fences / saved
  verbatim to a file" — that ceremony exists only to make text extraction survive.)
- **JSON-Schema structured output** — an **API-enforced** guarantee via constrained decoding: the
  response is *made* to validate against the schema.

Rule of thumb — **three buckets** (worked examples in §6):

- **Output is data** (fields, an enum, numbers, arrays) → **JSON Schema**. The presence of every
  field and the exact enum tokens become guarantees instead of hopes. → e.g. `copy-summary`.
- **Output is a document** (markdown, Mermaid, long free text) → **keep XML tags**. Stuffing a big
  markdown/Mermaid blob into JSON string fields means escaping and awkward handling for no real
  gain; the only win (guaranteed presence of each part) rarely justifies it. → e.g. `node-summary`.
- **Output is a document whose fields are all data** (a human-readable block that is really a
  *rendering* of extracted values) → **model emits JSON (schema-enforced); code renders the
  document; constants/fixed text injected at render, never asked of the model.** → e.g.
  `program-java-conversion`.

These three buckets are the operative decision for classifying every prompt in `prompts_repo`.

## 4. JSON-Schema support across the three providers

All three enforce a JSON Schema; the **param differs**, and each enforces a **subset**.

| Provider | Parameter (wire) | Enforced? | Backend / notes |
|---|---|---|---|
| **Anthropic** | `output_config: { format: { type: "json_schema", schema: {…} } }` on `messages.create`; helper `client.messages.parse()`. Also `strict: true` on a tool def. | Yes | Validated/constrained. **Incompatible with citations.** Old top-level `output_format` is **deprecated** — use `output_config.format`. |
| **Ollama** | `format: {…json schema…}` on `/api/chat` \| `/api/generate` (or `format: "json"` for schemaless JSON mode) | Yes | Grammar-constrained (llama.cpp). Quality tracks the underlying model. |
| **OpenAI** | `response_format: { type: "json_schema", json_schema: { name, schema, strict: true } }` (or `{ type: "json_object" }`) | Yes | Requires `additionalProperties:false` and all keys `required` (or documented workarounds). |
| **vLLM** (OpenAI-compatible — what we reach *through* it) | Same `response_format` **plus** vLLM extras `guided_json` / `guided_grammar` / `guided_regex` / `guided_choice` via `extra_body` | Yes | outlines / xgrammar / lm-format-enforcer. Very complex schemas may be partially supported; `guided_json` is the escape hatch when `response_format` coverage is thin. |

**Two caveats that shape the design:**

1. **Each enforces a *subset* of JSON Schema, and the subsets differ.** Deep recursion, some
   `pattern`/`format` keywords, and unusual constructs may be rejected or silently ignored. The
   **portable core** is `type`, `enum`, `required`, nested `object`/`array`, `additionalProperties:false`.
2. **Length / numeric / pattern constraints are the weak part.** `maxLength`, `minimum`, `pattern`
   have uneven support across providers *and* versions. Treat them as **best-effort**: express
   intent in the schema, repeat as a prompt instruction, and **validate client-side after parsing.**

Enforcement guarantees *shape*, never *content correctness*.

## 5. The provider-agnostic structured-output port

### 5.1 Design — extend `Request`, don't add an interface

Structured output is a property of a generation request, so it rides on the existing §4
`Request`/`Result` rather than a new port. `nil` `Format` = today's free-form behavior.

```go
// package linkedservices/llm  (extends the AGENTS-TO-BE §4 ports)

// Format is an optional output contract on a Request.
// Nil  → free-form text (current behavior).
// Set  → the provider is asked to emit a single JSON object matching Schema,
//        and the port re-validates the result against Schema before returning.
type Format struct {
    Name   string          // schema name; required by OpenAI/vLLM, ignored by Anthropic/Ollama
    Schema json.RawMessage // a JSON Schema (portable subset — see §4) for the desired object
    Strict bool            // request hard enforcement where the provider supports it
}

type Request struct {
    System      string
    User        string
    Model       string
    MaxTokens   int
    Temperature float64
    Format      *Format   // NEW — nil = free-form text
}

// Result.Text still carries the payload; when Format != nil it is guaranteed
// to be a JSON object that has passed client-side schema validation.
// Optional convenience:
func Unmarshal[T any](r Result) (T, error) // json.Unmarshal(r.Text) once, typed

// Capabilities gains one flag so callers/adapters know enforcement strength.
type Capabilities struct {
    Async                  bool
    NativeBatch            bool
    NativeSchemaEnforcement bool // true: constrained decoding; false: best-effort + validate
}
```

The port contract adds one rule: **when `Format != nil`, the adapter always re-validates the
returned JSON against `Schema` client-side** (defense in depth — Ollama/vLLM enforcement quality
varies) and returns an error on mismatch. Enforcement is the provider's best effort; **validation
is ours, and non-negotiable.**

### 5.2 Adapter mapping — one schema, three params

Each per-provider adapter translates the same `Format.Schema` to that provider's native field.
(Wire shapes below; verify exact Go SDK symbols against the installed SDKs — the Anthropic wrapper
already builds `anthropic.OutputConfigParam` at `anthropiclks/client/client.go:426` for effort, so
`OutputConfig` is the seam to extend.)

```go
// Anthropic adapter — output_config.format
params.OutputConfig = anthropic.OutputConfigParam{
    // …existing Effort…
    Format: /* json_schema variant */ { Type: "json_schema", Schema: fmt.Schema },
}
// or use client.messages.parse() with the schema.

// Ollama adapter — the `format` request field takes the schema object directly
body["format"] = json.RawMessage(fmt.Schema)   // NOT the string "json"

// OpenAI / vLLM adapter — response_format
body["response_format"] = map[string]any{
    "type": "json_schema",
    "json_schema": map[string]any{
        "name":   fmt.Name,        // required here
        "schema": json.RawMessage(fmt.Schema),
        "strict": fmt.Strict,
    },
}
// vLLM fallback for schemas response_format won't take:
//   extra_body: { "guided_json": <schema> }
```

Because the schema is one artifact and the mapping is mechanical, adding a provider is "write the
one mapping," exactly like the §4 wrappers. This is the concrete payoff of §3's convergence claim:
**the XML-tag-and-parse convention needs bespoke extraction per prompt; a JSON Schema ports.**

### 5.3 Streaming

Structured output pairs naturally with `Execute`/`Generate` (bulk). For `Stream`, deliver the
**final validated object** at end-of-stream; partial-object streaming (rendering half-built JSON)
is possible on all three but is **out of scope** for the port's first version.

## 6. Examples

Each sample has its own doc under [`examples/`](./examples/) with the full
**Before / Next / Reasoning**. This table is only a pointer + the *type* of each example — add a
row per new sample and keep the synthesis here to one line, detail in the example doc.

| Example | Type (§3 bucket) | Doc |
|---|---|---|
| `copy-summary` | data → JSON Schema | [copy-summary.md](./examples/copy-summary.md) |
| `node-summary` | document → keep XML tags (counter-example) | [node-summary.md](./examples/node-summary.md) |
| `program-java-conversion` | document-as-rendering-of-data → JSON + deterministic render | [program-java-conversion.md](./examples/program-java-conversion.md) |

## 7. Transport sidebar — reaching the providers (concurrency-only on-prem)

Recorded here because it came up alongside the I/O question; it **refines**, does not replace,
`AGENTS-TO-BE` §4 (which keeps the Anthropic Batch API first-class). Scope: the *forced on-prem*
scenario where **concurrency-only is acceptable** (batch dropped).

- **"One-stop shop" (everything through LiteLLM's single OpenAI endpoint) does not save a
  wrapper.** You need an OpenAI client for LiteLLM regardless; the only question is whether
  Anthropic *also* rides through the proxy. Same client count either way — so the trade is
  **fidelity + directness on the champion path**, not code volume.
- **Recommendation: champion-native.** Keep **Anthropic native** (prompt caching — a real cost
  lever that flattens through the OpenAI surface; hot-path directness; thinking/tool-use fidelity;
  and the Batch API of §4). Use **LiteLLM only in front of the OpenAI-compatible on-prem backends**
  (Ollama, vLLM), where config-driven swap + gateway features (routing, spend, failover) are the
  point and provider-specific richness doesn't matter.
- **Flip to full one-stop only for gateway features** — centralized cross-provider *failover*,
  unified spend/rate-limit/key management, one observability pane — accepting the
  lowest-common-denominator flattening as the price. (LiteLLM can also expose an Anthropic-format
  `/v1/messages` passthrough to avoid flattening, but then you're running two formats through the
  proxy, not "a single OpenAI endpoint.")
- **Structured output is transport-agnostic.** It travels as a JSON Schema mapped to
  `output_config.format` / `format` / `response_format` (§5.2) — so this transport choice does not
  change the port design.

## 8. Embeddings & similarity — orthogonal capability, noted for completeness

Came up in discussion; **separate axis** from generation/output-format, recorded so the on-prem
picture is complete.

- **Anthropic has no embeddings endpoint.** Its API is `POST /v1/messages` only; the recommended
  embeddings provider is **Voyage AI** (separate SDK/key). So the champion for *generation* is not
  a source for *embeddings*.
- **Ollama and vLLM both do embeddings locally.** Ollama: `POST /api/embed` with an embedding model
  (`nomic-embed-text`, `qwen3-embedding:0.6b`, …). vLLM: serve an embedding model (`--task embed`),
  OpenAI-compatible `POST /v1/embeddings`.
- **Similarity is not a framework need.** Comparing two vectors is a cosine one-liner. Semantic
  *search* over a corpus is a vector-DB job — EINO's `Retriever`/`Indexer` components wrap
  Milvus/Qdrant/Redis/etc. but the search runs in the DB, not in EINO. For our current needs this
  is out of scope; if a RAG/search need appears, it's a distinct decision (and a distinct §5b
  framework trigger in `AGENTS-TO-BE`).
- **Dimension caveat:** vectors from different models aren't comparable — switching the embedding
  model means re-embedding the corpus.

## 9. Open questions

- **Schema authoring & storage.** Where do schemas live — a `schema.json` beside each prompt in
  `prompts_repo`, or inline in the agent? (Leaning: a file next to the prompt, versioned together.)
- **Validation library.** Which Go JSON-Schema validator for the client-side re-validate step
  (§5.1)? Pick one that supports draft 2020-12 and the portable subset.
- **`Strict` semantics per provider.** Confirm the exact Anthropic `OutputConfigParam` format
  union and OpenAI `strict` subset rules against the installed SDKs before finalizing the adapters.
- **Classify every `prompts_repo` prompt into a §3 bucket** (data → schema; document → keep XML
  tags; document-as-rendering-of-data → JSON + render) and record the migration list; add an
  example doc under `examples/` for each one worth working through.
- **Effort × schema interaction.** Verify enforced structured output composes with adaptive
  thinking / effort on Anthropic (it should; confirm no 400 with `output_config.format` + effort).
