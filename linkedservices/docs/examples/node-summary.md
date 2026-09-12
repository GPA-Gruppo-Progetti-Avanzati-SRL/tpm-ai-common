# Example — `node-summary` → stays on XML tags

> **Type:** document-shaped output (markdown overview + raw Mermaid + summary) → **keep XML tags**.
> **Source:** `cob-game/linkedservices/prompts_repo/node-summary-prompt.txt`
> **Pattern (AGENTS-IO §3):** *document → XML tags.* This is the **counter-example** — the case
> where JSON Schema is the wrong tool.

`node-summary` produces three outputs whose *content is documents*, not data: a plain-text Italian
summary, a full markdown functional overview (headers, lists, tables), and **raw Mermaid** flowchart
source. There are no enums, no typed fields — just free-form text of three different shapes.

## Before

The current prompt (excerpt — full text in the source file). It emits three XML-delimited outputs
and spends real effort on formatting guidance so the extracted text lands clean in a file:

```text
<cobol_source>
{{ .COBOL_SOURCE }}
</cobol_source>
...
Place this analysis inside <scratchpad> tags.
...
Place your summary inside <summary> tags.
...
Place your functional overview inside <overview> tags.
...
Output raw Mermaid source only. Do NOT wrap the diagram in a markdown code fence ...
Place your Mermaid diagram inside <flowchart> tags.
...
Each output tag must contain ONLY the raw content of its declared type, with no surrounding
markdown code fences:
- <summary> contains plain text only.
- <overview> contains a markdown document (headers, lists, tables). Markdown code fences are
  allowed here ONLY to delimit real code samples inside the document ...
- <flowchart> contains raw Mermaid source only, starting directly with `flowchart TD`.
```

## Next

**Unchanged — keep the XML tags.** No schema migration.

The only refinement worth considering is cosmetic: the `<cobol_source>` input marker stays (§2), and
the three output tags stay. If anything, the "no code fences / saved verbatim" guidance is *more*
justified here than in a data-shaped prompt, because the payload really is raw markdown/Mermaid that
a stray fence would corrupt.

## Reasoning

Why JSON Schema is the **wrong** tool here:

- **The outputs are documents, not data.** `overview` is a whole markdown document; `flowchart` is
  raw Mermaid source. Forcing these into JSON string fields means escaping newlines, backticks, and
  Mermaid punctuation into one long string — awkward to produce, awkward to consume, and slightly
  worse for the model.
- **The one benefit doesn't justify the cost.** The only thing a schema would guarantee here is
  *presence of all three parts* (an XML tag can silently go missing). That is real but small, and it
  doesn't outweigh the escaping tax and the loss of clean, directly-usable markdown/Mermaid files.
- **Enforcement buys nothing on free text.** `type`/`enum`/`required` — the reliably-enforced part
  of JSON Schema (§4) — has nothing to bite on when every field is unconstrained prose.

**Rule confirmed:** *document → XML tags.* This example marks the boundary that
[`copy-summary`](./copy-summary.md) (data → JSON) sits on the other side of, and that
[`program-java-conversion`](./program-java-conversion.md) (document *rendered from* data) resolves a
third way.
