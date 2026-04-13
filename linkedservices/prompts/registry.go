package prompts

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fileutil"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

const semLogContextBasePromptRegistry = "prompt-registry::"

type Registry struct {
	prompts    map[string]PromptTemplate
	categories map[string]Category
}

var theRegistry Registry

func NewPromptsRegistry(fld string) error {
	const semLogContext = semLogContextBasePromptRegistry + "new"

	fs, err := fileutil.FindFiles(fld,
		fileutil.WithFindOptionNavigateSubDirs(),
		fileutil.WithFindFileType(fileutil.FileTypeFile),
		fileutil.WithFindOptionExcludeRootFolderInNames(),
	)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	theRegistry = Registry{
		prompts:    make(map[string]PromptTemplate, len(fs)),
		categories: make(map[string]Category, len(fs)),
	}
	for _, fn := range fs {

		switch {
		case strings.HasPrefix(fn, "categories"):
			err = addCategory(fn)
		case strings.HasSuffix(fn, ".yml"):
			err = addPrompt(fn)
		}
		if err != nil {
			log.Error().Err(err).Str("fn", fn).Msg(semLogContext)
			return err
		}
	}

	return nil
}

func addCategory(fn string) error {
	const semLogContext = semLogContextBasePromptRegistry + "add-category"

	content, err := os.ReadFile(fn)
	if err != nil {
		log.Error().Err(err).Str("fn", fn).Msg(semLogContext)
		return err
	}

	var cat []Category
	err = yaml.Unmarshal(content, &cat)
	if err != nil {
		log.Error().Err(err).Str("fn", fn).Msg(semLogContext + " failed to unmarshal prompt template")
		return err
	}

	for _, c := range cat {
		theRegistry.categories[c.Name] = c
	}

	log.Info().Str("fn", fn).Msg(semLogContext + " category loaded")
	return nil
}

func addPrompt(fn string) error {
	const semLogContext = semLogContextBasePromptRegistry + "add-prompt"

	prmpt, err := NewPrompt(fn)
	if err != nil {
		log.Error().Err(err).Str("fn", fn).Msg(semLogContext + " failed to unmarshal prompt template")
		return err
	}

	theRegistry.prompts[prmpt.Name] = prmpt
	log.Info().Str("fn", fn).Msg(semLogContext + " prompt template loaded")
	return nil
}

func NewPromptsEmbeddedRegistry(rootFolder string, embeddedTemplates embed.FS) error {
	const semLogContext = semLogContextBasePromptRegistry + "new-embedded"

	fs, err := fileutil.FindEmbeddedFiles(embeddedTemplates, rootFolder,
		fileutil.WithFindOptionNavigateSubDirs(),
		fileutil.WithFindFileType(fileutil.FileTypeFile),
		fileutil.WithFindOptionExcludeRootFolderInNames(),
		fileutil.WithFindOptionPreloadContent())
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	theRegistry = Registry{
		prompts:    make(map[string]PromptTemplate, len(fs)),
		categories: make(map[string]Category, len(fs)),
	}
	for _, ff := range fs {
		fn := ff.Info.Name()
		switch {
		case strings.HasPrefix(fn, "categories"):
			err = addEmbeddedCategory(rootFolder, embeddedTemplates, fn, ff.Content)
		case strings.HasSuffix(fn, ".yml"):
			err = addEmbeddedPrompt(rootFolder, embeddedTemplates, fn, ff.Content)
		}
		if err != nil {
			log.Error().Err(err).Str("fn", fn).Msg(semLogContext)
			return err
		}
	}

	return nil
}

func addEmbeddedCategory(rootFolder string, embeddedTemplates embed.FS, fn string, fileContent []byte) error {
	const semLogContext = semLogContextBasePromptRegistry + "add-embedded-category"
	var cat []Category
	err := yaml.Unmarshal(fileContent, &cat)
	if err != nil {
		log.Error().Err(err).Str("fn", fn).Msg(semLogContext + " failed to unmarshal prompt template")
		return err
	}

	for _, c := range cat {
		theRegistry.categories[c.Name] = c
	}

	log.Info().Str("fn", fn).Msg(semLogContext + " category loaded")
	return nil
}

func addEmbeddedPrompt(rootFolder string, embeddedTemplates embed.FS, fn string, fileContent []byte) error {
	const semLogContext = semLogContextBasePromptRegistry + "add-embedded-prompt"

	prmpt, err := NewPromptFromEmbeddedFS(rootFolder, embeddedTemplates, fn, fileContent)
	if err != nil {
		log.Error().Err(err).Str("fn", fn).Msg(semLogContext + " failed to unmarshal prompt template")
		return err
	}

	theRegistry.prompts[prmpt.Name] = prmpt
	log.Info().Str("fn", fn).Msg(semLogContext + " prompt template loaded")
	return nil
}

func RegisterPrompt(aPrompt PromptTemplate) error {
	const semLogContext = semLogContextBasePromptRegistry + "register-prompt"

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

func GetCategory(name string) (Category, error) {
	const semLogContext = semLogContextBasePromptRegistry + "get-category"
	cat, ok := theRegistry.categories[name]
	if !ok {
		err := fmt.Errorf("category %s not found", name)
		log.Error().Err(err).Msg(semLogContext)
		return Category{}, err
	}

	return cat, nil
}
