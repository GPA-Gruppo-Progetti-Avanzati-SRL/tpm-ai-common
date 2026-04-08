package prompts

import (
	"bufio"
	"bytes"
	"embed"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/rs/zerolog/log"
)

//go:embed templates/*
var templates embed.FS

const templatesRootFolder = "templates"

const (
	ScratchPadSection = "scratchpad"
	OverviewSection   = "overview"
	SummarySection    = "summary"
	MermaidSection    = "mermaid"

	ContentTypeTextPlain             = "text/plain"
	ContentTypeTextMarkdown          = "text/markdown"
	ContentTypeApplicationVndMermaid = "application/vnd.mermaid"
)

type MessageSection struct {
	Name    string
	Ct      string
	Ext     string
	Summary bool
	Data    []byte
}

var sectionMap = map[string]MessageSection{
	ScratchPadSection: MessageSection{Name: ScratchPadSection, Ext: ".md", Summary: false, Ct: ContentTypeTextMarkdown},
	OverviewSection:   MessageSection{Name: OverviewSection, Ext: ".md", Summary: false, Ct: ContentTypeTextMarkdown},
	SummarySection:    MessageSection{Name: SummarySection, Ext: ".txt", Summary: true, Ct: ContentTypeTextPlain},
	MermaidSection:    MessageSection{Name: MermaidSection, Ext: ".mmd", Summary: false, Ct: ContentTypeApplicationVndMermaid},
}

const (
	PromptProgramSummary = "program-summary-prompt"
)

type PromptTemplate struct {
	Name         string                     `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	System       []anthropic.TextBlockParam `yaml:"-" mapstructure:"-" json:"-"`
	Vars         []string                   `yaml:"vars,omitempty" mapstructure:"vars,omitempty" json:"vars,omitempty"`
	Sections     []string                   `yaml:"sections,omitempty" mapstructure:"sections,omitempty" json:"sections,omitempty"`
	TemplateName string                     `yaml:"template_name,omitempty" mapstructure:"template_name,omitempty" json:"template_name,omitempty"`
	Data         *template.Template         `yaml:"-" mapstructure:"-" json:"-"`
}

var registry = map[string]PromptTemplate{
	PromptProgramSummary: {
		Name:         PromptProgramSummary,
		TemplateName: "program-summary-prompt.txt",
		System:       []anthropic.TextBlockParam{{Text: "You are a senior Cobol developer with extensive knowledge in mainframe cobol cics programming"}},
		Vars:         []string{"COBOL_SOURCE"},
		Sections:     []string{ScratchPadSection, OverviewSection, SummarySection, MermaidSection},
	},
}

func RegisterTemplate(name string, system string, vars []string, data string) (PromptTemplate, error) {
	var err error

	pt := PromptTemplate{
		Name:         name,
		Vars:         vars,
		TemplateName: name,
	}

	if system != "" {
		pt.System = []anthropic.TextBlockParam{{Text: system}}
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

func (p PromptTemplate) Text(vars map[string]string) ([]byte, error) {
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
			return nil, err
		}
	}

	var buf bytes.Buffer
	err = p.Data.Execute(&buf, vars)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}
	parsedPrompt := buf.String()
	b, err := json.Marshal(parsedPrompt)
	return b, err
}

const (
	StatusZero = iota
	StatusInDelimiter

	LineTypeLine        string = "line"
	LineTypeBODelimiter        = "bo-delimiter"
	LineTypeEODelimiter        = "eo-delimiter"
)

func (p PromptTemplate) ParseMessage(message *anthropic.Message) (map[string]MessageSection, error) {
	const semLogContext = "anthropic-lks-client::process-message"
	var err error

	if len(message.Content) != 1 {
		err = errors.New("expected a single message")
		log.Warn().Err(err).Msg(semLogContext)
		return nil, err
	}

	for _, m := range message.Content {
		if m.Type == "text" {
			content, err := ParseTextContent(m.Text, p.Sections)
			if err != nil {
				log.Error().Err(err).Msg(semLogContext)
				return nil, err
			}

			parsedPrompt := make(map[string]MessageSection)
			for sn, sv := range content {
				if si, ok := sectionMap[sn]; ok {
					parsedPrompt[sn] = MessageSection{Name: si.Name, Ext: si.Ext, Summary: si.Summary, Data: []byte(sv)}
				} else {
					log.Warn().Str("section", sn).Msg(semLogContext + " - unknown section")
				}
			}

			return parsedPrompt, nil
		}
	}

	return nil, nil
}

func ParseTextContent(text string, delimiters []string) (map[string]string, error) {
	const semLogContext = "anthropic-cobol::processResponse"

	var beginDelimiters []string
	var endDelimiters []string
	for _, d := range delimiters {
		beginDelimiters = append(beginDelimiters, fmt.Sprintf("<%s>", d))
		endDelimiters = append(endDelimiters, fmt.Sprintf("</%s>", d))
	}

	bytesReader := strings.NewReader(text)
	bufReader := bufio.NewReader(bytesReader)

	var m map[string]string
	var sb strings.Builder
	var currentDelimiter, delim string
	var lineType string
	status := StatusZero
	line, isPrefix, err := bufReader.ReadLine()
	for err == nil && !isPrefix {

		sline := string(line)
		sline, lineType, delim = typeOfLine(sline, delimiters)

		switch lineType {
		case LineTypeBODelimiter:
			if status != StatusZero {
				return nil, fmt.Errorf("found bof-scrathcpad when status in %d", status)
			}
			currentDelimiter = delim
			if sline != "" {
				sb.WriteString(sline)
				sb.WriteString("\n")
			}
			status = StatusInDelimiter
		case LineTypeEODelimiter:
			if status != StatusInDelimiter || (status == StatusInDelimiter && currentDelimiter != delim) {
				return nil, fmt.Errorf("found end of delimiter %s when status in %d and current delimiter is %s", delim, status, currentDelimiter)
			}
			if sline != "" {
				sb.WriteString(sline)
				sb.WriteString("\n")
			}
			status = StatusZero
			if m == nil {
				m = make(map[string]string)
			}
			m[currentDelimiter] = sb.String()
			sb = strings.Builder{}
		default:
			if status != StatusZero {
				sb.WriteString(sline)
				sb.WriteString("\n")
			}
		}

		line, isPrefix, err = bufReader.ReadLine()
	}

	if err != nil && err != io.EOF {
		return nil, err
	}

	if isPrefix {
		return nil, errors.New("buffer too small")
	}

	return m, nil
}

func typeOfLine(line string, delimiters []string) (string, string, string) {
	line = strings.TrimSpace(line)
	for _, d := range delimiters {
		start := fmt.Sprintf("<%s>", d)
		if ndx := strings.Index(line, start); ndx >= 0 {
			line = line[ndx+len(start):]
			line = strings.TrimSpace(line)
			return line, LineTypeBODelimiter, d
		}

		end := fmt.Sprintf("</%s>", d)
		if ndx := strings.Index(line, end); ndx >= 0 {
			line = line[:ndx]
			line = strings.TrimSpace(line)
			return line, LineTypeEODelimiter, d
		}
	}

	return line, LineTypeLine, ""
}
