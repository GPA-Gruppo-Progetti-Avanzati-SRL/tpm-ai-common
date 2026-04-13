package anthropiclks

import (
	"time"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
)

/*
 * for general options see https://github.com/anthropics/anthropic-sdk-go/blob/main/option/requestoption.go
 */

//type ClientOptions struct {
//	MaxTokens   int64                  `yaml:"max-tokens,omitempty" mapstructure:"max-tokens,omitempty" json:"max-tokens,omitempty"`
//	Temperature float64                `yaml:"temperature,omitempty" mapstructure:"temperature,omitempty" json:"temperature,omitempty"`
//	Prompt      prompts.PromptTemplate `yaml:"prompt,omitempty" mapstructure:"prompt,omitempty" json:"prompt,omitempty"`
//	Model       anthropic.Model        `yaml:"model,omitempty" mapstructure:"model,omitempty" json:"model,omitempty"`
//}

// type ClientOption func(o *ClientOptions)

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
	// ClientOptions  ClientOptions `yaml:"client-options,omitempty" mapstructure:"client-options,omitempty" json:"client-options,omitempty"`
}
