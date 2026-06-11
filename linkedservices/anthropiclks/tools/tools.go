package tools

import (
	"encoding/json"
	"fmt"

	anthropic "github.com/anthropics/anthropic-sdk-go"
)

// TextEditorToolName is the name the model uses in tool_use blocks.
// Derived from the SDK param type so it stays correct on SDK upgrades —
// update the struct reference here when bumping to a newer text editor version.
var TextEditorToolName = string(anthropic.ToolTextEditor20250728Param{}.Name.Default())

// ListFilesToolName is the name the model uses for the custom list_files tool.
const ListFilesToolName = "list_files"

// textEditorParam is the SDK param for the built-in text editor tool.
// Name and Type fields have correct defaults and are omitted from the JSON — the
// API injects the full input schema server-side; no Properties block is needed.
var textEditorParam = anthropic.ToolUnionParam{
	OfTextEditor20250728: &anthropic.ToolTextEditor20250728Param{},
}

// listFilesParam is the SDK param for the custom list_files tool.
// The full JSON schema must be provided because this is not a built-in API tool.
var listFilesParam = anthropic.ToolUnionParam{OfTool: &anthropic.ToolParam{
	Name:        ListFilesToolName,
	Description: anthropic.String("List files in a directory. Path is relative to the base directory."),
	InputSchema: anthropic.ToolInputSchemaParam{
		Properties: map[string]any{
			"directory": map[string]any{"type": "string", "description": "Directory path relative to base"},
			"pattern":   map[string]any{"type": "string", "description": "Optional glob pattern, e.g. '*.cob'"},
		},
		Required: []string{"directory"},
	},
}}

// TextEditorExecutor backs the str_replace_based_edit_tool commands.
// Implementations receive paths relative to whatever base directory they manage.
type TextEditorExecutor interface {
	View(path string) (string, error)
	Create(path, content string) (string, error)
	StrReplace(path, oldStr, newStr string) (string, error)
	Insert(path string, insertLine int, newStr string) (string, error)
	UndoEdit(path string) (string, error)
}

// FileListExecutor backs the list_files custom tool.
type FileListExecutor interface {
	List(directory, pattern string) (string, error)
}

// ToolSet bundles Anthropic SDK tool params with their executor implementations.
// Only tools that have executors assigned are included in Params().
type ToolSet struct {
	editor TextEditorExecutor
	lister FileListExecutor
}

// Option configures a ToolSet.
type Option func(*ToolSet)

// WithEditor activates the text editor tool backed by the given executor.
func WithEditor(e TextEditorExecutor) Option {
	return func(ts *ToolSet) { ts.editor = e }
}

// WithLister activates the list_files tool backed by the given executor.
func WithLister(l FileListExecutor) Option {
	return func(ts *ToolSet) { ts.lister = l }
}

// New creates a ToolSet with the supplied options.
func New(opts ...Option) *ToolSet {
	ts := &ToolSet{}
	for _, o := range opts {
		o(ts)
	}
	return ts
}

// Params returns the SDK tool params for every tool that has an executor set.
// Pass the result directly to anthropic.MessageNewParams.Tools.
func (ts *ToolSet) Params() []anthropic.ToolUnionParam {
	var params []anthropic.ToolUnionParam
	if ts.editor != nil {
		params = append(params, textEditorParam)
	}
	if ts.lister != nil {
		params = append(params, listFilesParam)
	}
	return params
}

// textEditorInput mirrors the JSON the model sends for str_replace_based_edit_tool.
type textEditorInput struct {
	Command    string `json:"command"`
	Path       string `json:"path"`
	OldStr     string `json:"old_str"`
	NewStr     string `json:"new_str"`
	FileText   string `json:"file_text"`
	InsertLine int    `json:"insert_line"`
}

// listFilesInput mirrors the JSON the model sends for list_files.
type listFilesInput struct {
	Directory string `json:"directory"`
	Pattern   string `json:"pattern"`
}

// Execute dispatches a single tool call from the agent loop.
// name is the tool name from the ToolUseBlock; rawInput is the raw JSON input.
// Returns the result string, or "ERROR: ..." on failure.
func (ts *ToolSet) Execute(name string, rawInput json.RawMessage) string {
	switch name {
	case TextEditorToolName:
		if ts.editor == nil {
			return "ERROR: text editor tool not configured"
		}
		var in textEditorInput
		if err := json.Unmarshal(rawInput, &in); err != nil {
			return "ERROR: could not parse tool input: " + err.Error()
		}
		return ts.dispatchEditor(in)

	case ListFilesToolName:
		if ts.lister == nil {
			return "ERROR: list_files tool not configured"
		}
		var in listFilesInput
		if err := json.Unmarshal(rawInput, &in); err != nil {
			return "ERROR: could not parse tool input: " + err.Error()
		}
		result, err := ts.lister.List(in.Directory, in.Pattern)
		if err != nil {
			return "ERROR: " + err.Error()
		}
		return result

	default:
		return "ERROR: unknown tool '" + name + "'"
	}
}

func (ts *ToolSet) dispatchEditor(in textEditorInput) string {
	var result string
	var err error
	switch in.Command {
	case "view":
		result, err = ts.editor.View(in.Path)
	case "create":
		result, err = ts.editor.Create(in.Path, in.FileText)
	case "str_replace":
		result, err = ts.editor.StrReplace(in.Path, in.OldStr, in.NewStr)
	case "insert":
		result, err = ts.editor.Insert(in.Path, in.InsertLine, in.NewStr)
	case "undo_edit":
		result, err = ts.editor.UndoEdit(in.Path)
	default:
		return fmt.Sprintf("ERROR: unknown command %q for %s", in.Command, TextEditorToolName)
	}
	if err != nil {
		return "ERROR: " + err.Error()
	}
	return result
}
