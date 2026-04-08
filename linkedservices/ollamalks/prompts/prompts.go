package prompts

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/rs/zerolog/log"
)

//go:embed templates/*
var templates embed.FS

const templatesRootFolder = "templates"

const (
	PromptProgramSummary = "program-summary-prompt"
)

type MessageSection struct {
	Name    string
	Ct      string
	Ext     string
	Summary bool
	Data    []byte
}

const (
	ScratchPadSection = "scratchpad"
	OverviewSection   = "overview"
	SummarySection    = "summary"
	MermaidSection    = "mermaid"

	ContentTypeTextPlain             = "text/plain"
	ContentTypeTextMarkdown          = "text/markdown"
	ContentTypeApplicationVndMermaid = "application/vnd.mermaid"
)

var SectionMap = map[string]MessageSection{
	ScratchPadSection: MessageSection{Name: ScratchPadSection, Ext: ".md", Summary: false, Ct: ContentTypeTextMarkdown},
	OverviewSection:   MessageSection{Name: OverviewSection, Ext: ".md", Summary: false, Ct: ContentTypeTextMarkdown},
	SummarySection:    MessageSection{Name: SummarySection, Ext: ".txt", Summary: true, Ct: ContentTypeTextPlain},
	MermaidSection:    MessageSection{Name: MermaidSection, Ext: ".mmd", Summary: false, Ct: ContentTypeApplicationVndMermaid},
}

type PromptTemplate struct {
	Name string `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`

	Vars         []string           `yaml:"vars,omitempty" mapstructure:"vars,omitempty" json:"vars,omitempty"`
	Sections     []string           `yaml:"sections,omitempty" mapstructure:"sections,omitempty" json:"sections,omitempty"`
	TemplateName string             `yaml:"template_name,omitempty" mapstructure:"template_name,omitempty" json:"template_name,omitempty"`
	Data         *template.Template `yaml:"-" mapstructure:"-" json:"-"`
}

var registry = map[string]PromptTemplate{
	PromptProgramSummary: {
		Name:         PromptProgramSummary,
		TemplateName: "program-summary-prompt.txt",

		Vars:     []string{"COBOL_SOURCE"},
		Sections: []string{ScratchPadSection, OverviewSection, SummarySection, MermaidSection},
	},
}

func RegisterTemplate(name string, vars []string, data string) (PromptTemplate, error) {
	var err error

	pt := PromptTemplate{
		Name:         name,
		Vars:         vars,
		TemplateName: name,
	}

	pt.Data, err = template.New("").Parse(data)
	if err != nil {
		return PromptTemplate{}, err
	}
	registry[name] = pt
	return pt, nil
}

func GetPrompt(name string) (PromptTemplate, error) {
	pt, ok := registry[name]
	if !ok {
		return PromptTemplate{}, fmt.Errorf("template %s not found", name)
	}

	if pt.Data != nil {
		return pt, nil
	}

	b, err := templates.ReadFile(templatesRootFolder + "/" + pt.TemplateName)
	if err != nil {
		return PromptTemplate{}, err
	}

	tmpl := template.Must(template.New("").Parse(string(b)))
	pt.Data = tmpl
	registry[name] = pt

	return pt, nil
}

type VariableType string

const (
	TextVariable VariableType = "text"
)

type Variable struct {
	Name  string
	Index int
	Ct    VariableType
	Value []byte
}

func (p PromptTemplate) Text(vars map[string]string) (string, error) {
	const semLogContext = "anthropic-lks-prompt::text"
	var err error

	for n, _ := range vars {
		found := false
		for _, v := range p.Vars {
			if v == n {
				found = true
				break
			}
		}

		if !found {
			err = fmt.Errorf("variable %s not found in prompt %s", n, p.Name)
			return "", err
		}
	}

	var buf bytes.Buffer
	err = p.Data.Execute(&buf, vars)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return "", err
	}
	parsedPrompt := buf.String()
	return parsedPrompt, err
}
