package anthropiclks_test

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/anthropiclks"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/anthropiclks/prompts"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

//go:embed example.cob
var cobExample []byte

func TestClient(t *testing.T) {

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	require.NotEmpty(t, apiKey)

	cfg := &anthropiclks.Config{
		ApiKey:         apiKey,
		MaxRetries:     0,
		RequestTimeout: time.Second * 480,
		Verbose:        true,
	}

	lks, err := anthropiclks.Initialize(cfg)
	require.NoError(t, err)

	prompt, err := prompts.GetPrompt(prompts.PromptProgramSummary)
	require.NoError(t, err)

	cli, err := lks.NewClient(
		anthropiclks.WithPrompt(prompt),
	)
	require.NoError(t, err)

	defer cli.Close()

	resp, err := cli.Execute(anthropiclks.WithPromptInput(prompts.TextVariable, "COBOL_SOURCE", 0, cobExample))
	require.NoError(t, err)

	for _, lr := range resp.Content {
		t.Log(lr.Name)
		fp := filepath.Join("/tmp", fmt.Sprintf("%s-%s%s", "opus-4.6", lr.Name, lr.Ext))
		require.NoError(t, os.WriteFile(fp, []byte(lr.Data), os.ModePerm))
	}
	t.Log(resp)
}

func TestMockupClient(t *testing.T) {

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	require.NotEmpty(t, apiKey)

	cfg := &anthropiclks.Config{
		ApiKey:         apiKey,
		MaxRetries:     0,
		RequestTimeout: time.Second * 20,
		Mockup: &anthropiclks.MockupConfig{
			HtttpClient: &restclient.Config{
				RestTimeout:       0,
				SkipVerify:        true,
				Headers:           nil,
				TraceGroupName:    "",
				TraceRequestName:  "",
				RetryCount:        0,
				RetryWaitTime:     0,
				RetryMaxWaitTime:  0,
				RetryOnHttpError:  nil,
				HarTracingEnabled: false,
				Span:              nil,
				HarSpan:           nil,
			},
			Enabled:    true,
			HostName:   "localhost",
			ServerPort: 3001,
			HttpScheme: "http",
			Endpoint:   "/anthropic/test",
		},
		Verbose: false,
	}

	lks, err := anthropiclks.Initialize(cfg)
	require.NoError(t, err)

	b, err := yaml.Marshal(lks.Cfg)
	require.NoError(t, err)

	fmt.Println(string(b))

	prompt, err := prompts.GetPrompt(prompts.PromptProgramSummary)
	require.NoError(t, err)

	cli, err := lks.NewClient(
		anthropiclks.WithPrompt(prompt),
	)
	require.NoError(t, err)

	defer cli.Close()

	resp, err := cli.Execute(anthropiclks.WithPromptInput(prompts.TextVariable, "COBOL_SOURCE", 0, []byte("this is a source!")))
	require.NoError(t, err)

	t.Log(resp)
}
