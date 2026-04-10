package ollamalks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/rs/zerolog/log"
)

type mockupClient struct {
	cfg          *MockupConfig
	options      RequestOptions
	ollamaClient *api.Client
}

func (c *mockupClient) Close() {
	if c.ollamaClient != nil {
		c.ollamaClient = nil
	}
}

var schema = []byte(`
{
    "type": "object",
    "properties": {
      "scratchpad": {
        "type": "string"
      },
      "summary": {
        "type": "string"
      },
      "overview": {
        "type": "string"
      },
      "mermaid": {
        "type": "string"
      },
    },
    "required": [
      "scratchpad",
      "summary", 
      "overview",
      "mermaid"
    ]
}	
`)

func (c *mockupClient) Execute(params ...RequestParam) (*Response, error) {
	const semLogContext = semLogContextBaseLks + "execute"

	Req := Request{}
	for _, p := range params {
		p(&Req)
	}

	b, err := c.options.Prompt.Text(Req.TextVariables())
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	//b1, err := json.Marshal(string(b))
	//fmt.Println(string(b1))

	reqOpts := make(map[string]interface{})
	if c.options.GenerateOptions.NumCtx > 0 {
		reqOpts["num_ctx"] = c.options.GenerateOptions.NumCtx
	}
	if c.options.GenerateOptions.Temperature > 0 {
		reqOpts["temperature"] = c.options.GenerateOptions.Temperature
	}
	if len(reqOpts) == 0 {
		reqOpts = nil
	}

	req := &api.GenerateRequest{
		Model:  c.options.Model,
		Prompt: string(b),

		// set streaming to false
		Stream:  new(bool),
		Options: reqOpts,
		Format:  make(json.RawMessage, 0),
	}

	var full strings.Builder
	respFunc := func(resp api.GenerateResponse) error {
		// Only print the response here; GenerateResponse has a number of other
		// interesting fields you want to examine.
		full.WriteString(resp.Response)
		return nil
	}

	err = c.ollamaClient.Generate(context.Background(), req, respFunc)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	text := full.String()
	fmt.Println(text)
	if len(text) == 0 {
		err = errors.New("no response content")
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	parsedPrompt, err := c.options.Prompt.ParseTextContent(full.String())
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	return &Response{Content: parsedPrompt}, nil
}
