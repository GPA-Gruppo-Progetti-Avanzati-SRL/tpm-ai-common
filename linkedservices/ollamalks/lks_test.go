package ollamalks_test

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/ollamalks"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/prompts"
	"github.com/stretchr/testify/require"
)

//go:embed example.cob
var cobExample []byte

func TestClient(t *testing.T) {

	cfg := &ollamalks.Config{
		Token:   "",
		Url:     "http://localhost:11434",
		Verbose: true,
	}

	lks, err := ollamalks.Initialize(cfg)
	require.NoError(t, err)

	cat, err := prompts.GetCategory("ollama:program-summary:granite4:latest")
	require.NoError(t, err)

	prompt, err := prompts.GetPrompt(cat.Prompt)
	require.NoError(t, err)

	cli, err := lks.NewClient()
	require.NoError(t, err)

	defer cli.Close()

	resp, err := cli.Execute(
		ollamalks.WithPrompt(prompt),
		ollamalks.WithPromptInput(prompts.TextVariable, "COBOL_SOURCE", 0, cobExample),
		ollamalks.WithModel(cat.Model),
		ollamalks.WithGenerateOptions(ollamalks.GenerateOptionsFromCategoryOptions(&cat)),
	)
	require.NoError(t, err)
	require.NotEmpty(t, resp)

	for _, lr := range resp.Content {
		t.Log(lr.Name)
		fp := filepath.Join("/tmp", fmt.Sprintf("%s-%s%s", "granite4", lr.Name, lr.Ext))
		require.NoError(t, os.WriteFile(fp, []byte(lr.Data), os.ModePerm))
	}
	t.Log(resp)
}
