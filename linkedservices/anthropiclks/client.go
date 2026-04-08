package anthropiclks

import (
	"context"
	"errors"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type clientImpl struct {
	verbose   bool
	apiClient anthropic.Client
	options   ClientOptions
}

func (c *clientImpl) Close() {
}

func (c *clientImpl) Execute(params ...RequestParam) (*Response, error) {
	const semLogContext = "anthropic-lks-client::execute"
	Req := Request{}
	for _, p := range params {
		p(&Req)
	}

	b, err := c.options.Prompt.Text(Req.TextVariables())
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	message, err := c.apiClient.Messages.New(context.Background(), anthropic.MessageNewParams{
		MaxTokens:   c.options.MaxTokens,
		Temperature: anthropic.Float(c.options.Temperature),
		System:      c.options.Prompt.System,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(string(b))),
		},
		Model: c.options.Model,
	})

	if message != nil && c.verbose {
		fmt.Println("-------------------")
		fmt.Println(message.RawJSON())
		fmt.Println("-------------------")
	}

	if err != nil {
		logError(err).Msg(semLogContext)
		return nil, err
	}

	if message == nil {
		return nil, errors.New("no message returned")
	}

	resp, err := c.options.Prompt.ParseMessage(message)
	if err != nil {
		logError(err).Msg(semLogContext)
		return nil, err
	}
	return &Response{Content: resp}, err
}

func logError(err error) *zerolog.Event {
	var evt *zerolog.Event
	if err != nil {
		if apiErr, ok := errors.AsType[*anthropic.Error](err); ok {
			evt = log.Error().Err(apiErr).Str("Request ID:", apiErr.RequestID).Int("status-code", apiErr.StatusCode)
			println(string(apiErr.DumpRequest(true)))  // Prints the serialized HTTP request
			println(string(apiErr.DumpResponse(true))) // Prints the serialized HTTP response
		} else {
			evt = log.Error().Err(err)
		}
	} else {
		evt = log.Info()
	}

	return evt
}
