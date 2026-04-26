package anthropiclks

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/prompts"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/rs/zerolog/log"
)

const (
	semLogContextBase = "anthropic-lks::"
	LinkedServiceType = "claude"
)

type LinkedService struct {
	Cfg           *Config
	httpClientLks *restclient.LinkedService
}

var theLks LinkedService

func Initialize(cfg *Config) (*LinkedService, error) {
	const semLogContext = "anthropic-lks-registry::initialize"
	var err error

	if cfg == nil {
		log.Info().Msg(semLogContext + " no config provided....skipping")
		return nil, nil
	}

	if theLks.Cfg != nil {
		log.Warn().Msg(semLogContext + " registry already configured.. overwriting")
	}

	theLks, err = newInstanceWithConfig(cfg)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
	}

	return &theLks, nil
}

func newInstanceWithConfig(cfg *Config) (LinkedService, error) {
	var err error

	lks := LinkedService{Cfg: cfg}

	if cfg.Mockup != nil && cfg.Mockup.Enabled {
		lks.httpClientLks, err = restclient.NewInstanceWithConfig(cfg.Mockup.HtttpClient)
	}

	return lks, err
}

type Client interface {
	Close()
	Execute(params ...RequestParam) (*Response, error)
}

type BatchClient interface {
	Close()
	SubmitBatch(requests []BatchRequest) (string, error)
	GetBatchStatus(batchID string) (*BatchStatus, error)
	GetBatchResults(batchID string, prompt prompts.PromptTemplate) ([]BatchResult, error)
}

func NewClient() (Client, error) {
	return theLks.NewClient()
}

func NewBatchClient() (BatchClient, error) {
	return theLks.NewBatchClient()
}

func (lks *LinkedService) NewClient() (Client, error) {
	const semLogContext = "anthropic-lks-registry::new-client"

	cliOpts := []option.RequestOption{
		option.WithAPIKey(lks.Cfg.ApiKey),
	}

	if lks.Cfg.RequestTimeout > 0 {
		cliOpts = append(cliOpts, option.WithMaxRetries(lks.Cfg.MaxRetries))
	}

	if lks.Cfg.RequestTimeout > 0 {
		cliOpts = append(cliOpts, option.WithRequestTimeout(lks.Cfg.RequestTimeout))
	}

	if lks.Cfg.Mockup != nil && lks.Cfg.Mockup.Enabled {

		httpCli, err := lks.httpClientLks.NewClient()
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}

		return &mockupClient{cfg: lks.Cfg.Mockup, httpClient: httpCli}, nil
	}
	anthropicCli := anthropic.NewClient(cliOpts...)
	return &clientImpl{verbose: lks.Cfg.Verbose, apiClient: anthropicCli}, nil
}

func (lks *LinkedService) NewBatchClient() (BatchClient, error) {
	const semLogContext = "anthropic-lks-registry::new-batch-client"

	cliOpts := []option.RequestOption{
		option.WithAPIKey(lks.Cfg.ApiKey),
	}

	if lks.Cfg.MaxRetries > 0 {
		cliOpts = append(cliOpts, option.WithMaxRetries(lks.Cfg.MaxRetries))
	}

	if lks.Cfg.RequestTimeout > 0 {
		cliOpts = append(cliOpts, option.WithRequestTimeout(lks.Cfg.RequestTimeout))
	}

	if lks.Cfg.Mockup != nil && lks.Cfg.Mockup.Enabled {
		httpCli, err := lks.httpClientLks.NewClient()
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}
		return &mockupBatchClient{cfg: lks.Cfg.Mockup, httpClient: httpCli}, nil
	}

	anthropicCli := anthropic.NewClient(cliOpts...)
	return &batchClientImpl{verbose: lks.Cfg.Verbose, apiClient: anthropicCli}, nil
}
