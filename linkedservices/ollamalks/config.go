package ollamalks

/*
 * for general options see https://github.com/anthropics/anthropic-sdk-go/blob/main/option/requestoption.go
 */

//type RequestOptions struct {
//	Prompt          prompts.PromptTemplate `yaml:"prompt,omitempty" mapstructure:"prompt,omitempty" json:"prompt,omitempty"`
//	Model           string                 `yaml:"model,omitempty" mapstructure:"model,omitempty" json:"model,omitempty"`
//	GenerateOptions *GenerateOptions       `yaml:"options,omitempty" mapstructure:"options,omitempty" json:"options,omitempty"`
//}

// type RequestOption func(o *RequestOptions)

//type MockupConfig struct {
//	HtttpClient *restclient.Config `yaml:"http-client,omitempty" mapstructure:"http-client,omitempty" json:"http-client,omitempty"`
//	Enabled     bool               `yaml:"enabled,omitempty" mapstructure:"enabled,omitempty" json:"enabled,omitempty"`
//	HostName    string             `yaml:"host-name,omitempty" mapstructure:"host-name,omitempty" json:"host-name,omitempty"`
//	ServerPort  int                `yaml:"server-port,omitempty" mapstructure:"server-port,omitempty" json:"server-port,omitempty"`
//	HttpScheme  string             `yaml:"http-scheme,omitempty" mapstructure:"http-scheme,omitempty" json:"http-scheme,omitempty"`
//	Endpoint    string             `yaml:"endpoint,omitempty" mapstructure:"endpoint,omitempty" json:"endpoint,omitempty"`
//}

type Config struct {
	Token   string `yaml:"token,omitempty" mapstructure:"token,omitempty" json:"token,omitempty"`
	Url     string `yaml:"url,omitempty" mapstructure:"url,omitempty" json:"url,omitempty"`
	Verbose bool   `yaml:"verbose,omitempty" mapstructure:"verbose,omitempty" json:"verbose,omitempty"`

	//RequestTimeout time.Duration  `yaml:"request-timeout,omitempty" mapstructure:"request-timeout,omitempty" json:"request-timeout,omitempty"`
	//Mockup         *MockupConfig  `yaml:"mockup,omitempty" mapstructure:"mockup,omitempty" json:"mockup,omitempty"`

	//RequestOptions RequestOptions `yaml:"request-options,omitempty" mapstructure:"request-options,omitempty" json:"request-options,omitempty"`
}
