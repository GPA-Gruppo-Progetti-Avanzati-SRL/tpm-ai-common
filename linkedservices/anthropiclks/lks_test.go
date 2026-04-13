package anthropiclks_test

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/anthropiclks"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/prompts"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/stretchr/testify/require"
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

	cat, err := prompts.GetCategory("claude:program-summary:claude-opus-4-6")
	require.NoError(t, err)

	prompt, err := prompts.GetPrompt(cat.Prompt)
	require.NoError(t, err)

	cli, err := lks.NewClient()
	require.NoError(t, err)

	defer cli.Close()

	resp, err := cli.Execute(
		anthropiclks.WithPrompt(prompt),
		anthropiclks.WithModel(cat.Model),
		anthropiclks.WithPromptInput(prompts.TextVariable, "COBOL_SOURCE", 0, cobExample),
		anthropiclks.WithTemperature(cat.GetFloat64Option("temperature", 0.2)),
		anthropiclks.WithMaxTokens(cat.GetInt64Option("max-tokens", 20000)),
	)
	require.NoError(t, err)

	for _, lr := range resp.Content {
		t.Log(lr.Name)
		fp := filepath.Join("/tmp", fmt.Sprintf("%s-%s%s", "opus-4.6", lr.Name, lr.Ext))
		require.NoError(t, os.WriteFile(fp, []byte(lr.Data), os.ModePerm))
	}
	t.Log(resp)
}

func TestBatchClient(t *testing.T) {

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

	cat, err := prompts.GetCategory("claude:program-summary:claude-opus-4-6")
	require.NoError(t, err)

	prompt, err := prompts.GetPrompt(cat.Prompt)
	require.NoError(t, err)

	cli, err := lks.NewBatchClient()
	require.NoError(t, err)

	defer cli.Close()

	resp, err := cli.SubmitBatch([]anthropiclks.BatchRequest{
		{
			CustomID: anthropiclks.BatchRequestCustomId("program", "YPBCEPGP", cat.Prompt),
			Params: []anthropiclks.RequestParam{
				anthropiclks.WithPrompt(prompt),
				anthropiclks.WithModel(cat.Model),
				anthropiclks.WithPromptInput(prompts.TextVariable, "COBOL_SOURCE", 0, cobExample),
				anthropiclks.WithTemperature(cat.GetFloat64Option("temperature", 0.2)),
				anthropiclks.WithMaxTokens(cat.GetInt64Option("max-tokens", 20000)),
			},
		},
	})
	require.NoError(t, err)

	bt, err := cli.GetBatchStatus(resp)
	require.NoError(t, err)

	for bt.ProcessingStatus != anthropiclks.BatchStatusEnded {
		time.Sleep(time.Second * 60)
		bt, err = cli.GetBatchStatus(resp)
		require.NoError(t, err)
	}

	brs, err := cli.GetBatchResults(resp)
	require.NoError(t, err)

	for _, br := range brs {
		for _, mp := range br.Response.Content {
			t.Log(mp.Name, mp.Ct, mp.Ext)
			t.Log(string(mp.Data))
		}
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
			Endpoint:   "/api/anthropic",
		},
		Verbose: false,
	}

	lks, err := anthropiclks.Initialize(cfg)
	require.NoError(t, err)

	//b, err := yaml.Marshal(lks.Cfg)
	//require.NoError(t, err)
	//fmt.Println(string(b))

	cat, err := prompts.GetCategory("claude:program-summary:claude-opus-4-6")
	require.NoError(t, err)

	prompt, err := prompts.GetPrompt(cat.Prompt)
	require.NoError(t, err)

	cli, err := lks.NewClient()
	require.NoError(t, err)

	defer cli.Close()

	resp, err := cli.Execute(
		anthropiclks.WithPrompt(prompt),
		anthropiclks.WithModel(cat.Model),
		anthropiclks.WithPromptInput(prompts.TextVariable, "COBOL_SOURCE", 0, []byte("this is a source!")),
		anthropiclks.WithTemperature(cat.GetFloat64Option("temperature", 0.2)),
		anthropiclks.WithMaxTokens(cat.GetInt64Option("max-tokens", 20000)))
	require.NoError(t, err)

	for _, mp := range resp.Content {
		t.Log(mp.Name, mp.Ct, mp.Ext)
		t.Log(string(mp.Data))
	}
}

func TestMockupBatchClient(t *testing.T) {

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
			Endpoint:   "/api/anthropic",
		},
		Verbose: false,
	}

	lks, err := anthropiclks.Initialize(cfg)
	require.NoError(t, err)

	//b, err := yaml.Marshal(lks.Cfg)
	//require.NoError(t, err)
	//fmt.Println(string(b))

	cat, err := prompts.GetCategory("claude:program-summary:claude-opus-4-6")
	require.NoError(t, err)

	prompt, err := prompts.GetPrompt(cat.Prompt)
	require.NoError(t, err)

	cli, err := lks.NewBatchClient()
	require.NoError(t, err)

	defer cli.Close()

	batchRequest := anthropiclks.BatchRequest{
		CustomID: anthropiclks.BatchRequestCustomId("program", "xytcunic", cat.Prompt),
		Params: []anthropiclks.RequestParam{
			anthropiclks.WithPrompt(prompt),
			anthropiclks.WithModel(cat.Model),
			anthropiclks.WithPromptInput(prompts.TextVariable, "COBOL_SOURCE", 0, []byte("this is a source!")),
			anthropiclks.WithTemperature(cat.GetFloat64Option("temperature", 0.2)),
			anthropiclks.WithMaxTokens(cat.GetInt64Option("max-tokens", 20000)),
		},
	}

	resp, err := cli.SubmitBatch([]anthropiclks.BatchRequest{batchRequest})
	require.NoError(t, err)
	t.Log(resp)

	brs, err := cli.GetBatchResults(resp)
	require.NoError(t, err)

	for _, br := range brs {
		for _, mp := range br.Response.Content {
			t.Log(mp.Name, mp.Ct, mp.Ext)
			t.Log(string(mp.Data))
		}
	}

}
