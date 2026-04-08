package anthropiclks

import (
	"net/http"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/rs/zerolog/log"
)

type mockupClient struct {
	cfg        *MockupConfig
	options    ClientOptions
	httpClient *restclient.Client
}

func (c *mockupClient) Close() {
	if c.httpClient != nil {
		c.httpClient.Close()
		c.httpClient = nil
	}
}

func (c *mockupClient) Execute(params ...RequestParam) (*Response, error) {
	const semLogContext = "anthropic-lks-mockup-client::execute"

	Req := Request{}
	for _, p := range params {
		p(&Req)
	}

	b, err := c.options.Prompt.Text(Req.TextVariables())
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	anthropicMessageParams := anthropic.MessageNewParams{
		MaxTokens:   c.options.MaxTokens,
		Temperature: anthropic.Float(c.options.Temperature),
		System:      c.options.Prompt.System,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(string(b))),
		},
		Model: c.options.Model,
	}

	bodyJson, err := anthropicMessageParams.MarshalJSON()
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	urlBuilder := har.UrlBuilder{}
	urlBuilder.WithScheme(c.cfg.HttpScheme)
	urlBuilder.WithHostname(c.cfg.HostName)
	urlBuilder.WithPort(c.cfg.ServerPort)
	urlBuilder.WithPath(c.cfg.Endpoint + "/messages")

	reqHeaders := []har.NameValuePair{{Name: "Content-type", Value: "application/json"}, {Name: "Accept", Value: "application/json"}}
	request, err := c.httpClient.NewRequest(http.MethodPost, urlBuilder.Url(), bodyJson, reqHeaders, nil)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	harEntry, err := c.httpClient.Execute(request)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	if harEntry != nil && harEntry.Response != nil && harEntry.Response.Content != nil {
		var msg anthropic.Message
		err = msg.UnmarshalJSON(harEntry.Response.Content.Data)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}

		resp, err := c.options.Prompt.ParseMessage(&msg)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}

		return &Response{Content: resp}, nil
	}

	log.Error().Msg(semLogContext + " - no response content")
	return nil, nil
}
