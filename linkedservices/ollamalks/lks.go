package ollamalks

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/rs/zerolog/log"
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
	cfg.ClientOptions = DefaultClientOptions(cfg.ClientOptions)
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

func NewClient(opts ...ClientOption) (Client, error) {
	return theLks.NewClient(opts...)
}

func (lks *LinkedService) NewClient(opts ...ClientOption) (Client, error) {
	const semLogContext = "anthropic-lks-registry::new-client"

	options := lks.Cfg.ClientOptions
	for _, o := range opts {
		o(&options)
	}

	httpCli, err := lks.httpClientLks.NewClient()
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	return &mockupClient{cfg: lks.Cfg.Mockup, httpClient: httpCli, options: options}, nil
}
