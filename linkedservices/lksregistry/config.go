package lksregistry

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/anthropiclks"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/ollamalks"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AnthropicAPI *anthropiclks.Config `mapstructure:"anthropic,omitempty" json:"anthropic,omitempty" yaml:"anthropic,omitempty"`
	OllamaAPI    *ollamalks.Config    `mapstructure:"ollama,omitempty" json:"ollama,omitempty" yaml:"ollama,omitempty"`
}

func ReadConfig(fn string) (*Config, error) {
	const semLogContext = "util::read-config"

	var lksCfg Config

	cfgContent, err := util.ReadFileAndResolveEnvVars(fn)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	err = yaml.Unmarshal(cfgContent, &lksCfg)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	return &lksCfg, nil
}
