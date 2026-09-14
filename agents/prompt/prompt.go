package prompt

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/agents"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fileutil"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

const (
	semLogPackageContext = "prompt::"

	TemplateVariableTypeFile   = "file"
	TemplateVariableTypeString = "string"
)

type TemplateVariable struct {
	Typ   string `yaml:"type" json:"type" mapstructure:"type"`
	Name  string `yaml:"name" json:"name" mapstructure:"name"`
	Value []byte `yaml:"-" json:"-" mapstructure:"-"`
}
type Location string

func (l Location) resolveFilename(pn, fn string, required bool, withContent bool) (string, []byte, error) {
	const semLogContext = semLogPackageContext + "resolve-filename"
	var err error

	if fn == "" && required {
		err = fmt.Errorf("error: unspecified file %s param", pn)
		log.Error().Err(err).Msg(semLogContext)
		return fn, nil, err
	}

	if fn != "" {
		fn = filepath.Join(string(l), fn)
		if !fileutil.FileExists(fn) {
			err = fmt.Errorf("error: file %s specifies an invalid path", fn)
			log.Error().Err(err).Msg(semLogContext)
			return fn, nil, err
		}

		if fi, err := os.Stat(fn); err != nil || fi.IsDir() {
			err = fmt.Errorf("error: file %s is not a valid file", fn)
			log.Error().Err(err).Msg(semLogContext)
			return fn, nil, err
		}

		if withContent {
			b, err := os.ReadFile(fn)
			if err != nil {
				log.Error().Err(err).Str("path", fn).Msg(semLogContext + " failed to read file")
				return fn, nil, err
			}

			return fn, b, nil
		}
	}

	return fn, nil, nil
}

type Definition struct {
	Name           string             `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	Location       Location           `yaml:"location,omitempty" mapstructure:"location,omitempty" json:"location,omitempty"`
	Vars           []TemplateVariable `yaml:"vars,omitempty" mapstructure:"vars,omitempty" json:"vars,omitempty"`
	SystemFn       string             `yaml:"system-fn" mapstructure:"system-fn" json:"system-fn"`
	TemplateFn     string             `yaml:"template-fn,omitempty" mapstructure:"template-fn,omitempty" json:"template-fn,omitempty"`
	SchemaFn       string             `yaml:"schema-fn,omitempty" mapstructure:"schema-fn,omitempty" json:"schema-fn,omitempty"`
	ParsedTemplate *template.Template `yaml:"-" mapstructure:"-" json:"-"`
	XMLSections    agents.XMLSections `yaml:"xml-sections" mapstructure:"xml-sections" json:"xml-sections"`
}

func ReadPromptDefinition(fn string) (Definition, error) {
	const semLogContext = semLogPackageContext + "read-prompt-definition"

	fileContent, err := os.ReadFile(fn)
	if err != nil {
		log.Error().Err(err).Str("fn", fn).Msg(semLogContext)
		return Definition{}, err
	}

	var def Definition
	err = yaml.Unmarshal(fileContent, &def)
	if err != nil {
		log.Error().Err(err).Str("fn", fn).Msg(semLogContext + " failed to unmarshal prompt template")
		return Definition{}, err
	}

	def.Location = Location(path.Dir(fn))

	def.TemplateFn, fileContent, err = def.Location.resolveFilename("template-fn", def.TemplateFn, true, true)
	if err != nil {
		log.Error().Err(err).Str("fn", fn).Msg(semLogContext + " failed to unmarshal prompt template")
		return Definition{}, err
	}

	tmpl := template.Must(template.New("").Parse(string(fileContent)))
	def.ParsedTemplate = tmpl

	def.SchemaFn, _, err = def.Location.resolveFilename("schema-fn", def.SchemaFn, false, true)
	if err != nil {
		log.Error().Err(err).Str("fn", fn).Msg(semLogContext + " failed to unmarshal prompt template")
		return Definition{}, err
	}

	def.SystemFn, _, err = def.Location.resolveFilename("system-fn", def.SystemFn, true, true)
	if err != nil {
		log.Error().Err(err).Str("fn", fn).Msg(semLogContext + " failed to unmarshal prompt template")
		return Definition{}, err
	}

	return def, nil
}

func (d Definition) OutputType() string {

	switch {
	case d.SchemaFn != "":
		return agents.AgentResponseJSON
	case len(d.XMLSections) > 0:
		return agents.AgentResponseXML
	default:
		return agents.AgentResponseRAW
	}

}

func (d Definition) HasSchemaOutput() bool {
	return d.SchemaFn != ""
}

func (d Definition) SystemPrompt() ([]byte, error) {
	const semLogContext = semLogPackageContext + "system-prompt"
	b, err := os.ReadFile(d.SystemFn)
	if err != nil {
		log.Error().Err(err).Str("fn", d.SystemFn).Msg(semLogContext)
		return nil, err
	}

	return b, nil
}

/*
func (d Definition) UserPrompt(vars map[string]string) ([]byte, error) {
	const semLogContext = semLogPackageContext + "user-prompt"
	var buf bytes.Buffer
	err := d.ParsedTemplate.Execute(&buf, vars)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	return buf.Bytes(), nil
}
*/

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

func (p Definition) UserPrompt(vars map[string]string, forgiving, debugMode bool) ([]byte, error) {
	const semLogContext = semLogPackageContext + "text"
	var err error

	for _, v := range p.Vars {
		if _, ok := vars[v.Name]; !ok {
			if forgiving {
				vars[v.Name] = ""
			} else {
				err = fmt.Errorf("error: variable %s not found in definition", v)
				log.Error().Err(err).Msg(semLogContext)
				return nil, err
			}
		} else {
			if v.Typ == TemplateVariableTypeFile {
				fn, content, err := p.Location.resolveFilename(v.Name, vars[v.Name], true, true)
				if err != nil {
					log.Error().Err(err).Msg(semLogContext)
					return nil, err
				}

				log.Info().Str("fn", fn).Str("variable", v.Name).Msg(semLogContext)
				vars[v.Name] = string(content)
			}
		}
	}

	var buf bytes.Buffer
	err = p.ParsedTemplate.Execute(&buf, vars)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}
	parsedPrompt := buf.Bytes()

	if debugMode {
		fmt.Println(string(parsedPrompt))
	}

	return parsedPrompt, err
}
