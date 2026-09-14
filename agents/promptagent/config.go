package promptagent

import (
	"fmt"
	"time"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/agents/agentutil"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/agents/prompt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/rs/zerolog/log"
)

// Config keys used by NewConfig. Callers build the params map with these so the
// cmd layer and this package cannot drift on key names.
const (
	ParamModel                = "model"
	ParamOutFolder            = "out-folder"
	ParamLlmProvider          = "llm-provider"
	ParamPromptDefinitionFile = "prompt-definition-file"
	ParamPromptVars           = "prompt-vars"
)

type Config struct {
	CustomID             string
	Model                string
	OutFolder            string
	LlmProvider          string
	PromptDefinitionFile string

	Batch             bool
	BatchPollInterval time.Duration
	PromptVars        map[string]string

	Temperature float64
	MaxTokens   int64

	promptDefinition prompt.Definition
}

type AgentParams []Config

func (items AgentParams) FirstByCustomId(cid string) int {

	for i, e := range items {
		if e.CustomID == cid {
			return i
		}
	}

	return -1
}

var DefaultConfig = Config{
	Model:             anthropic.ModelClaudeSonnet4_6,
	Temperature:       0.2,
	MaxTokens:         20000,
	BatchPollInterval: 30 * time.Second,
}

// NewConfig builds a Config starting from DefaultConfig and overriding any field
// whose key is present in params with a value of the matching type. Keys that
// are absent or carry a value of the wrong type leave the default in place.
func NewConfig(params map[string]any) *Config {
	cfg := DefaultConfig

	cfg.CustomID = "custom-id-not-assigned"

	// String fields fall back to the DefaultConfig value when the supplied value
	// is absent, the wrong type, or empty.
	cfg.Model = paramString(params, ParamModel, DefaultConfig.Model)
	cfg.OutFolder = paramString(params, ParamOutFolder, DefaultConfig.OutFolder)
	cfg.LlmProvider = paramString(params, ParamLlmProvider, DefaultConfig.LlmProvider)
	cfg.PromptDefinitionFile = paramString(params, ParamPromptDefinitionFile, DefaultConfig.PromptDefinitionFile)

	if v, ok := params[ParamPromptVars].(map[string]string); ok {
		cfg.PromptVars = v
	}

	return &cfg
}

// Validate checks the config is usable and returns an error describing the first
// problem found. The prompt definition file is resolved to an absolute path and
// written back into the config on success.
func (cfg *Config) Validate() error {
	const semLogContext = semLogPackageContext + "validate"

	if cfg.Model == "" {
		err := fmt.Errorf("error: unspecified %s", ParamModel)
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	if cfg.LlmProvider == "" {
		err := fmt.Errorf("error: unspecified %s", ParamLlmProvider)
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	resolvedFolder, err := agentutil.ResolveFolder(ParamOutFolder, cfg.OutFolder, true)
	if err != nil {
		return err
	}
	cfg.OutFolder = resolvedFolder

	resolvedFile, err := agentutil.ResolveFilename(ParamPromptDefinitionFile, cfg.PromptDefinitionFile, true)
	if err != nil {
		return err
	}
	cfg.PromptDefinitionFile = resolvedFile

	log.Info().
		Str("model", cfg.Model).
		Str("llm", cfg.LlmProvider).
		Bool("batch", cfg.Batch).
		Msg(semLogContext + " prompt agent config validated")

	return nil
}

// paramString returns the string value for key coalesced with defaultValue: the
// supplied value wins when present and non-empty, otherwise defaultValue is used.
// A missing key or a value of the wrong type is treated as empty.
func paramString(params map[string]any, key, defaultValue string) string {
	v, _ := params[key].(string)
	return util.StringCoalesce(v, defaultValue)
}
