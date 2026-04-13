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
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

const semLogContextBasePrompt = "prompt::"

type MessagePart struct {
	Name      string `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	Ct        string `yaml:"ct,omitempty" mapstructure:"ct,omitempty" json:"ct,omitempty"`
	Ext       string `yaml:"ext,omitempty" mapstructure:"ext,omitempty" json:"ext,omitempty"`
	IsSummary bool   `yaml:"is-summary,omitempty" mapstructure:"is-summary,omitempty" json:"is-summary,omitempty"`
	Required  bool   `yaml:"required,omitempty" mapstructure:"required,omitempty" json:"required,omitempty"`
	Data      []byte `yaml:"-" mapstructure:"-" json:"-"`
}

type PromptTemplate struct {
	Name           string             `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	System         string             `yaml:"system" mapstructure:"system" json:"system"`
	Vars           []string           `yaml:"vars,omitempty" mapstructure:"vars,omitempty" json:"vars,omitempty"`
	Sections       []MessagePart      `yaml:"parts,omitempty" mapstructure:"parts,omitempty" json:"parts,omitempty"`
	TemplateFn     string             `yaml:"template-fn,omitempty" mapstructure:"template-fn,omitempty" json:"template-fn,omitempty"`
	SchemaFn       string             `yaml:"schema-fn,omitempty" mapstructure:"schema-fn,omitempty" json:"schema-fn,omitempty"`
	ParsedTemplate *template.Template `yaml:"-" mapstructure:"-" json:"-"`
	TemplateSchema string             `yaml:"-" mapstructure:"-" json:"-"`
}

func NewPrompt(fn string) (PromptTemplate, error) {
	const semLogContext = semLogContextBasePromptRegistry + "new-prompt"

	fileContent, err := os.ReadFile(fn)
	if err != nil {
		log.Error().Err(err).Str("fn", fn).Msg(semLogContext)
		return PromptTemplate{}, err
	}

	var prmpt PromptTemplate
	err = yaml.Unmarshal(fileContent, &prmpt)
	if err != nil {
		log.Error().Err(err).Str("fn", fn).Msg(semLogContext + " failed to unmarshal prompt template")
		return PromptTemplate{}, err
	}

	tmplFn := filepath.Join(path.Dir(fn), prmpt.TemplateFn)
	tmplBody, err := os.ReadFile(tmplFn)
	if err != nil {
		log.Error().Err(err).Str("path", tmplFn).Msg(semLogContext + " failed to read prompt template")
		return PromptTemplate{}, err
	}
	tmpl := template.Must(template.New("").Parse(string(tmplBody)))
	prmpt.ParsedTemplate = tmpl

	if prmpt.SchemaFn != "" {
		schemaFn := filepath.Join(path.Dir(fn), prmpt.SchemaFn)
		schemaBody, err := os.ReadFile(schemaFn)
		if err != nil {
			log.Error().Err(err).Str("path", schemaFn).Msg(semLogContext + " failed to read schema template")
			return PromptTemplate{}, err
		}
		prmpt.TemplateSchema = string(schemaBody)
	}
	return prmpt, nil
}

func NewPromptFromEmbeddedFS(rootFolder string, embeddedTemplates embed.FS, fn string, fileContent []byte) (PromptTemplate, error) {
	const semLogContext = semLogContextBasePromptRegistry + "new-prompt-from-embedded-fs"
	var prmpt PromptTemplate
	err := yaml.Unmarshal(fileContent, &prmpt)
	if err != nil {
		log.Error().Err(err).Str("fn", fn).Msg(semLogContext + " failed to unmarshal prompt template")
		return PromptTemplate{}, err
	}

	tmplFn := path.Join(rootFolder, path.Dir(fn), prmpt.TemplateFn)
	tmplBody, err := embeddedTemplates.ReadFile(tmplFn)
	if err != nil {
		log.Error().Err(err).Str("path", tmplFn).Msg(semLogContext + " failed to read prompt template")
		return PromptTemplate{}, err
	}
	tmpl := template.Must(template.New("").Parse(string(tmplBody)))
	prmpt.ParsedTemplate = tmpl

	schemaFn := path.Join(rootFolder, path.Dir(fn), prmpt.SchemaFn)
	schemaBody, err := embeddedTemplates.ReadFile(schemaFn)
	if err != nil {
		log.Error().Err(err).Str("path", schemaFn).Msg(semLogContext + " failed to read schema template")
		return PromptTemplate{}, err
	}
	prmpt.TemplateSchema = string(schemaBody)
	return prmpt, nil
}

func (p PromptTemplate) SectionNames() []string {
	if p.Sections == nil {
		return nil
	}

	var arr []string
	for _, s := range p.Sections {
		arr = append(arr, s.Name)
	}

	return arr
}

func (p PromptTemplate) GetSectionByName(n string) (MessagePart, bool) {
	if p.Sections == nil {
		return MessagePart{}, false
	}

	n = strings.ToLower(n)
	for _, s := range p.Sections {
		if strings.ToLower(s.Name) == n {
			return s, true
		}
	}

	return MessagePart{}, false
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
	const semLogContext = semLogContextBasePrompt + "text"
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
	err = p.ParsedTemplate.Execute(&buf, vars)
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

func (p PromptTemplate) ParseTextContent(text string) (map[string]MessagePart, error) {
	const semLogContext = semLogContextBasePrompt + "parse-text-content"
	content, err := extractTextForParts(text, p.SectionNames())
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		fmt.Println(text)
		return nil, err
	}

	parsedPrompt := make(map[string]MessagePart)
	for _, s := range p.Sections {
		if data, ok := content[s.Name]; ok {
			parsedPrompt[s.Name] = MessagePart{Name: s.Name, Ext: s.Ext, Ct: s.Ct, IsSummary: s.IsSummary, Data: []byte(data)}
		} else if s.Required {
			err = fmt.Errorf("required section %s not found in prompt %s response", s.Name, p.Name)
			return nil, err
		}
	}

	return parsedPrompt, nil
}

func extractTextForParts(text string, delimiters []string) (map[string]string, error) {
	const semLogContext = semLogContextBasePrompt + "parse-text-content"

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
				return nil, fmt.Errorf("found %s when status in %d (%s)", delim, status, currentDelimiter)
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
	const semLogContext = semLogContextBasePrompt + "type-of-line"
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
