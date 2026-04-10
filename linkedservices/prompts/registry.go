package prompts

import (
	"embed"
	"fmt"
	"text/template"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fileutil"
	"github.com/rs/zerolog/log"
)

const semLogContextBasePromptRegistry = "prompt-registry::"

type Registry struct {
	templates map[string]*template.Template
	prompts   map[string]PromptTemplate
}

var theRegistry Registry

func NewPromptsRegistry(rootFolder string, embeddedTemplates embed.FS) error {
	const semLogContext = semLogContextBasePromptRegistry + "new"

	fs, err := fileutil.FindEmbeddedFiles(embeddedTemplates, rootFolder,
		fileutil.WithFindOptionNavigateSubDirs(),
		fileutil.WithFindFileType(fileutil.FileTypeFile),
		fileutil.WithFindOptionExcludeRootFolderInNames())
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	theRegistry = Registry{
		templates: make(map[string]*template.Template, len(fs)),
	}
	for _, ff := range fs {
		tmpl := template.Must(template.New("").Parse(string(ff.Content)))
		theRegistry.templates[ff.Info.Name()] = tmpl
		log.Info().Str("path", ff.Path).Msg(semLogContext + " prompt template loaded")
	}

	return nil
}

func RegisterPrompt(aPrompt PromptTemplate) error {
	const semLogContext = semLogContextBasePromptRegistry + "register-prompt"
	var err error

	if aPrompt.Data == nil {
		var ok bool
		aPrompt.Data, ok = theRegistry.templates[aPrompt.TemplateName]
		if !ok {
			err = fmt.Errorf("template %s not found", aPrompt.TemplateName)
			log.Error().Err(err).Msg(semLogContext)
			return err
		}
	}

	if theRegistry.prompts == nil {
		theRegistry.prompts = make(map[string]PromptTemplate)
	}

	theRegistry.prompts[aPrompt.Name] = aPrompt
	return nil
}

func GetPrompt(name string) (PromptTemplate, error) {
	const semLogContext = semLogContextBasePromptRegistry + "get-prompt"
	pt, ok := theRegistry.prompts[name]
	if !ok {
		err := fmt.Errorf("prompt %s not found", name)
		log.Error().Err(err).Msg(semLogContext)
		return PromptTemplate{}, err
	}

	return pt, nil
}
