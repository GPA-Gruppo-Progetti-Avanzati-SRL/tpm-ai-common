# tools — Anthropic tool wrapper

This package provides a thin, swappable abstraction over the Anthropic tool system
for use in agentic loops built with the Go SDK.

## Design rationale

The Anthropic API distinguishes two kinds of tools:

| Kind                            | SDK field                         | Schema                                         |
|---------------------------------|-----------------------------------|------------------------------------------------|
| Built-in (text editor, bash, …) | `OfTextEditor*`, `OfBashTool*`, … | Injected server-side — not sent in the request |
| Custom                          | `OfTool` + `ToolParam`            | You provide `InputSchema`                      |

Using a built-in tool type has two advantages over a custom tool with an equivalent
hand-written schema:

1. **Smaller requests** — the schema is never serialised into the JSON payload.
2. **Better model behaviour** — the model was trained extensively on the exact field
   names and command semantics of these tools; a custom tool with a different name
   or field names forces the model to learn from the description alone.

`list_files` has no built-in equivalent, so it stays as a custom `OfTool`.

## Tool versions and names

Three versioned text editor params exist in SDK v1.45.0:

| Struct                        | Tool name (model sees)        | API type string        |
|-------------------------------|-------------------------------|------------------------|
| `ToolTextEditor20250124Param` | `str_replace_editor`          | `text_editor_20250124` |
| `ToolTextEditor20250429Param` | `str_replace_based_edit_tool` | `text_editor_20250429` |
| `ToolTextEditor20250728Param` | `str_replace_based_edit_tool` | `text_editor_20250728` |

This package uses `ToolTextEditor20250728Param` (the latest). All Name/Type fields
have correct zero-value defaults and are omitted from the request body.

## Interfaces

```go
type TextEditorExecutor interface {
    View(path string) (string, error)
    Create(path, content string) (string, error)
    StrReplace(path, oldStr, newStr string) (string, error)
    Insert(path string, insertLine int, newStr string) (string, error)
    UndoEdit(path string) (string, error)
}

type FileListExecutor interface {
    List(directory, pattern string) (string, error)
}
```

Both interfaces receive **relative paths**. It is the executor's responsibility to
resolve them against whatever base directory or storage it manages.

## ToolSet

`ToolSet` is the main type. It holds executor references and exposes two methods
consumed by the agent loop:

```go
ts.Params()                            // → []anthropic.ToolUnionParam, for MessageNewParams.Tools
ts.Execute(name, rawInput)             // → string result, for tool_result blocks
```

Only tools with a non-nil executor are included in `Params()`.

## Usage

```go
import (
    "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/anthropiclks/tools"
    "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/anthropiclks/tools/fs"
)

// Single executor implementing both interfaces (typical case).
exec := fs.New(baseDir)
ts   := tools.New(tools.WithEditor(exec), tools.WithLister(exec))

// In the agent loop:
params := anthropic.MessageNewParams{
    Tools: ts.Params(),
    // ...
}

// Dispatching a tool call:
result := ts.Execute(toolUseBlock.Name, toolUseBlock.Input)
```

## Filesystem implementation (`fs` sub-package)

`fs.FSExecutor` implements both interfaces against the local filesystem, scoped to
`BaseDir`. All paths are joined with `filepath.Join(BaseDir, path)`.

| Method       | Behaviour                                                            |
|--------------|----------------------------------------------------------------------|
| `View`       | `os.ReadFile`                                                        |
| `Create`     | `os.MkdirAll` + `os.WriteFile` (overwrites)                          |
| `StrReplace` | Read → `strings.Replace(…, 1)` → write; error if `old_str` absent    |
| `Insert`     | Split lines → insert after line N (1-indexed; 0 = prepend) → write   |
| `UndoEdit`   | Returns error — no history store. Embed and override to add one.     |
| `List`       | `filepath.Glob` → sorted relative paths; `"(empty)"` when no matches |

## Extending

To swap the execution backend (e.g. in-memory store, remote filesystem, test mock):

```go
type MyEditor struct{ /* ... */ }

func (m *MyEditor) View(path string) (string, error)        { /* ... */ }
func (m *MyEditor) Create(path, content string) (string, error) { /* ... */ }
// ... implement remaining methods ...

ts := tools.New(tools.WithEditor(&MyEditor{}), tools.WithLister(fs.New(baseDir)))
```

Any type that satisfies `TextEditorExecutor` or `FileListExecutor` can be mixed in.
