package ollamalks

import (
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
	"github.com/rs/zerolog/log"
)

const semLogContextBaseLks = "ollama-lks::"

type LinkedService struct {
	Cfg *Config
}

var theLks LinkedService

func Initialize(cfg *Config) (*LinkedService, error) {
	const semLogContext = semLogContextBaseLks + "initialize"
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

	return lks, err
}

type Client interface {
	Close()
	Execute(params ...RequestParam) (*Response, error)
}

func NewClient() (Client, error) {
	return theLks.NewClient()
}

func (lks *LinkedService) NewClient() (Client, error) {
	const semLogContext = semLogContextBaseLks + "new-client"

	u, err := url.Parse(lks.Cfg.Url)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	client := api.NewClient(u, http.DefaultClient)
	return &clientImpl{ollamaClient: client, verbose: lks.Cfg.Verbose}, nil
}
