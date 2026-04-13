package anthropiclks

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/prompts"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/rs/zerolog/log"
)

type mockupBatchClient struct {
	cfg        *MockupConfig
	httpClient *restclient.Client
}

func (c *mockupBatchClient) Close() {
	if c.httpClient != nil {
		c.httpClient.Close()
		c.httpClient = nil
	}
}

// SubmitBatch renders each request through the prompt template, builds a
// MessageBatchNewParams body, POSTs it to {endpoint}/batches, and returns the
// batch ID from the MessageBatch response.
func (c *mockupBatchClient) SubmitBatch(requests []BatchRequest) (string, error) {
	const semLogContext = "anthropic-lks-mockup-batch-client::submit-batch"

	batchRequests := make([]anthropic.MessageBatchNewParamsRequest, 0, len(requests))
	for _, req := range requests {

		if req.CustomID == "" {
			err := fmt.Errorf("custom-id is required")
			log.Error().Err(err).Msg(semLogContext)
			return "", err
		}

		r := Request{}
		for _, p := range req.Params {
			p(&r)
		}

		b, err := r.Prompt.Text(r.TextVariables())
		if err != nil {
			log.Error().Err(err).Str("custom-id", req.CustomID).Msg(semLogContext)
			return "", err
		}

		batchRequests = append(batchRequests, anthropic.MessageBatchNewParamsRequest{
			CustomID: req.CustomID,
			Params: anthropic.MessageBatchNewParamsRequestParams{
				MaxTokens:   r.MaxTokens,
				Temperature: anthropic.Float(r.Temperature),
				System:      []anthropic.TextBlockParam{{Text: r.Prompt.System}},
				Messages: []anthropic.MessageParam{
					anthropic.NewUserMessage(anthropic.NewTextBlock(string(b))),
				},
				Model: r.Model,
			},
		})
	}

	bodyJson, err := anthropic.MessageBatchNewParams{Requests: batchRequests}.MarshalJSON()
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return "", err
	}

	urlBuilder := har.UrlBuilder{}
	urlBuilder.WithScheme(c.cfg.HttpScheme)
	urlBuilder.WithHostname(c.cfg.HostName)
	urlBuilder.WithPort(c.cfg.ServerPort)
	urlBuilder.WithPath(c.cfg.Endpoint + "/batches")

	reqHeaders := []har.NameValuePair{{Name: "Content-type", Value: "application/json"}, {Name: "Accept", Value: "application/json"}}
	request, err := c.httpClient.NewRequest(http.MethodPost, urlBuilder.Url(), bodyJson, reqHeaders, nil)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return "", err
	}

	harEntry, err := c.httpClient.Execute(request)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return "", err
	}

	if harEntry != nil && harEntry.Response != nil && harEntry.Response.Content != nil {
		var batch anthropic.MessageBatch
		if err = batch.UnmarshalJSON(harEntry.Response.Content.Data); err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return "", err
		}
		return batch.ID, nil
	}

	log.Error().Msg(semLogContext + " - no response content")
	return "", nil
}

// GetBatchStatus GETs {endpoint}/batches/{batchID} and maps the MessageBatch
// response to a BatchStatus.
func (c *mockupBatchClient) GetBatchStatus(batchID string) (*BatchStatus, error) {
	const semLogContext = "anthropic-lks-mockup-batch-client::get-batch-status"

	urlBuilder := har.UrlBuilder{}
	urlBuilder.WithScheme(c.cfg.HttpScheme)
	urlBuilder.WithHostname(c.cfg.HostName)
	urlBuilder.WithPort(c.cfg.ServerPort)
	urlBuilder.WithPath(c.cfg.Endpoint + "/batches/" + batchID)

	reqHeaders := []har.NameValuePair{{Name: "Accept", Value: "application/json"}}
	request, err := c.httpClient.NewRequest(http.MethodGet, urlBuilder.Url(), nil, reqHeaders, nil)
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
		var batch anthropic.MessageBatch
		if err = batch.UnmarshalJSON(harEntry.Response.Content.Data); err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}

		return &BatchStatus{
			ID:               batch.ID,
			ProcessingStatus: string(batch.ProcessingStatus),
			RequestCounts: BatchRequestCounts{
				Processing: batch.RequestCounts.Processing,
				Succeeded:  batch.RequestCounts.Succeeded,
				Errored:    batch.RequestCounts.Errored,
				Canceled:   batch.RequestCounts.Canceled,
				Expired:    batch.RequestCounts.Expired,
			},
			ExpiresAt: batch.ExpiresAt,
		}, nil
	}

	log.Error().Msg(semLogContext + " - no response content")
	return nil, nil
}

// GetBatchResults GETs {endpoint}/batches/{batchID}/results, expects a JSON
// array of MessageBatchIndividualResponse, and parses each succeeded item
// through the prompt template the same way Execute does.
func (c *mockupBatchClient) GetBatchResults(batchID string) ([]BatchResult, error) {
	const semLogContext = "anthropic-lks-mockup-batch-client::get-batch-results"

	urlBuilder := har.UrlBuilder{}
	urlBuilder.WithScheme(c.cfg.HttpScheme)
	urlBuilder.WithHostname(c.cfg.HostName)
	urlBuilder.WithPort(c.cfg.ServerPort)
	urlBuilder.WithPath(c.cfg.Endpoint + "/batches/" + batchID + "/results")

	reqHeaders := []har.NameValuePair{{Name: "Accept", Value: "application/json"}}
	request, err := c.httpClient.NewRequest(http.MethodGet, urlBuilder.Url(), nil, reqHeaders, nil)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	harEntry, err := c.httpClient.Execute(request)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	if harEntry == nil || harEntry.Response == nil || harEntry.Response.Content == nil {
		log.Error().Msg(semLogContext + " - no response content")
		return nil, nil
	}

	var items []anthropic.MessageBatchIndividualResponse
	if err = json.Unmarshal(harEntry.Response.Content.Data, &items); err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	results := make([]BatchResult, 0, len(items))
	for _, item := range items {
		br := BatchResult{CustomID: item.CustomID}

		var promptName string
		var ok bool
		if _, _, promptName, ok = BatchRequestCustomIdWellFormed(item.CustomID); !ok {
			return nil, fmt.Errorf("custom-id is not well formed")
		}

		prompt, err := prompts.GetPrompt(promptName)
		if err != nil {
			logError(err).Msg(semLogContext)
			return nil, err

		}

		switch item.Result.Type {
		case "succeeded":
			succeeded := item.Result.AsSucceeded()
			resp, err := ParseMessage(prompt, &succeeded.Message)
			if err != nil {
				log.Error().Err(err).Str("custom-id", item.CustomID).Msg(semLogContext)
				br.Err = err
			} else {
				br.Response = &Response{Content: resp}
			}
		case "errored":
			errored := item.Result.AsErrored()
			br.Err = fmt.Errorf("batch item errored: %s", errored.Error.Error.Message)
		case "canceled":
			br.Err = fmt.Errorf("batch item canceled")
		case "expired":
			br.Err = fmt.Errorf("batch item expired")
		}

		results = append(results, br)
	}

	return results, nil
}
