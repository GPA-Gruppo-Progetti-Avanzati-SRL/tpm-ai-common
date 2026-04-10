package ollamalks

import (
	"errors"
	"fmt"
	"strings"

	"github.com/eslider/go-ollama"
	"github.com/rs/zerolog/log"
)

type mockupClient struct {
	cfg        *MockupConfig
	options    ClientOptions
	httpClient *ollama.Client
}

func (c *mockupClient) Close() {
	if c.httpClient != nil {
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

	//b1, err := json.Marshal(string(b))
	//fmt.Println(string(b1))

	var full strings.Builder
	req := ollama.Request{
		Model:  c.options.Model,
		Prompt: string(b),
		OnJson: func(res ollama.Response) error {
			// log.Info().Msg(semLogContext + " - on json")
			if res.Response != nil {
				full.WriteString(*res.Response)
			}
			return nil
		},
		Options: &ollama.RequestOptions{
			NumContext: new(100000),
		},
	}

	err = c.httpClient.Query(req)
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
