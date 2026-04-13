package prompts_test

import (
	"embed"
	"testing"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/prompts"
	"github.com/stretchr/testify/require"
)

//go:embed prmpts_repo
var templates embed.FS

const (
	TemplatesRootFolder = "prmpts_repo"
)

func TestRegistry(t *testing.T) {
	err := prompts.NewPromptsRegistry(TemplatesRootFolder, templates)
	require.NoError(t, err)
}
