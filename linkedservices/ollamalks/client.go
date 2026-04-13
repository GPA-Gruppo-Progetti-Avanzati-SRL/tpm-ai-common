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

type clientImpl struct {
	verbose      bool
	ollamaClient *api.Client
}

func (c *clientImpl) Close() {
	if c.ollamaClient != nil {
		c.ollamaClient = nil
	}
}

func (c *clientImpl) Execute(params ...RequestParam) (*Response, error) {
	const semLogContext = semLogContextBaseLks + "execute"

	reqParams := Request{}
	for _, p := range params {
		p(&reqParams)
	}

	b, err := reqParams.Prompt.Text(reqParams.TextVariables())
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	//b1, err := json.Marshal(string(b))
	//fmt.Println(string(b1))

	reqOpts := make(map[string]interface{})
	if reqParams.GenerateOptions.NumCtx > 0 {
		reqOpts[GenerateOptionsNumContext] = reqParams.GenerateOptions.NumCtx
	}
	if reqParams.GenerateOptions.Temperature > 0 {
		reqOpts[GenerateOptionsTemperature] = reqParams.GenerateOptions.Temperature
	}
	if len(reqOpts) == 0 {
		reqOpts = nil
	}

	var responseFormat json.RawMessage
	if reqParams.Prompt.TemplateSchema != "" {
		responseFormat = json.RawMessage(reqParams.Prompt.TemplateSchema)
	}

	req := &api.GenerateRequest{
		Model:  reqParams.Model,
		Prompt: string(b),

		// set streaming to false
		Stream:  new(bool),
		Options: reqOpts,
		Format:  responseFormat,
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

	parsedPrompt, err := reqParams.Prompt.ParseTextContent(full.String())
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	return &Response{Content: parsedPrompt}, nil
}
