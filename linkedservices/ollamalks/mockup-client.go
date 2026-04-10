package ollamalks

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
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

	anthropicMessageParams := Message{
		Model:  c.options.Model,
		Stream: false,
		Prompt: string(b),
	}

	bodyJson, err := json.Marshal(anthropicMessageParams)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	urlBuilder := har.UrlBuilder{}
	urlBuilder.WithScheme(c.cfg.HttpScheme)
	urlBuilder.WithHostname(c.cfg.HostName)
	urlBuilder.WithPort(c.cfg.ServerPort)
	urlBuilder.WithPath(c.cfg.Endpoint)

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

	if harEntry != nil && harEntry.Response != nil {
		if harEntry.Response.Status == http.StatusOK {
			var msg NonStreamedResponse
			err = json.Unmarshal(harEntry.Response.Content.Data, &msg)
			if err != nil {
				log.Error().Err(err).Msg(semLogContext)
				return nil, err
			}

			resp, err := msg.RetrieveContentFromResponse(c.options.Prompt)
			if err != nil {
				log.Error().Err(err).Msg(semLogContext)
				return nil, err
			}

			return &Response{Content: resp}, nil
		} else {
			if harEntry.Response.Content != nil && len(harEntry.Response.Content.Data) > 0 {
				return nil, fmt.Errorf("[%d] - %s", harEntry.Response.Status, string(harEntry.Response.Content.Data))
			}

			return nil, fmt.Errorf("[%d] - %s", harEntry.Response.Status, "no content returned")
		}
	}

	log.Error().Msg(semLogContext + " - no response content")
	return nil, nil
}
