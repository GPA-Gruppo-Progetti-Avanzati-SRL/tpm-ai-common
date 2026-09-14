package agents

import (
	"testing"
)

const sampleDocumentYAML = `
- name: scratchpad
  ext: "md"
  is-document: true
  ct: text/markdown
  title: Scratchpad
  description: Scratchpad for the program
- name: overview
  ext: "md"
  is-document: true
  ct: text/markdown
  title: Overview
  description: Overview del programma
- name: summary
  ext: "txt"
  is-document: false
  ct: text/plain
  title: Descrizione Breve
  description: Descrizione Breve del programma
- name: flowchart
  ext: "mmd"
  is-document: true
  ct: text/vnd.mermaid
  title: Flowchart
  description: Flowchart del programma
`

func TestNewXMLSectionedDocumentFromYAML(t *testing.T) {
	doc, err := NewXMLSectionsFromYAML([]byte(sampleDocumentYAML))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc) != 4 {
		t.Fatalf("expected 4 sections, got %d", len(doc))
	}

	want := XMLSections{
		{Name: "scratchpad", Ext: "md", IsDocument: true, Ct: "text/markdown", Title: "Scratchpad", Description: "Scratchpad for the program"},
		{Name: "overview", Ext: "md", IsDocument: true, Ct: "text/markdown", Title: "Overview", Description: "Overview del programma"},
		{Name: "summary", Ext: "txt", IsDocument: false, Ct: "text/plain", Title: "Descrizione Breve", Description: "Descrizione Breve del programma"},
		{Name: "flowchart", Ext: "mmd", IsDocument: true, Ct: "text/vnd.mermaid", Title: "Flowchart", Description: "Flowchart del programma"},
	}

	for i, w := range want {
		got := doc[i]
		if got.Name != w.Name || got.Ext != w.Ext || got.IsDocument != w.IsDocument ||
			got.Ct != w.Ct || got.Title != w.Title || got.Description != w.Description {
			t.Errorf("section %d mismatch:\n got  %+v\n want %+v", i, got, w)
		}
	}
}

func TestDelimiterNames(t *testing.T) {
	doc, err := NewXMLSectionsFromYAML([]byte(sampleDocumentYAML))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := doc.DelimiterNames()
	want := []string{"scratchpad", "overview", "summary", "flowchart"}
	if len(got) != len(want) {
		t.Fatalf("expected %d names, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("name %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestSectionByName(t *testing.T) {
	doc, err := NewXMLSectionsFromYAML([]byte(sampleDocumentYAML))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s, ok := doc.SectionByName("summary")
	if !ok {
		t.Fatal("expected to find section 'summary'")
	}
	if s.Ct != "text/plain" {
		t.Errorf("summary ct: got %q, want text/plain", s.Ct)
	}

	if _, ok := doc.SectionByName("does-not-exist"); ok {
		t.Error("expected nil for unknown section")
	}
}

func TestExtractFromText(t *testing.T) {
	doc, err := NewXMLSectionsFromYAML([]byte(sampleDocumentYAML))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text := `<scratchpad>
some notes
</scratchpad>
<overview>
program overview line 1
program overview line 2
</overview>
<summary>
a short summary
</summary>
<flowchart>
flowchart TD
  A --> B
</flowchart>`

	if err := doc.ExtractFromText(text); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cases := map[string]string{
		"scratchpad": "some notes\n",
		"overview":   "program overview line 1\nprogram overview line 2\n",
		"summary":    "a short summary\n",
		// Content-line indentation is preserved (only tag lines are trimmed).
		"flowchart": "flowchart TD\n  A --> B\n",
	}

	for name, wantData := range cases {
		s, ok := doc.SectionByName(name)
		if !ok {
			t.Fatalf("section %q not found", name)
		}
		if string(s.Data) != wantData {
			t.Errorf("section %q data mismatch:\n got  %q\n want %q", name, string(s.Data), wantData)
		}
	}
}

func TestExtractParts(t *testing.T) {
	doc, err := NewXMLSectionsFromYAML([]byte(sampleDocumentYAML))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text := `<overview>
line one
line two
</overview>
<summary>
short
</summary>`

	parts, err := doc.ExtractParts(text)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := parts["overview"]; got != "line one\nline two\n" {
		t.Errorf("overview: got %q", got)
	}
	if got := parts["summary"]; got != "short\n" {
		t.Errorf("summary: got %q", got)
	}
	if _, ok := parts["flowchart"]; ok {
		t.Error("flowchart should be absent from the map")
	}

	// ExtractParts must not mutate the document's section Data.
	for _, s := range doc {
		if s.Data != nil {
			t.Errorf("section %q Data should be nil after ExtractParts, got %q", s.Name, string(s.Data))
		}
	}
}

func TestExtractFromTextMissingRequired(t *testing.T) {
	doc := XMLSections{
		{Name: "overview", Required: true},
		{Name: "summary"},
	}

	text := `<summary>
only the summary is present
</summary>`

	err := doc.ExtractFromText(text)
	if err == nil {
		t.Fatal("expected an error for missing required section, got nil")
	}
}

func TestExtractFromTextMissingOptional(t *testing.T) {
	doc := XMLSections{
		{Name: "overview"},
		{Name: "summary"},
	}

	text := `<summary>
only the summary is present
</summary>`

	if err := doc.ExtractFromText(text); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s, ok := doc.SectionByName("overview"); !ok || s.Data != nil {
		t.Errorf("expected overview data to be nil, got %q", string(s.Data))
	}
	if s, ok := doc.SectionByName("summary"); !ok || string(s.Data) != "only the summary is present\n" {
		t.Errorf("unexpected summary data: %q", s.Data)
	}
}
