package ollamalks

import (
	"time"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/prompts"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/rs/zerolog/log"
)

/*
 * for general options see https://github.com/anthropics/anthropic-sdk-go/blob/main/option/requestoption.go
 */

type ClientOptions struct {
	Prompt prompts.PromptTemplate `yaml:"prompt,omitempty" mapstructure:"prompt,omitempty" json:"prompt,omitempty"`
	Model  string                 `yaml:"model,omitempty" mapstructure:"model,omitempty" json:"model,omitempty"`
}

type ClientOption func(o *ClientOptions)

func WithPrompt(p prompts.PromptTemplate) ClientOption {
	return func(o *ClientOptions) {
		o.Prompt = p
	}
}

func WithPromptName(n string) ClientOption {
	return func(o *ClientOptions) {
		pt, err := prompts.GetPrompt(n)
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to get prompt %s", n)
			return
		}
		o.Prompt = pt
	}
}

func WithModel(n string) ClientOption {
	return func(o *ClientOptions) {
		o.Model = n
	}
}

var defCliOptions = ClientOptions{Model: "granite4:latest"}

func DefaultClientOptions(opts ClientOptions) ClientOptions {
	opts.Model = util.StringCoalesce(opts.Model, defCliOptions.Model)
	return opts
}

type MockupConfig struct {
	HtttpClient *restclient.Config `yaml:"http-client,omitempty" mapstructure:"http-client,omitempty" json:"http-client,omitempty"`
	Enabled     bool               `yaml:"enabled,omitempty" mapstructure:"enabled,omitempty" json:"enabled,omitempty"`
	HostName    string             `yaml:"host-name,omitempty" mapstructure:"host-name,omitempty" json:"host-name,omitempty"`
	ServerPort  int                `yaml:"server-port,omitempty" mapstructure:"server-port,omitempty" json:"server-port,omitempty"`
	HttpScheme  string             `yaml:"http-scheme,omitempty" mapstructure:"http-scheme,omitempty" json:"http-scheme,omitempty"`
	Endpoint    string             `yaml:"endpoint,omitempty" mapstructure:"endpoint,omitempty" json:"endpoint,omitempty"`
}

type Config struct {
	ApiKey         string        `yaml:"api-key,omitempty" mapstructure:"api-key,omitempty" json:"api-key,omitempty"`
	MaxRetries     int           `yaml:"max-retries,omitempty" mapstructure:"max-retries,omitempty" json:"max-retries,omitempty"`
	RequestTimeout time.Duration `yaml:"request-timeout,omitempty" mapstructure:"request-timeout,omitempty" json:"request-timeout,omitempty"`
	Mockup         *MockupConfig `yaml:"mockup,omitempty" mapstructure:"mockup,omitempty" json:"mockup,omitempty"`
	Verbose        bool          `yaml:"verbose,omitempty" mapstructure:"verbose,omitempty" json:"verbose,omitempty"`
	ClientOptions  ClientOptions `yaml:"client-options,omitempty" mapstructure:"client-options,omitempty" json:"client-options,omitempty"`
}
