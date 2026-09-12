package lksregistry

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/anthropiclks"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/ollamalks"
	"github.com/rs/zerolog/log"
)

type ServiceRegistry struct {
}

var registry ServiceRegistry

func InitRegistry(cfg *Config) error {
	var err error

	registry = ServiceRegistry{}
	log.Info().Msg("initialize services registry")

	_, err = anthropiclks.Initialize(cfg.AnthropicAPI)
	if err != nil {
		return err
	}

	_, err = ollamalks.Initialize(cfg.OllamaAPI)
	if err != nil {
		return err
	}

	return nil
}
