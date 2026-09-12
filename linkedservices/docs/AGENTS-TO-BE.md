# AGENTS-TO-BE — LLM provider strategy for `agents/` + EINO evaluation

> Status: **Proposed** (ADR + roadmap) · 2026-09-03 · Scope: LLM execution for the
> `agents/` folder, driven by `ana_fun`; informs future agents.
>
> **Out of scope / do not build on:** the legacy `executor` package, the `ai_prompt`
> worker, the prompt `Category`/`Lks`/`HasBatchAPI` scheme, and the legacy
> `anthropiclks`/`ollamalks` `Client`/`BatchClient` interfaces. They still exist elsewhere
> but are not a foundation. **Ollama is treated as green field.**

## TL;DR — decision

**Get multi-provider (Anthropic · Ollama · vLLM) by adding symmetric provider _wrappers_ +
thin provider-agnostic _ports_ in `tpm-ai-common`, keeping the Anthropic Batch API
first-class** (bulk = batch on Anthropic, concurrency on Ollama/vLLM; stream on all). This
leaves each agent's workflow intact. **Do _not_ adopt an AI framework (EINO or any) for
`ana_fun` or for provider independence** — provider coverage is a wash (EINO reaches the same
set), frameworks own the whole workflow (rewrite + duplicated business logic) and don't model
async batch, and even the single tool loop is already hand-rolled (`RunAgent`). A framework is
a **separate, greenfield** decision reserved for **tool-agent _systems_** (cross-provider
tool-calling, multi-agent, MCP/A2A, shared observability). **Roadmap:** Step 1 — wrappers +
ports in `tpm-ai-common`; Step 2 — prove on a coexisting `ana_fun_ng` clone (incl. Ollama
end-to-end); Step 3 — gated broad adoption.

## 1. Why this document

`agents/ana_fun` is the least-trivial agent in `agents/`: it mixes a **deterministic**
step (partition — pure `depgraph` graph analysis, no LLM), an **async batch** step (synth —
one Anthropic Batch request per group), and an **online streaming** step (merge). That
surfaces two questions the four simpler batch agents never raise, and which need a recorded,
opinionated decision:

1. **Multiple provider implementations** for the LLM steps — be able to switch between
   **Anthropic**, **Ollama**, and **vLLM** (via its OpenAI-compatible API) providers.
2. Whether **EINO** (`github.com/cloudwego/eino`) makes sense here, and whether it is the
   right move for future agents. (§5 also tests the tempting idea that EINO could double as the
   answer to #1 — and rejects it.)

This is an ADR: it commits to a direction and a phased roadmap, not a neutral survey.

## 2. What `agents/` uses today (and nothing else)

| Consumer | Layer | Capability used | Streaming? |
|---|---|---|---|
| 4 batch agents (`copy_summary`, `node_summary`, `program_summary`, `program_context_info`) | NG `client` wrapper (`anthropiclks/client`, via `NewClientNG()` → `SubmitBatch`) | **bulk** (Anthropic Batch API) | no (not exposed to callers) |
| `ana_fun` synth | raw `anthropic.Client` SDK (`Messages.Batches.*`) | **bulk** (Anthropic Batch API) | — |
| `ana_fun` merge | raw `anthropic.Client` SDK (`Messages.NewStreaming`) | **stream** | yes |
| `agents_util` | — | none (logging + asset/store helpers only) | — |

Two capabilities are all `agents/` needs:

1. **Bulk generate** — N independent single-turn generations. Anthropic = Batch API.
2. **Stream generate** — one request, streamed tokens. Anthropic = `NewStreaming`.

Today, **bulk** is covered two different ways (the NG wrapper for the 4 agents; the raw SDK
for `ana_fun` synth), **stream** exists only inside `ana_fun` via the raw SDK, and **nothing
is provider-agnostic**. The NG `client` wrapper the batch agents use exposes `Execute`
(online single-shot) and the batch calls, but **no raw token streaming** (its `RunAgent`
streams internally only) — which is exactly why `ana_fun` bypasses it and talks to the raw
SDK.

## 3. The two problems

- **P1 — Consolidation.** `ana_fun` should stop hand-rolling the raw SDK and share one LLM
  abstraction with the batch agents. **Blocker:** the abstraction the batch agents use has
  **no streaming**, which `ana_fun`'s merge needs.
- **P2 — Multiple providers.** Add Ollama (green field) behind the same abstraction, with
  **batch as a capability**: Anthropic bulk = Batch API; Ollama bulk = bounded concurrency
  over single-shot generate; both provide streaming.

Both are solved by **one** small provider abstraction.

## 4. Proposed abstraction — two layers, both in `tpm-ai-common` (shared)

Put the whole abstraction in **`tpm-ai-common`** (next to the provider wrappers) so future
projects can reuse it — not in cob-game. Nothing from the legacy stack is carried over.

### Layer 1 — symmetric provider wrappers
Provider-specific, but the **same surface**: each wrapper exposes `Execute` (online
single-shot), `Stream` (streamed tokens), and `Bulk` (many at once).

- **Anthropic** — `anthropiclks/client` (exists): already has `Execute` + the batch calls
  (`SubmitBatch`/`GetBatch`/`CollectBatchResults`); **add `Stream`** to complete it. Its
  `Bulk` = the Batch API (async, ~50% discount, scale).
- **Ollama** — a **new sibling wrapper** (*not* the deprecated `ollamalks`): `Execute` +
  native `Stream`, and `Bulk` = a bounded-concurrency pool over `Execute`. Same method shape
  as Anthropic even though Ollama has **no batch API** — the concurrency policy lives
  **inside the wrapper**, so the wrappers share one shape.
- **OpenAI-compatible** — a wrapper against the OpenAI Chat Completions API with a configurable
  base URL. One adapter covers **vLLM** (point it at vLLM's `/v1`), OpenAI itself, and the wider
  OpenAI-compatible ecosystem (TGI, llama.cpp server, LM Studio, Groq, Together, OpenRouter, …).
  `Execute` + native `Stream`; `Bulk` = concurrency (no batch API) — and for **vLLM** that is the
  *right* pattern: its server-side continuous batching thrives under many concurrent requests.
  This is the highest-leverage wrapper to add — a single adapter unlocks many backends.

### Layer 2 — provider-agnostic ports
A new shared package (working name `tpm-ai-common/linkedservices/llm`):

```go
type BulkGenerator interface {
    Generate(ctx, []Request) ([]Result, error)   // → wrapper.Bulk (batch on Anthropic, concurrency on Ollama)
    Capabilities() Capabilities                    // Async, NativeBatch, …
}

type StreamGenerator interface {
    Stream(ctx, Request) (TokenStream, error)      // → wrapper.Stream
}

// Request/Result: system+user text, model, max-tokens, temperature; text+stop-reason+usage
```

Thin per-provider **adapters** bind a wrapper to the ports; since the wrappers all expose the
symmetric `Execute`/`Stream`/`Bulk`, the adapters mostly normalize provider types to
`Request`/`Result`.

**Naming.** Keep `BulkGenerator` — it says "many at once" **without** saying "Batch" (batch
is the Anthropic *mechanism*, not the abstraction). Keep `StreamGenerator` for symmetry
(`Streamer` is the crisp alternative). Port methods `Generate`/`Stream` follow conventional
naming (also the verbs EINO's `ChatModel` uses). An optional umbrella `Provider` interface can
embed both, while an agent still holds only the narrow port it needs (batch agents:
`BulkGenerator`; `ana_fun`: both).

**Selection & injection.** The driving worker picks a provider from **simple config**
(provider name + model) and **injects** the port(s) into the agent, exactly like the existing
`SetInputs`/`SetLogger`. The agent depends only on the ports — store-free, provider-free, and
(for `ana_fun`) free of the raw SDK. During Step 1 **only `ana_fun_ng` imports the new
package**, and the interface is marked **experimental** until proven.

**Outcome:** an agent's `SubmitBatch` block becomes a `BulkGenerator.Generate` call;
`ana_fun` uses `BulkGenerator` (synth) + `StreamGenerator` (merge); Anthropic and Ollama are
both pluggable — and the whole abstraction is reusable outside cob-game.

### Proposed package layout (naming)

Concrete names for the placeholders above, all under `tpm-ai-common/linkedservices/`, following
the existing `anthropiclks/client` convention (a `client` subpackage inside each `*lks` provider
dir). **Proposed** — adjust before creating the packages; this is what §7's "package names" open
item resolves to.

| Component | Proposed package | Status | Notes |
|---|---|---|---|
| Provider-agnostic **ports** (`BulkGenerator`/`StreamGenerator`, `Request`/`Result`, `Format`, adapters) | `linkedservices/llm` | new | Working name kept; the seam every agent imports. Per-provider adapters live here (or beside each wrapper). |
| **Anthropic** wrapper | `linkedservices/anthropiclks/client` | **exists** | Only change: **add `Stream`** (§4 Layer 1). `Bulk` = Batch API. |
| **Ollama** wrapper (new, green-field) | `linkedservices/ollamalks/client` | new subpkg | New `client` subpackage mirroring `anthropiclks/client`; the **deprecated top-level `ollamalks`** stays untouched (§ out-of-scope). `Execute`/`Stream`/`Bulk` = concurrency. |
| **OpenAI-compatible** wrapper (covers vLLM, OpenAI, TGI, …) | `linkedservices/openailks/client` | new | Configurable base URL → vLLM `/v1`, etc. `Execute`/`Stream`/`Bulk` = concurrency. Highest-leverage wrapper (§4 Layer 1). |
| Strategy docs + examples | `linkedservices/docs/` (`AGENTS-TO-BE.md`, `AGENTS-IO.md`, `examples/`) | **placed** | This doc + the I/O companion + per-sample example docs, in the `docs/` folder for this repo's AI stack. |

Naming choices still genuinely open (pick before coding): the ports package name (`llm` vs.
alternative); whether the new Ollama wrapper is a `client` subpackage of `ollamalks` or a fresh
dir (e.g. `ollamang`); the OpenAI-compatible dir name (`openailks` vs. `oaicompatlks`).

## 5. Does EINO make sense? Not for this — and the reason generalizes

Earlier drafts conflated two orthogonal axes. Separating them **is** the answer.

- **Axis A — provider abstraction** (*which* backend). This is what #1 asks for. Solved by the
  wrappers/ports (§4): additive, leaves the agent's workflow intact, **preserves the Batch
  API**, and extends to further backends (OpenAI, Gemini, …) by **writing more wrappers** on our
  own terms. EINO is **not needed** here, and is a worse tool for it (no batch).
- **Axis B — orchestration framework** (*how* the workflow is expressed and run). This is what
  EINO fundamentally is (`Graph`/`Chain`/`Workflow` + `ChatModel` + callbacks + interrupt/
  resume). Adopting it is an **Axis-B bet**, and Axis-B is **all-in for a given agent** — you
  cannot run half an agent under it.

**Why it can't sit *side by side* with the Batch API (correcting the earlier draft):**

1. **Framework = workflow ownership → rewrite, not migration.** Once EINO owns orchestration,
   `ana_fun` is not adapted — it is **re-implemented as a different agent** (an EINO graph), and
   the business logic (partition/ownership rules, token budgeting, reduce order, prompt
   assembly, asset stamping) is **duplicated** across the current agent and the EINO one, with
   the usual divergence risk. This is inherent, not an EINO quirk.
2. **The async Batch API fights the model.** EINO executes a graph synchronously (with stream
   plumbing); the only way to express "submit now, collect later" is to abuse
   **interrupt + checkpoint/resume** — a human-in-the-loop mechanism (submit → interrupt →
   resume on batch-done) that couples our cost optimization to EINO's checkpoint store, node-I/O
   serialization, and "re-run unfinished tasks on resume" semantics. So batch is either a **long
   blocking node** (defeats the framework) or an **interrupt hack** (bends it) — never a clean
   peer `Lambda` "side by side" with an online graph.
3. **Provider breadth is not a unique EINO advantage.** OpenAI/Gemini/etc. via `eino-ext` is
   replicable by writing wrappers — without the framework tax and without giving up batch. (You
   *could* use `eino-ext`'s `ChatModel` components as a provider layer while ignoring the graph,
   but that is a heavy dependency for the component library alone, still has no batch, and buys
   nothing over wrappers.)

**This generalizes to any AI orchestration framework** — LangChainGo, Genkit, LlamaIndex-style,
etc. Each imposes its own workflow model, so retrofitting an existing agent means rewrite +
duplicated business logic; and each is optimized for online/agentic/streaming flows and does
**not** model async batch (a provider-specific cost optimization, not an orchestration
primitive). The batch-oriented fan-out — where our cost model lives — is exactly where these
frameworks are weakest.

**A single tool loop is already hand-rolled — so "agentic" alone is not the trigger.** The
Anthropic wrapper's `RunAgent` is a working agentic loop (tool set, multi-turn, max-turns,
streaming), thin over the native SDK. A single-provider tool loop is therefore cheap and
*already exists* here, so the framework break-even is higher than "do you have tools" — it is
**tool-agent *systems***:

- **cross-provider tool-calling** — `RunAgent` is Anthropic-only; normalizing tool
  schemas/semantics across Claude / OpenAI / vLLM / Ollama (parallel calls, tool-choice, streaming
  tool-call deltas) is the fiddly, drift-prone part a maintained framework earns its keep on;
- **orchestration beyond one loop** — multi-agent, branching/conditional graphs, planner/executor,
  sub-graph composition, A2A;
- **ecosystem / protocol plumbing** — MCP tool servers, RAG / retrievers / memory;
- **a shared substrate at fleet scale** — uniform tracing, checkpoint/resume (HITL), retries.

None of these is a lone tool loop; hand-rolling our own versions of *them* is "building a worse
framework."

**Verdict.** Do **not** adopt EINO (or any framework) to solve #1, and do **not** retrofit
`ana_fun`. A framework is a **separate, forward-looking decision about new agent shapes** —
**tool-agent *systems*** and *RAG / multi-agent / human-in-the-loop* agents where orchestration,
cross-provider tool-calling, streaming, and interrupt/resume are the core value and there is
**no batch**. A single-provider tool loop does *not* qualify — `RunAgent` already covers it. For
those systems, EINO is a strong candidate, to be evaluated **standalone and greenfield**. The wrapper/port layer is
**framework-independent**: new EINO agents could still call our wrappers (or `eino-ext` models),
but that choice is decoupled from `ana_fun` and from #1.

- **Risks (if/when we go there):** young / fast-moving framework, ByteDance/CloudWeGo-driven
  (Apache-2.0, ~12.9k★), partly Chinese docs, dependency weight, generic-`Graph` learning curve,
  and adapters between EINO and our `agentregistry`/`AgentExecution`/asset/worker model.

### 5a. Even without batch, is EINO worth it for `ana_fun`? — No, overkill

Suppose we dropped the batch requirement entirely (synth becomes bounded-concurrency online
calls). Would rebuilding `ana_fun` on EINO then make sense? Still no — for `ana_fun`
*specifically* it is overkill:

- **`ana_fun` is a static, deterministic pipeline with LLM leaves** — partition (pure Go) →
  fan-out synth (N independent generations) → reduce/order (pure Go) → one streaming merge. No
  dynamic branching, tool loops, retrieval, multi-agent, or human-in-the-loop. Without batch,
  synth is just a **bounded-concurrency map** and merge is **one** streaming call.
- **EINO's value is in the parts `ana_fun` doesn't have** — dynamic/agentic control flow,
  stream-through plumbing, interrupt/resume. Orchestrating this fixed 4-node DAG is ~40 lines of
  plain Go (`errgroup` + a semaphore); a graph abstraction over it is heavier, not clearer.
- **Its one relevant benefit — provider independence — the wrappers/ports already deliver.** So
  EINO would add essentially nothing while still charging the framework tax: a rewrite (not a
  migration), business-logic duplication, dependency weight, the `Graph` learning curve, and
  impedance-matching EINO's types to our `AgentExecution`/asset/worker model. The partitioning
  IP just becomes an opaque `Lambda`.

**Where the calculus flips:** a framework amortizes across *many* agents, not one. A commitment
to a **portfolio** of LLM pipelines wanting a shared substrate (tracing/callbacks/checkpoint/
stream handling), or to genuinely **dynamic/agentic** agents, is what justifies EINO — a
strategic *platform* decision, at which point a rebuilt `ana_fun` might ride along for
consistency. For one static pipeline in isolation, it is not worth it; the batch constraint only
makes an already-weak retrofit case weaker.

### 5b. The Go framework field (2026) — EINO is not the only option

If/when the portfolio trigger above fires, EINO is one of **four** credible, maintained Go
frameworks — worth choosing among deliberately:

| Framework | Backing | Sweet spot |
|---|---|---|
| **Google ADK** (Agent Development Kit, Go) | Google (US) | Corporate-backed, Go reached **1.0 (Nov 2025)**; **strongest of the four for genuine multi-agent** (sequential/parallel/loop, agent-as-tool), native OpenTelemetry, MCP + A2A, Vertex deploy/eval. But **Gemini / Google-Cloud gravity** (non-Gemini provider maturity in the *Go* build is the open question), **no batch**, and younger in Go than in Python. **Ladder rung 3 (multi-agent) — §5b.** See [Appendix B](#appendix-b--google-adk-in-depth). |
| **Firebase Genkit** (Go) | Google (US) | Production **app-framework** for LLM features: typed **flows**, streaming, structured output, first-class **tracing/eval + Developer UI**, prompt management, provider-swap via **plugins**. **Front-runner for the online-agent case** (right-sized — *not* a multi-agent orchestrator, which suits us); batch N/A online; real caveat is Go plugins younger than JS (verify Anthropic/Ollama). **Ladder rung 2 (richer online) — §5b.** See [Appendix C](#appendix-c--firebase-genkit-in-depth). |
| **Eino** | ByteDance (CN) | Composable, **strongly-typed** `Graph`/`ChatModel`/tools/retrievers + agent layer; battle-tested at ByteDance scale; providers via `eino-ext`. **No batch** (and batch fights its sync graph — §5), rewrite-to-adopt (§5.1), CN-origin governance optics. **Ladder rung 3 (multi-agent) — §5b.** See [Appendix D](#appendix-d--eino-in-depth). |
| **LangChainGo** (`tmc/langchaingo`) | Community | Broadest provider/vector-store list; familiar if you know Python LangChain — but a community port (bus-factor: primarily one maintainer, no corporate backing) that lags Python on features and reflects the **pre-LCEL/LangGraph** architecture, with shallower integrations and no batch. Best for RAG breadth; **weakest for the cross-provider tool-calling that would justify a framework here.** See [Appendix A](#appendix-a--langchaingo-in-depth). |

Smaller/newer entrants exist (Jetify AI SDK, Anyi, Agent SDK Go, …) but are not yet in the same
maturity tier. And **"no framework"** is a legitimate choice for simple pipelines — the official
provider SDKs plus plain Go, which is exactly `ana_fun`'s situation.

Notes that matter for us:

- **Provider coverage is not a differentiator.** Eino covers the same set we care about —
  **Anthropic** and **Ollama** natively (`eino-ext`), and **vLLM** via its OpenAI `ChatModel` with
  a custom `BaseURL` at vLLM's `/v1`. But so do our wrappers (§4), where a single
  OpenAI-compatible adapter covers vLLM *and* the wider OpenAI-API ecosystem — and, unlike Eino's
  `ChatModel`, our wrappers **keep the Anthropic Batch API**. So the choice turns on
  orchestration + batch + simplicity, **not** on which providers are reachable.
- **Provenance is a *governance* question, not a *runtime* one.** Eino is Apache-2.0 and runs
  entirely on our infra (no phone-home), so its ByteDance/CloudWeGo origin is not a technical or
  security issue; it only affects roadmap control and procurement optics (some orgs restrict
  Chinese-origin dependencies). If that is a concern, **Google ADK or Genkit** are the
  Google-backed alternatives and **LangChainGo** the community one.
- **None of them models Anthropic's async Batch API.** ADK, Genkit, LangChainGo and Eino are all
  *online/agentic* orchestration frameworks — so the §5 argument ("a framework owns the workflow
  and does not model batch") holds across the entire field, not just Eino. The wrapper/port layer
  (§4) keeps us framework-agnostic underneath, so a future agent can pick any of these
  independently of `ana_fun` and of the Batch API.
- **For the *online-agent* trigger specifically, Genkit is the front-runner.** "Not a
  multi-agent orchestrator" is a **fit** for us, not a con (our use cases are simple, §5); batch
  is irrelevant online; and its DX/tracing/eval/RAG/structured-flows layer is exactly what
  hand-rolling would reinvent. Home-grown (`RunAgent` + §4 wrappers) still wins the **bare tool
  loop / provider-independence-only** case; ADK/Eino lead only if the trigger is genuine
  **multi-agent** systems. The one caveat to clear is Go-plugin maturity (verify Anthropic +
  Ollama). See [Appendix C](#appendix-c--firebase-genkit-in-depth).
- **Home-grown wins the bare-loop case — a framework is not the default.** For a single tool loop
  (single-provider, or multi-provider via the §4 wrappers) with no RAG / eval / multi-agent
  needs, `RunAgent` + the wrappers already cover it: no new dependency, native to our
  `agentregistry` / `AgentExecution` / worker model, nothing to integrate. A framework earns its
  place only when the agent needs **more** than that — the DX / observability / eval / RAG layer
  (→ Genkit) or genuine **multi-agent** orchestration (→ ADK / Eino). Absent those, prefer
  home-grown; "no framework" is a first-class option, not a fallback.
- **Decision ladder — which rung the trigger lands on** (summary of the above; the per-framework
  rows point to their rung):
  1. **bare tool loop / provider-swap only → home-grown** (`RunAgent` + §4 wrappers);
  2. **richer online** (DX / eval / RAG / structured output / cross-provider tool-calling) **→ Genkit**;
  3. **genuine multi-agent systems → ADK or Eino**;
  4. **batch → §4 wrappers** (no framework models it).

## 6. Decision & roadmap

**Decision:** build the shared provider abstraction in **`tpm-ai-common`** (symmetric provider
wrappers + provider-agnostic ports) and prove it on a **cloned** agent first (no disruption to
what works), add a green-field Ollama, and only then adopt broadly. **Keep the EINO question
entirely separate** — evaluate it on a **greenfield online/agentic agent** (never by
retrofitting `ana_fun` or any batch step), because a framework is an all-in Axis-B bet (§5).
**Batch stays a first-class capability throughout.**

- **Phase 0 — provider abstraction, proven on a coexisting clone.**
  *(split into three independently shippable steps so nothing that works is put at risk)*
  - **Step 1 — foundation in `tpm-ai-common` (no cob-game changes).** Add the new symmetric
    **Ollama wrapper** (`Execute`/`Stream`/`Bulk`, bulk = concurrency); **add `Stream` to the
    Anthropic wrapper** to complete it; define the `BulkGenerator`/`StreamGenerator` **ports +
    thin adapters** (working pkg `linkedservices/llm`), marked experimental. **Done when** the
    wrappers + adapters compile and pass unit tests (Anthropic bulk↔batch, Ollama
    bulk↔concurrency, both `Stream`) — **no consumer yet**, so cob-game is untouched.
  - **Step 2 — `ana_fun_ng` clone (cob-game).** **Clone `ana_fun` → a new `agents/ana_fun_ng`
    package** and wire *that* to the ports — leave the original `ana_fun` and the 4 batch
    agents untouched, so only `ana_fun_ng` consumes the new package. **Done when** `ana_fun_ng`
    (a) matches `ana_fun` on Anthropic and (b) runs a small program **end-to-end on Ollama** —
    with `ana_fun` still working, unchanged.
  - **Step 3 — major adoption (gated on Steps 1–2).** Only if it works well: migrate the 4
    batch agents' `SubmitBatch` blocks onto `BulkGenerator`, promote `ana_fun_ng` to replace
    `ana_fun`, and retire the raw SDK usage. The big, deliberately-deferred change.
- **Phase 1 — EINO spike (decoupled from `ana_fun`).** Evaluate EINO on the shape it fits:
  build **one small greenfield online/agentic agent** (e.g. a tool-loop or streaming Q&A over
  existing assets) on EINO `Graph`+`ChatModel` for Claude *and* Ollama. Do **not** port
  `ana_fun` or any batch step into it. Measure orchestration ergonomics, streaming/callbacks,
  interrupt/resume, and integration cost with our `AgentExecution`/asset/worker model.
  **Kill criteria:** integration cost or type friction outweighs the orchestration benefit for
  new agents.
- **Phase 2 — decide EINO's role for *new* agents.** Either adopt EINO as the orchestration
  layer for future online/agentic agents (which can still call our wrappers or `eino-ext`
  models), or keep hand-rolled agents. Either way `ana_fun` and the batch agents stay on the
  ports and the Batch API stays first-class — the two tracks never merge in one agent.

## 7. Non-goals / open questions

- **Non-goals:** touching the legacy `executor`/`ai_prompt`/`Category`/`ollamalks` (a separate
  deprecation, not a foundation); reworking the `Category`/prompt scheme; changing `ana_fun`
  or the 4 batch agents during Phase 0 Steps 1–2 (Step 1 touches only `tpm-ai-common`; the
  clone `ana_fun_ng` carries the experiment in Step 2); **retrofitting `ana_fun` or the batch
  agents onto an orchestration framework (EINO or other)** — frameworks are reserved for new
  online/agentic agents (§5), evaluated greenfield.
- **Decided:** the abstraction lives in **`tpm-ai-common`** (shared, reusable), in two layers —
  symmetric provider wrappers (Anthropic + new Ollama and OpenAI-compatible siblings — the last
  covering vLLM — each `Execute`/`Stream`/`Bulk`) and the provider-agnostic ports; Anthropic streaming is added **to the wrapper** (completing
  it); the **Ollama wrapper mirrors the Anthropic surface** with `Bulk` = internal concurrency;
  port names `BulkGenerator`/`StreamGenerator`, methods `Generate`/`Stream`; kept
  **experimental** (consumed only by `ana_fun_ng`) until Step 1 proves it.
- **Open:** exact package names in `tpm-ai-common` — see the **Proposed package layout** table
  in §4 for the current proposal (ports `linkedservices/llm`, new Ollama `ollamalks/client`,
  OpenAI-compatible `openailks/client`); provider-selection config shape (minimal — provider name
  + model, not `Category`); whether to expose the umbrella `Provider` interface; the retirement
  path for the original `ana_fun` after `ana_fun_ng` is promoted.

## Appendix A — LangChainGo in depth

Expands the §5b one-liner (what "trails the original" means). Does **not** change the §5 verdict (a framework owns the workflow
and models no batch — that holds for LangChainGo like the rest); this only sharpens the §5b
"if the portfolio trigger fires, which framework" choice.

### A.1 "Trails the original" = a governance/throughput fact, surfacing four ways
The root cause is who builds it: Python LangChain is backed by a funded company (LangChain,
Inc.) with a full-time team and commercial products (LangSmith tracing/eval, LangGraph) pulling
the OSS forward continuously. LangChainGo is a **community port**, historically driven mainly by
one maintainer (Travis Cline / `tmc`) plus contributors, with **no corporate sponsor** — a
fraction of the throughput. That shows up as:

1. **Feature lag (time).** Structured/JSON output, tool-calling refinements (parallel tools,
   `tool_choice`, streaming tool-call deltas), newer memory/retriever types, and new provider
   features land in Python first and reach Go later, partially, or never.
2. **Architectural lag (the important one).** Python LangChain re-architected: split into
   `langchain-core` + partner packages, made **LCEL** (the `|` pipe) the idiom, then moved
   agent-building to **LangGraph**, effectively deprecating the old
   `AgentExecutor`/`initialize_agent` model. LangChainGo largely reflects that **earlier
   generation** (Chains/Agents/Tools/Memory, pre-LCEL). So "the Python mental model transfers"
   is true — but it transfers the *previous* model; a team on current LangChain sees a prior era.
3. **Integration depth vs. breadth.** Its headline strength — a **broad** provider + vector-store
   list — is real, but breadth ≠ depth: individual integrations tend to be **shallower and less
   battle-tested** than their Python counterparts, some Python integrations have no Go analogue,
   and docs/examples are sparser (you end up reading Go source or translating Python docs).

### A.2 Why this specifically weakens it for *our* decision
§5 reserves a framework for **tool-agent systems**, chiefly **cross-provider tool-calling** —
exactly LangChainGo's soft spot:
- **Cross-provider tool-calling is where the maintenance lag bites hardest.** Normalizing tool
  schemas/semantics across Claude/OpenAI/vLLM/Ollama (parallel calls, tool-choice, streaming
  deltas) is fiddly, drift-prone work a *well-resourced* framework earns its keep on; a trailing
  community port is the least likely of the four to keep that surface current and correct.
- **No async Batch API** — same as the whole field, so the §5 core argument is unchanged.
- **Abstraction weight vs. our baseline.** We already hand-roll a clean single-provider tool
  loop (`RunAgent`). LangChain's abstractions (Chain/Agent/Memory/OutputParser) are heavy/leaky
  even in Python; the Go port takes on a **second-hand copy of an over-abstracted design**
  without the ecosystem depth that partly justifies it upstream — a worse trade than the
  leaner, Go-native ADK/Genkit/Eino.

### A.3 Where it genuinely wins
- **Max out-of-the-box breadth** for provider + vector-store coverage — attractive if the future
  direction is **RAG-heavy** apps wanting many backends with minimal glue.
- **Familiarity** for a team that already thinks in LangChain terms.
- **Provenance:** community, Apache-2.0 — no Chinese-origin procurement concern (the Eino caveat)
  and no single-vendor lock (the Google-backed options).

### A.4 Net
Among the four, LangChainGo is the **broadest-integration but weakest-maintained-for-
orchestration** option, and its one differentiator (breadth) is *not* what our trigger
(cross-provider tool-agent systems) needs — where **ADK/Genkit** (production-grade, actively
developed, 1.0) and **Eino** (composable, high-throughput) are stronger. Credible **only** if
RAG breadth + LangChain familiarity outweigh feature/architecture lag and shallow integrations;
for the systems a framework is reserved for here, it is the trailing choice.

## Appendix B — Google ADK in depth

Strengths/weaknesses of ADK (Agent Development Kit, Go) for our decision. As with Appendix A,
this does **not** change the §5 verdict (a framework owns the workflow and models no batch —
true for ADK too); it sharpens the §5b "which framework if the portfolio trigger fires" choice.

### B.1 What it is
Google's open-source (Apache-2.0) agent framework, part of a larger Google agent stack (Vertex
AI Agent Engine for managed deploy, the **A2A** agent-to-agent protocol Google co-authors, MCP
tool support). Launched first for **Python** (2025), with **Go** and Java following; the **Go**
build reached **1.0 in Nov 2025** — so it is corporate-backed and now API-stable, but the Go
line is **younger and less proven than Python ADK**. Its model is code-first and explicitly
**multi-agent**: `SequentialAgent` / `ParallelAgent` / `LoopAgent`, hierarchical composition,
and agent-as-tool, over a runner/session/event substrate.

### B.2 Strengths
- **Corporate backing + stability.** Full-time Google team, active development, an explicit 1.0
  commitment — the opposite of LangChainGo's bus-factor. Lowest "will this be maintained" risk
  of the four (alongside Genkit).
- **Best-in-class multi-agent orchestration.** Sequential/parallel/loop agents, hierarchies,
  agent-as-tool are **first-class primitives**, not something you assemble. If the trigger is a
  genuine multi-agent *system* (planner/executor, fan-out/fan-in of sub-agents), ADK is the
  strongest of the four.
- **Protocol/ecosystem, forward-looking.** Native **MCP** tools and **A2A** (Google is a
  driver of A2A) — the interop surface most likely to matter for cross-org / cross-framework
  agents.
- **Production/ops maturity.** Native **OpenTelemetry** tracing, an **evaluation** framework,
  and a real deploy path (Vertex AI Agent Engine, Cloud Run). It is built as a product-grade
  stack, not just a library.

### B.3 Weaknesses for *our* decision
- **Gemini / Google-Cloud gravity — the crux.** ADK is model-agnostic in principle and strongest
  on **Gemini + Vertex AI**. In Python, non-Google providers come via a LiteLLM bridge; **in the
  Go build, the breadth and maturity of non-Gemini providers (Anthropic, Ollama, vLLM) is the
  open question to verify** — likely thinner/newer than Gemini. Since our stated targets are
  **Anthropic + Ollama + vLLM with no Gemini emphasis**, this cuts against exactly our need.
  (Contrast §4's wrappers, where one OpenAI-compatible adapter covers vLLM + the wider ecosystem
  and we keep Anthropic Batch.)
- **No async Batch API** — same as the whole field; the §5 argument holds. ADK is
  online/streaming/multi-agent; batch is not an orchestration primitive it models.
- **Younger in Go; feature/example lag behind Python ADK.** Fewer integrations, examples, and
  community answers on the Go line specifically.
- **Weight + GCP gravity.** It is a sizable, opinionated stack (sessions, runners, events,
  artifacts, deploy). Its full value — managed sessions, Agent Engine deploy, eval — is realized
  **on Vertex AI / GCP**; on our own infra we use a subset while still carrying the abstractions.
- **Integration cost.** Adapting ADK's runner/session/event/artifact model to our
  `agentregistry` / `AgentExecution` / asset / worker model is real impedance-matching (the same
  caveat noted for Eino).

### B.4 Net
ADK is the **most production-mature, best multi-agent** option of the four, and — with Genkit —
the lowest maintenance risk. It is the right pick **if** the trigger is genuine multi-agent
systems and/or we are on or near **Gemini + Google Cloud**. Against *our* stated constraints
(Anthropic/Ollama/vLLM, batch first-class, own infra), its Gemini/GCP gravity, no batch, and
younger Go provider layer make it a strong-but-mismatched choice: excellent for the multi-agent
future the doc reserves a framework for, weaker for the *provider-independence-plus-batch* need
that §4's wrappers already solve.

## Appendix C — Firebase Genkit in depth

Strengths/weaknesses of Genkit (Go) for our decision. Same caveat as A/B: it does **not** change
the §5 verdict (framework owns the workflow, no batch); it sharpens the §5b choice.

### C.1 What it is
Google's open-source (Apache-2.0) framework for building **AI features in applications** —
positioned as a production app-framework, not a multi-agent orchestrator. Launched first for
**Node.js**, with **Go** (GA in 2025) and Python (beta) following, so its **center of gravity is
JS** and the Go line is newer. Core abstraction: typed, observable **flows** (functions with
input/output schemas) plus `generate()` with tools, retrievers/RAG, and embedders. Provider and
vector-store support comes via **plugins** (Google AI/Gemini, Vertex, plus OpenAI, Anthropic,
Ollama, etc.). A signature feature is the **Developer UI** — a local console to run flows,
inspect traces, and iterate on prompts (**Dotprompt**).

### C.2 Strengths
- **Best developer experience / observability of the four.** The Developer UI + native tracing +
  an **eval** framework + prompt management make iterating and debugging LLM features unusually
  pleasant. This is Genkit's real differentiator.
- **Config-over-code, provider-swappable.** Models and vector DBs are plugins you swap without
  rewriting flow logic — provider independence is a design goal (aligns with our §4 intent).
- **Typed flows / structured output first-class.** Schema-typed inputs/outputs and structured
  generation are core, not bolted on — good for reliable, testable LLM steps.
- **Google-backed + actively developed**, Apache-2.0; runs on any Go server (Cloud Run,
  Functions, self-host). Lower maintenance risk than LangChainGo.
- **Right-sized for "an LLM feature," not a heavy agent platform** — lighter and more focused
  than ADK/Eino when you don't need multi-agent graphs.

### C.3 Weaknesses and non-issues for *our* decision
Scored against *our* situation (mostly online, non-multi-agent), not a generic ideal — because
two things that look like cons on paper are neutral-to-positive here:

- **Not a multi-agent orchestrator — scope, not a weakness for us.** Genkit is flows + generate
  + tools + RAG, with no first-class sequential/parallel/loop **sub-agent** orchestration or A2A.
  But by §5's own logic our use cases are simple, so this is a **fit** (right-sized), *not* a
  fault — it would count against Genkit only if the trigger were genuine multi-agent systems
  (there, ADK/Eino fit better).
- **No async Batch API — irrelevant for the online case.** It matters only for the batch track
  (`ana_fun` synth), which §4's wrappers own; for online agents batch does not apply, so the
  strongest §5 anti-framework argument (batch fights the model) does **not** bite here.
- **Go line younger than JS — a bounded, verifiable risk (the real caveat).** Genkit's
  provider/vector **plugins are richest in Node.js**; the Go plugins for our backends (Anthropic,
  Ollama) exist but are **less battle-tested** — expect possible missing options, rough edges, or
  a plugin to pin/patch, plus sparser Go docs/examples. *Action:* a short spike to verify the
  Anthropic + Ollama **Go** plugins against our models largely retires this.
- **Real but small residual costs.** Integration to our `agentregistry`/`AgentExecution`/worker
  model; modest dependency weight and lock-in to Genkit's flow/runner shape; some Google gravity
  (less than ADK — Genkit runs anywhere).
- **Overlaps §4 for provider-independence *alone*.** If provider-swap is the only need, our
  wrappers already deliver it and Genkit is overkill; Genkit earns its place when you *also* want
  its DX/observability/eval/RAG/structured-flow layer.

### C.4 Net — a strong option for the online case
For **online** agents (batch set aside), Genkit is a legitimately strong choice and a real
alternative to home-grown — arguably the **best-fit framework of the four**, precisely because it
is *not* a heavy multi-agent orchestrator. The break-even against home-grown:

- **Baseline home-grown** = `RunAgent` (single-provider tool loop) + §4 wrappers (online provider
  independence) — already covers a **bare tool loop** and **provider-swap**.
- **Genkit adds**, that you would otherwise hand-build: tracing, an **eval** harness, **prompt
  management** (Dotprompt), **typed/structured flows**, **RAG** building blocks, and a
  **Developer UI** — hand-rolling those is "building a worse framework."

So: a new online agent needing **only a bare loop or just provider-swap** → home-grown + wrappers
win (cheaper, native, no new dep); a new online agent wanting the **richer layer** (observability,
eval, RAG, structured output, cross-provider tool-calling) → **Genkit is the stronger choice**. It
is *not* the tool for heavy multi-agent orchestration (ADK/Eino) or for the batch track (§4). The
one genuine caveat to clear first is the Go-plugin maturity spike above.

## Appendix D — Eino in depth

Consolidates the Eino case. The **full argument is in §5 and §5a** (why a framework can't sit
side-by-side with the Batch API; why Eino is overkill for `ana_fun` even without batch); this
appendix is the strengths/weaknesses summary parallel to A–C. It does **not** change the §5
verdict.

### D.1 What it is
ByteDance / CloudWeGo's open-source (Apache-2.0, ~12.9k★) Go framework: **composable components**
(`ChatModel`, tools, retrievers, embedders) + **orchestration** (`Graph`/`Chain`/`Workflow`) +
an **agent layer** (ReAct etc.) + **callbacks** and **interrupt/resume** (checkpoint). A
distinguishing trait is a **strongly-typed graph** — node input/output types are checked via Go
generics. Provider/component integrations live in **`eino-ext`** (Anthropic, Ollama, OpenAI, and
OpenAI-compatible via `BaseURL` for vLLM). Battle-tested at ByteDance production scale.

### D.2 Strengths
- **Composable, strongly-typed orchestration.** The generic-typed `Graph` catches node-wiring
  mistakes at compile time — the most robust orchestration model of the four for complex,
  branching flows.
- **Battle-tested at high throughput.** Proven in ByteDance production; resilience/scale are a
  design focus.
- **Full orchestration feature set.** Branching graphs, callbacks/observability hooks,
  interrupt/resume (HITL), agent patterns — the "real framework" capabilities.
- **Provider coverage via `eino-ext`** reaches the set we care about (Anthropic, Ollama, vLLM
  via `BaseURL`); Apache-2.0 and runs entirely on our infra (no phone-home).

### D.3 Weaknesses for *our* decision
- **No async Batch API — and worse, batch *fights* the model (§5.2).** Eino runs a graph
  synchronously; "submit now, collect later" can only be forced via interrupt + checkpoint/
  resume, coupling our cost optimization to Eino's checkpoint store and resume semantics. Batch
  is never a clean peer node. This is the sharpest anti-fit for `ana_fun`.
- **Framework ownership → rewrite + duplicated business logic (§5.1).** Retrofitting `ana_fun`
  re-implements it as an Eino graph and duplicates the partition/budget/reduce/asset logic.
- **Provenance is a *governance* question (not runtime).** ByteDance/CN origin is not a technical
  or security issue (Apache-2.0, self-hosted, no phone-home), but some orgs restrict Chinese-origin
  dependencies — a procurement/roadmap-control concern. If that blocks, ADK/Genkit (Google) or
  LangChainGo (community) are the alternatives.
- **Learning curve + churn.** The generic-`Graph` model has a real ramp; young/fast-moving; docs
  partly Chinese; plus the integration cost to our `agentregistry`/`AgentExecution`/asset/worker
  model.

### D.4 Net
Eino is the **composable, strongly-typed, battle-tested orchestrator** — the strongest of the
four where you want robust type-checked graph orchestration and high-throughput resilience for a
tool-agent *system*, and where CN provenance is acceptable. It is the **worst fit for our actual
first need** (`ana_fun`'s async batch fan-out), because batch fights its synchronous graph and
adoption means a rewrite. Reserve it, like the others, for a **greenfield online/agentic** agent
(§5 verdict), evaluated standalone.
