// Package fs provides a filesystem-backed implementation of the tool executor
// interfaces defined in the parent tools package.
//
// FSExecutor implements both TextEditorExecutor and FileListExecutor, so a
// single instance can be passed to both WithEditor and WithLister:
//
//	exec := fs.New(baseDir)
//	ts   := tools.New(tools.WithEditor(exec), tools.WithLister(exec))
package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FSExecutor implements TextEditorExecutor and FileListExecutor using the local
// filesystem. All paths are joined with BaseDir, so the model never escapes the
// project root.
type FSExecutor struct {
	BaseDir string
}

// New creates an FSExecutor scoped to baseDir.
func New(baseDir string) *FSExecutor {
	return &FSExecutor{BaseDir: baseDir}
}

// abs resolves a model-supplied path against BaseDir.
//
// The text-editor tool is trained to send absolute paths, so a model may
// prepend a hallucinated sandbox root (e.g. "/home/user/app/paragraphs/x.cob")
// to what should be a project-relative path. A plain filepath.Join would then
// concatenate that garbage prefix onto BaseDir and fail to find the file.
//
// To recover, when the direct join does not exist we retry with progressively
// shorter suffixes of the requested path, keeping the first that resolves under
// BaseDir. This is unambiguous because the subdirectory suffix is preserved.
// When nothing resolves (e.g. creating a genuinely new file) we fall back to
// the plain join, so create/replace-new-file behaviour is unchanged.
func (f *FSExecutor) abs(path string) string {
	direct := filepath.Join(f.BaseDir, path)
	if _, err := os.Stat(direct); err == nil {
		return direct
	}

	parts := strings.Split(filepath.Clean(path), string(os.PathSeparator))
	for i := range parts {
		suffix := parts[i:]
		if len(suffix) == 0 || (len(suffix) == 1 && suffix[0] == "") {
			continue
		}
		cand := filepath.Join(append([]string{f.BaseDir}, suffix...)...)
		if _, err := os.Stat(cand); err == nil {
			return cand
		}
	}

	return direct
}

// guardWritePath rejects absolute paths for write operations. The text-editor
// tool is trained to send absolute paths, so a model may pass a hallucinated
// absolute path (e.g. "/home/user/app/src/java/Foo.java"). For reads, abs()
// recovers by suffix-matching an existing file, but a write to a not-yet-
// existing file cannot be recovered and would land in the wrong location.
// Returning an error lets the model see the mistake and retry with a
// project-relative path.
func guardWritePath(path string) error {
	if filepath.IsAbs(path) {
		return fmt.Errorf("path %q must be relative to the working directory, not absolute — retry with a project-relative path (e.g. \"src/java/Foo.java\")", path)
	}
	return nil
}

// View reads and returns the full content of the file at path.
func (f *FSExecutor) View(path string) (string, error) {
	data, err := os.ReadFile(f.abs(path))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Create writes content to path, creating parent directories as needed.
// Overwrites any existing file.
func (f *FSExecutor) Create(path, content string) (string, error) {
	if err := guardWritePath(path); err != nil {
		return "", err
	}
	full := f.abs(path)
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return "", err
	}
	if err := os.WriteFile(full, []byte(content), 0644); err != nil {
		return "", err
	}
	return "OK", nil
}

// StrReplace replaces the first occurrence of oldStr with newStr in the file at path.
// Returns an error if oldStr is not found.
func (f *FSExecutor) StrReplace(path, oldStr, newStr string) (string, error) {
	if err := guardWritePath(path); err != nil {
		return "", err
	}
	full := f.abs(path)
	data, err := os.ReadFile(full)
	if err != nil {
		return "", err
	}
	if !strings.Contains(string(data), oldStr) {
		return "", fmt.Errorf("old_str not found in %s", path)
	}
	updated := strings.Replace(string(data), oldStr, newStr, 1)
	if err := os.WriteFile(full, []byte(updated), 0644); err != nil {
		return "", err
	}
	return "OK", nil
}

// Insert inserts newStr after line insertLine (1-indexed; 0 means prepend).
// Follows the API contract of the built-in text_editor tool.
func (f *FSExecutor) Insert(path string, insertLine int, newStr string) (string, error) {
	if err := guardWritePath(path); err != nil {
		return "", err
	}
	full := f.abs(path)
	data, err := os.ReadFile(full)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(data), "\n")
	if insertLine < 0 || insertLine > len(lines) {
		return "", fmt.Errorf("insert_line %d out of range (file has %d lines)", insertLine, len(lines))
	}
	newLines := make([]string, 0, len(lines)+1)
	newLines = append(newLines, lines[:insertLine]...)
	newLines = append(newLines, newStr)
	newLines = append(newLines, lines[insertLine:]...)
	if err := os.WriteFile(full, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return "", err
	}
	return "OK", nil
}

// UndoEdit is not implemented by FSExecutor; it returns an error.
// Override by embedding FSExecutor and providing an edit-history store.
func (f *FSExecutor) UndoEdit(_ string) (string, error) {
	return "", fmt.Errorf("undo_edit not implemented by FSExecutor")
}

// List returns the relative paths of files matching pattern inside directory,
// one per line. pattern defaults to "*" when empty. Returns "(empty)" when no
// files match.
func (f *FSExecutor) List(directory, pattern string) (string, error) {
	if pattern == "" {
		pattern = "*"
	}
	entries, err := filepath.Glob(filepath.Join(f.abs(directory), pattern))
	if err != nil {
		return "", err
	}
	var rel []string
	for _, e := range entries {
		r, _ := filepath.Rel(f.BaseDir, e)
		rel = append(rel, r)
	}
	sort.Strings(rel)
	if len(rel) == 0 {
		return "(empty)", nil
	}
	return strings.Join(rel, "\n"), nil
}
