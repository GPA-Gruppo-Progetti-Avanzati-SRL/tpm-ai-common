# Example — `copy-summary` → JSON-Schema enforced

> **Type:** data-shaped output (a summary string + a classification enum) → **JSON Schema**.
> **Source:** `cob-game/linkedservices/prompts_repo/copy-summary-prompt.txt`
> **Pattern (AGENTS-IO §3):** *data → JSON Schema.*

`copy-summary` is the cleanest data-shaped example: two outputs, both structured — an Italian
summary (≤500 chars) and a **classification enum**. This is where schema enforcement pays off most.

## Before

The current prompt verbatim. Note how much of it is *output-tag ceremony* (marked inline): the two
"place inside `<tag>`" lines, then a whole block on "plain text only / no code fences / no language
annotation / saved verbatim to a file." That text exists solely to make string-extraction of
`<summary>`/`<type-of-copy>` survive.

```text
You will be analyzing a CICS COBOL COPY and providing
1. a brief summary of the content of the COPY
2. the classification of the type of COPY: it can be a
     2.1 DATA DIVISION COPY
     2.2 PROCEDURE DIVISION COPY
     2.3 a BMS physical-map

Here is the COBOL source code you need to analyze:

<cobol_source>
{{.COBOL_SOURCE}}
</cobol_source>

**Output 1: Summary**
Write a concise summary of what the program does. This summary must be:
- Maximum 500 characters in length
- Written in Italian
- Focused on the primary business function or technical purpose
- Clear and understandable to both technical and business audiences
Place your summary inside <summary> tags.                        ← output-tag ceremony

**Output 2: COPY classification**
- Classifiy the COPY as DATA_DIVISION, PROCEDURE_DIVISION, BMP_MAP
Place your classification inside <type-of-copy> tag.             ← output-tag ceremony

**Important Formatting Requirements:**                            ← output-tag ceremony (whole block)

Each output tag must contain ONLY plain text, with no surrounding markdown code fences (no ``` delimiters) and no language annotation:
- <summary> contains the plain-text Italian summary only.
- <type-of-copy> contains exactly one of DATA_DIVISION, PROCEDURE_DIVISION, BMP_MAP and nothing else.
The text between the opening and closing tag is saved verbatim to a file, so a stray code fence corrupts that file.
```

## Next

### The schema (portable subset)

```json
{
  "type": "object",
  "additionalProperties": false,
  "required": ["summary", "type_of_copy"],
  "properties": {
    "summary": {
      "type": "string",
      "maxLength": 500,
      "description": "Concise Italian summary of the COPY's primary business/technical purpose, for both technical and business readers."
    },
    "type_of_copy": {
      "type": "string",
      "enum": ["DATA_DIVISION", "PROCEDURE_DIVISION", "BMP_MAP"],
      "description": "Classification of the COPY."
    }
  }
}
```

> Note: the enum keeps the **current** tokens (`DATA_DIVISION`, `PROCEDURE_DIVISION`, `BMP_MAP`).
> The prompt intro says "a BMS physical-map" but the emitted token is `BMP_MAP` — likely a
> `BMS`→`BMP` slip. Fix it if desired, but change it in **one place** (the schema) now.

### The revised prompt

The schema carries the structure, so every "output-tag ceremony" line is deleted. What remains is
only the genuine *content* instructions (Italian, ≤500 chars, the two things to produce):

```text
System:
You analyze CICS COBOL COPY members. Return your result ONLY as a JSON object
matching the provided schema. Write `summary` in Italian (max 500 characters),
focused on the primary business/technical purpose, clear to technical and
business readers. Classify the COPY in `type_of_copy`.

User:
<cobol_source>
{{.COBOL_SOURCE}}
</cobol_source>
```

### Calling it via the port (AGENTS-TO-BE §4 + AGENTS-IO §5)

```go
schema, _ := os.ReadFile("copy_summary.schema.json")
res, err := gen.Generate(ctx, []llm.Request{{
    System: sys, User: user, Model: model, MaxTokens: 1024,
    Format: &llm.Format{Name: "copy_summary", Schema: schema, Strict: true},
}})
// res[i].Text is validated JSON:
out, _ := llm.Unmarshal[struct {
    Summary    string `json:"summary"`
    TypeOfCopy string `json:"type_of_copy"`
}](res[i])
```

## Reasoning

**What changed, line for line:**

| Before | After |
|---|---|
| `Place your summary inside <summary> tags.` | gone — `summary` is a schema field |
| `Place your classification inside <type-of-copy> tag.` | gone — `type_of_copy` is a schema field |
| enum values listed in prose (`DATA_DIVISION, PROCEDURE_DIVISION, BMP_MAP`) | moved into the schema `enum` — **enforced**, single source of truth |
| entire "Important Formatting Requirements" block (no fences / verbatim-to-file) | gone — the transport is a validated JSON object, not extracted text |
| "Maximum 500 characters", "Written in Italian" | **kept** as prompt lines (best-effort / not schema-expressible) |
| `<cobol_source>` input marker | **kept** unchanged (§2) |

**What the schema enforces vs. what stays a prompt instruction:**

| Requirement | Mechanism after migration |
|---|---|
| Both outputs present | **Enforced** (`required`) — no more "a tag went missing" |
| `type_of_copy` is one of three values | **Enforced** (`enum`) — grammar-constrained across all three providers; the big win |
| `summary` ≤ 500 chars | **Best-effort** (`maxLength`) — keep as a prompt line *and* validate/trim client-side |
| `summary` in Italian | **Prompt instruction** — not expressible in JSON Schema |
| No stray code fences / "saved verbatim to a file" | **Deleted** — obsolete once the transport is a validated JSON object |

This is the asymmetry that proves AGENTS-IO §1's point: the **input** still uses the
`<cobol_source>` XML marker (§2 — kept), while the **output** tag machinery is gone, replaced by
the schema (§3 — migrated).
