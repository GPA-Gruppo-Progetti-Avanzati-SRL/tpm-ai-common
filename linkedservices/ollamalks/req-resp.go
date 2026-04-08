package ollamalks

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/ollamalks/prompts"
	"github.com/rs/zerolog/log"
)

type Request struct {
	Vars map[string]prompts.Variable
}

func (r Request) TextVariables() map[string]string {
	var m map[string]string
	for k, v := range r.Vars {
		switch v.Ct {
		case prompts.TextVariable:
			if len(m) == 0 {
				m = make(map[string]string)
			}
			m[k] = string(v.Value)
		}
	}

	return m
}

type RequestParam func(r *Request)

func WithPromptInput(vt prompts.VariableType, n string, ndx int, v []byte) RequestParam {
	return func(o *Request) {
		if o.Vars == nil {
			o.Vars = make(map[string]prompts.Variable)
		}

		if ndx < 0 {
			ndx = 0
		}
		o.Vars[n] = prompts.Variable{Name: n, Index: ndx, Ct: vt, Value: v}
	}
}

type Response struct {
	Content map[string]prompts.MessageSection
}

type Message struct {
	Model  string `yaml:"model,omitempty" mapstructure:"model,omitempty" json:"model,omitempty"`
	Stream bool   `yaml:"stream" mapstructure:"stream" json:"stream"`
	Prompt string `yaml:"prompt,omitempty" mapstructure:"prompt,omitempty" json:"prompt,omitempty"`
}

type NonStreamedResponse struct {
	Model              string    `json:"model" yaml:"model" mapstructure:"model"`
	CreatedAt          time.Time `json:"created_at" yaml:"created_at" mapstructure:"created_at"`
	Response           string    `json:"response" yaml:"response" mapstructure:"response"`
	Done               bool      `json:"done" yaml:"done" mapstructure:"done"`
	DoneReason         string    `json:"done_reason" yaml:"done_reason" mapstructure:"done_reason"`
	Context            []int     `json:"context" yaml:"context" mapstructure:"context"`
	TotalDuration      int64     `json:"total_duration" yaml:"total_duration" mapstructure:"total_duration"`
	LoadDuration       int       `json:"load_duration" yaml:"load_duration" mapstructure:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count" yaml:"prompt_eval_count" mapstructure:"prompt_eval_count"`
	PromptEvalDuration int       `json:"prompt_eval_duration" yaml:"prompt_eval_duration" mapstructure:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count" yaml:"eval_count" mapstructure:"eval_count"`
	EvalDuration       int       `json:"eval_duration" yaml:"eval_duration" mapstructure:"eval_duration"`
}

func (nr *NonStreamedResponse) RetrieveContentFromResponse(sections []string) (map[string]prompts.MessageSection, error) {

	const semLogContext = "anthropic-lks-client::process-message"
	var err error

	if len(nr.Response) == 0 {
		err = errors.New("expected a single message")
		log.Warn().Err(err).Msg(semLogContext)
		return nil, err
	}

	content, err := ParseTextContent(nr.Response, sections)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	parsedPrompt := make(map[string]prompts.MessageSection)
	for sn, sv := range content {
		if si, ok := prompts.SectionMap[sn]; ok {
			parsedPrompt[sn] = prompts.MessageSection{Name: si.Name, Ext: si.Ext, Summary: si.Summary, Data: []byte(sv)}
		} else {
			log.Warn().Str("section", sn).Msg(semLogContext + " - unknown section")
		}
	}

	return parsedPrompt, nil
}

const (
	StatusZero = iota
	StatusInDelimiter

	LineTypeLine        string = "line"
	LineTypeBODelimiter        = "bo-delimiter"
	LineTypeEODelimiter        = "eo-delimiter"
)

/**/

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
