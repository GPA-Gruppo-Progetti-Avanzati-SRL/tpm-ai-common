# Example — `program-java-conversion` → JSON + deterministic render

> **Type:** mixed — a **document whose fields are all structured data** extracted from the source.
> → model emits **JSON (schema-enforced)**; code **renders the markdown**; constants injected at
> render, never asked of the model.
> **Source:** `cob-game/linkedservices/prompts_repo/program-java-conversion-prompt.txt`
> **Pattern (AGENTS-IO §3):** *document-as-rendering-of-data* — the third bucket, distinct from
> both [`copy-summary`](./copy-summary.md) (data → JSON) and [`node-summary`](./node-summary.md)
> (document → XML tags).

The output is a markdown "Context Block" meant to be read/embedded — but every field in it is
either an **enum** (`program-type`, `middleware`, `database`, entry-interface `type`, `how-invoked`),
a **fixed checklist** (`special-dialect-features`), a **table of typed rows**
(`known-external-programs`), or a **pre-set constant already held in Go**
(`target-java-package`, `output-folder`, `paragraph-dependency-table`, `source-layout`). The only
true free text is a couple of short "detail"/"notes" strings.

The key reframe: **that markdown block is not a document the model should *format* — it is a
*rendering of structured data*.** So split the two jobs the current prompt fuses.

## Before

The current prompt (structure — full text in the source file). It fuses four concerns into one
model call:

```text
1. a <jc-scratchapd> reasoning block (PROGRAM-ID? CICS/batch? middleware? entry interface? ...)

2. a "Pre-set constants" section (~30 lines) begging the model NOT to derive three values from
   source and to echo them verbatim:
     | target-java-package        | {{.JAVA_PACKAGE}}                |
     | output-folder              | {{.JAVA_OUTPUT_FOLDER}}          |
     | paragraph-dependency-table | {{.PARAGRAPH_DEPENDENCY_TABLE}}  |

3. field-detection RULES (the valuable IP): program-type priority order, database rules,
   entry-interface priority 1–4, known-external-programs scan, special-dialect checklist

4. a fixed markdown TEMPLATE the model must fill and emit inside <jc-context-block> tags:
     ## Part 1 — Context Block
     <!-- BEGIN CONTEXT BLOCK -->
     ### Program Name / Program Type / Middleware / Database / Entry Interface /
     ### Known External Programs (a markdown table) / Special Dialect Features /
     ### Target Java Package / Output Folder / Paragraph Dependency Table
     <!-- END CONTEXT BLOCK -->
   plus the usual "no code fences / exact format" ceremony.
```

## Next

### The schema — extracted fields only (no constants)

```json
{
  "type": "object",
  "additionalProperties": false,
  "required": ["program_name","program_type","middleware","database",
               "entry_interface","known_external_programs","special_dialect_features"],
  "properties": {
    "program_name": { "type": "string" },

    "program_type": {
      "type": "object", "additionalProperties": false, "required": ["kind"],
      "properties": {
        "kind":   { "type": "string",
                    "enum": ["CICS_ONLINE","CALLED_SUBPROGRAM","BATCH","OTHER"] },
        "detail": { "type": "string", "description": "required only when kind = OTHER" }
      }
    },

    "middleware": { "type": "string", "enum": ["CICS","IMS","NONE"] },

    "database": {
      "type": "array",
      "items": { "type": "string",
        "enum": ["DB2_EMBEDDED_SQL","VSAM_VIA_CICS","VSAM_SEQUENTIAL_INDEXED","NONE"] }
    },

    "entry_interface": {
      "type": "object", "additionalProperties": false, "required": ["type","java_entry"],
      "properties": {
        "type": { "type": "string",
          "enum": ["DFHCOMMAREA_VIA_COPY","DFHCOMMAREA_INLINE",
                   "LINKAGE_SECTION_ONLY","COMMAND_LINE_PARM","NONE"] },
        "copybook":   { "type": "string", "description": "copybook name when type = *_VIA_COPY" },
        "detail":     { "type": "string", "description": "top-level sub-groups, one line each" },
        "java_entry": { "type": "string", "description": "e.g. execute(C6Ap01asCommarea ca)" }
      }
    },

    "known_external_programs": {
      "type": "array",
      "items": {
        "type": "object", "additionalProperties": false,
        "required": ["program_name","how_invoked","is_dynamic"],
        "properties": {
          "program_name": { "type": "string" },
          "how_invoked":  { "type": "string",
                            "enum": ["EXEC_CICS_LINK","EXEC_CICS_XCTL","CALL"] },
          "is_dynamic":   { "type": "boolean" },
          "commarea":     { "type": "string" },
          "notes":        { "type": "string" },
          "move_targets": { "type": "array", "items": { "type": "string" },
                            "description": "MOVE targets for dynamic (variable-name) calls" }
        }
      }
    },

    "special_dialect_features": {
      "type": "array",
      "items": { "type": "string",
        "enum": ["DECIMAL_POINT_IS_COMMA","COMP_3_PACKED_DECIMAL","COMP_BINARY",
                 "VARCHAR_LEN_TEXT_PAIRS","CBL_OR_PROCESS_OPTIONS",
                 "COPY_REPLACING_OR_IN_LIBRARY"] }
    }
  }
}
```

**Absent by design:** `target_java_package`, `output_folder`, `paragraph_dependency_table`,
`source_layout` — injected at render, not asked of the model. The `OTHER`+`detail` shape on
`program_type` is the reusable pattern for any "other — describe" escape.

### The revised prompt

Keep the field-detection *rules* (the IP); drop the markdown template, the constants section, and
the output-tag ceremony. Input marker stays (§2).

```text
System:
You extract Java-conversion metadata from a CICS/COBOL program and return it ONLY as a
JSON object matching the provided schema. Apply the rules below to decide each value;
for every enum field choose exactly one allowed token.

program_type  — priority: any EXEC CICS → CICS_ONLINE; else USING non-DFHCOMMAREA areas,
                no CICS → CALLED_SUBPROGRAM; else no CICS and no USING → BATCH; else OTHER (+detail).
database      — EXEC SQL → DB2_EMBEDDED_SQL; EXEC CICS file cmd with FILE() → VSAM_VIA_CICS;
                plain READ/WRITE on an FD → VSAM_SEQUENTIAL_INDEXED; none → NONE; list all that apply.
entry_interface — [the current priority 1–4 rules, verbatim, producing type + copybook + detail + java_entry]
known_external_programs — scan PROCEDURE DIVISION for EXEC CICS LINK/XCTL and CALL; for each set
                program_name, how_invoked, is_dynamic; for dynamic calls fill move_targets; commarea/notes as known.
special_dialect_features — include each checklist token that applies (empty array if none).

User:
<cobol_source>
{{.COBOL_SOURCE}}
</cobol_source>
```

The scratchpad becomes native adaptive thinking (drop `<jc-scratchapd>`). If the reasoning must be
retained as an artifact, keep it as a **separate** top-level `scratchpad` string or a separate XML
block — not part of the enforced data object.

### The rendering step (keeps the downstream contract)

A Go `text/template` takes the validated struct **plus** the three injected constants and emits the
identical `## Part 1 — Context Block` markdown (same `<!-- BEGIN/END CONTEXT BLOCK -->` markers), so
anything that currently greps the block keeps working. The cleaner end-state is to have downstream
consume the JSON directly and drop the markdown, but the template path migrates without touching
consumers.

## Reasoning

Why the hybrid is strictly better than the current single-call prompt:

- **Enums get enforced.** `program-type`, `middleware`, `database`, entry-interface `type`,
  `how-invoked` become fixed tokens — no "other — describe" drift leaking into whatever greps the
  block downstream.
- **The external-programs table becomes an array of typed objects.** The "write a single row of
  `none`" and "add a final note row if all static" ceremony disappears — code renders those from an
  empty array / an `is_dynamic` flag.
- **`special-dialect-features` becomes a multi-select enum** (a fixed checklist → an `enum`-item
  array), far more reliable than a free-text list the model assembles.
- **The pre-set constants stop being the model's problem entirely.** ~30 lines currently beg the
  model *not* to regenerate the dependency table and to echo three values verbatim — fragile and
  wasteful. In the hybrid design you **don't ask for them at all**; they're injected at render, and
  the risk of the model "helpfully" rebuilding the dependency table from source goes to zero.
- **The markdown can't be malformed** — it comes from a template, not the model, so the whole
  "no code fences / exact format" guidance is gone.

**The general pattern this establishes:** *document-that-is-a-rendering-of-data → model emits JSON
(schema-enforced); code renders the document; constants/fixed text injected at render, never asked
of the model.*
