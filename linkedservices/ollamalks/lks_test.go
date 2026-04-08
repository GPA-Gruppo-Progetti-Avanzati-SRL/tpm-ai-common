package ollamalks_test

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/ollamalks"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/ollamalks/prompts"
	"github.com/stretchr/testify/require"
)

//go:embed example.cob
var cobExample []byte

func TestClient(t *testing.T) {

	cfg := &ollamalks.Config{
		ApiKey:         "",
		MaxRetries:     0,
		RequestTimeout: time.Second * 120,
		Mockup: &ollamalks.MockupConfig{
			Enabled:    true,
			HostName:   "localhost",
			ServerPort: 11434,
			HttpScheme: "http",
			Endpoint:   "/api/generate",
		},
		Verbose:       true,
		ClientOptions: ollamalks.ClientOptions{},
	}

	lks, err := ollamalks.Initialize(cfg)
	require.NoError(t, err)

	prompt, err := prompts.GetPrompt(prompts.PromptProgramSummary)
	require.NoError(t, err)

	cli, err := lks.NewClient(
		ollamalks.WithPrompt(prompt),
		ollamalks.WithModel("granite4:latest"),
	)
	require.NoError(t, err)

	defer cli.Close()

	resp, err := cli.Execute(ollamalks.WithPromptInput(prompts.TextVariable, "COBOL_SOURCE", 0, cobExample))
	require.NoError(t, err)
	require.NotEmpty(t, resp)

	for _, lr := range resp.Content {
		t.Log(lr.Name)
		fp := filepath.Join("/tmp", fmt.Sprintf("%s-%s%s", "granite4", lr.Name, lr.Ext))
		require.NoError(t, os.WriteFile(fp, []byte(lr.Data), os.ModePerm))
	}
	t.Log(resp)
}
