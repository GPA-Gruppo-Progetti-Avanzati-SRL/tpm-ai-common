package anthropiclks

import (
	"context"
	"fmt"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/prompts"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/rs/zerolog/log"
)

type batchClientImpl struct {
	verbose   bool
	apiClient anthropic.Client
}

func (c *batchClientImpl) Close() {}

func (c *batchClientImpl) SubmitBatch(requests []BatchRequest) (string, error) {
	const semLogContext = "anthropic-lks-batch-client::submit-batch"

	batchRequests := make([]anthropic.MessageBatchNewParamsRequest, 0, len(requests))
	for _, req := range requests {

		if _, _, _, ok := req.CustomIdWellFormed(); !ok {
			return "", fmt.Errorf("custom-id is not well formed")
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

	batch, err := c.apiClient.Messages.Batches.New(context.Background(), anthropic.MessageBatchNewParams{
		Requests: batchRequests,
	})
	if err != nil {
		logError(err).Msg(semLogContext)
		return "", err
	}

	if c.verbose {
		fmt.Println("Request -------------------")
		fmt.Println(batch.RawJSON())
		fmt.Println("-------------------")
	}

	return batch.ID, nil
}

func (c *batchClientImpl) GetBatchStatus(batchID string) (*BatchStatus, error) {
	const semLogContext = "anthropic-lks-batch-client::get-batch-status"

	batch, err := c.apiClient.Messages.Batches.Get(context.Background(), batchID)
	if err != nil {
		logError(err).Msg(semLogContext)
		return nil, err
	}

	if c.verbose {
		fmt.Println("Status -------------------")
		fmt.Println(batch.RawJSON())
		fmt.Println("-------------------")
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

func (c *batchClientImpl) GetBatchResults(batchID string) ([]BatchResult, error) {
	const semLogContext = "anthropic-lks-batch-client::get-batch-results"

	stream := c.apiClient.Messages.Batches.ResultsStreaming(context.Background(), batchID)
	defer stream.Close()

	var results []BatchResult
	for stream.Next() {
		item := stream.Current()

		if c.verbose {
			fmt.Println("Stream Result: -------------------")
			fmt.Println(item.RawJSON())
			fmt.Println("-------------------")
		}

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

		br := BatchResult{CustomID: item.CustomID}

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

	if err := stream.Err(); err != nil {
		logError(err).Msg(semLogContext)
		return nil, err
	}

	return results, nil
}
