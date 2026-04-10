package ollamalks_test

import (
	"embed"
	"os"
	"testing"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/prompts"
	"github.com/rs/zerolog/log"
)

//go:embed templates
var templates embed.FS

const (
	TemplatesRootFolder  = "templates"
	PromptProgramSummary = "program-summary-prompt"

	ScratchPadSection = "scratchpad"
	OverviewSection   = "overview"
	SummarySection    = "summary"
	MermaidSection    = "mermaid"

	ContentTypeTextPlain             = "text/plain"
	ContentTypeTextMarkdown          = "text/markdown"
	ContentTypeApplicationVndMermaid = "application/vnd.mermaid"
)

func TestMain(m *testing.M) {

	err := prompts.NewPromptsRegistry(TemplatesRootFolder, templates)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create registry")
	}

	err = prompts.RegisterPrompt(prompts.PromptTemplate{
		Name:         PromptProgramSummary,
		TemplateName: "program-summary-prompt.txt",
		System:       "You are a senior Cobol developer with extensive knowledge in mainframe cobol cics programming",
		Vars:         []string{"COBOL_SOURCE"},
		Sections: []prompts.MessagePart{
			{Name: ScratchPadSection, Ext: ".md", Summary: false, Ct: ContentTypeTextMarkdown},
			{Name: OverviewSection, Ext: ".md", Summary: false, Ct: ContentTypeTextMarkdown},
			{Name: SummarySection, Ext: ".txt", Summary: true, Ct: ContentTypeTextPlain},
			{Name: MermaidSection, Ext: ".mmd", Summary: false, Ct: ContentTypeApplicationVndMermaid},
		},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register prompt")
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}
