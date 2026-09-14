package prompt

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/agents/agentutil"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/agents/promptagent"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/ai-cli/cmds"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/lksregistry"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/store/agentexecution"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	model                string
	outFolder            string
	llmProvider          string
	promptDefinitionFile string
	cfgFileName          string
	batch                bool
	batchPollInterval    time.Duration
	promptVars           map[string]string
)

const (
	DefaultDomain = "--no-domain--"
	DefaultSite   = "--no-site--"
	DefaultGroup  = "--no-group--"

	semLogContextCmd = "prompt::"
)

var theCmd = &cobra.Command{
	Use:   "execute",
	Short: "invokes an llm to execute a prefilled prompt",
	Long:  "invokes an llm to execute a prefilled prompt",
	Run: func(cmd *cobra.Command, args []string) {

		const semLogContext = semLogContextCmd + "run"

		if !cmds.Verbose {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
			log.Info().Bool("verbose", cmds.Verbose).Msg(semLogContext + " log level set to Info")
		} else {
			log.Info().Bool("verbose", cmds.Verbose).Msg(semLogContext + " log level set to Trace")
		}

		if err := validateArgs(); err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return
		}

		log.Info().Str("prompt", promptDefinitionFile).Bool("batch", batch).Msg(semLogContext)

		err := doWork()
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return
		}
	},
}

func doWork() error {
	const semLogContext = semLogContextCmd + "do-work"

	lksCfg, err := lksregistry.ReadConfig(cfgFileName)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	err = lksregistry.InitRegistry(lksCfg)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	agent := promptagent.NewAgentFactory(DefaultDomain, DefaultSite)
	_, batchId, err := agent.Execute(context.Background(), []agentexecution.AgentExecution{{
		Domain: DefaultDomain,
		Site:   DefaultSite,
		Bid:    promptagent.Name,
		Et:     agentexecution.EntityType,
		Status: agentexecution.StatusWorking,
		Weight: 0,
		BidRef: agentexecution.BidEtPair{},
		Params: map[string]any{
			promptagent.ParamModel:                model,
			promptagent.ParamOutFolder:            outFolder,
			promptagent.ParamLlmProvider:          llmProvider,
			promptagent.ParamPromptDefinitionFile: promptDefinitionFile,
			promptagent.ParamBatch:                batch,
			promptagent.ParamBatchPollInterval:    batchPollInterval,
			promptagent.ParamPromptVars:           promptVars,
		},
		Group: DefaultGroup,
	}})

	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	if batchId != "" {
		log.Info().Str("batch-id", batchId).Msg(semLogContext)
	}

	return nil
}

func init() {
	cmds.PromptCmd.AddCommand(theCmd)
	theCmd.Flags().StringVarP(&model, "model", "M", "", "the prompt definition file to use")
	theCmd.Flags().StringVarP(&outFolder, "out-folder", "O", "", "the output folder where to save the analyzed data")
	theCmd.Flags().StringVarP(&promptDefinitionFile, "def", "D", "", "the prompt definition file to use")
	theCmd.Flags().VarP(newEnumValue(&llmProvider, "", "anthropic", "ollama", "vllm"), "llm", "L", "provider to use (one of: anthropic, ollama, vllm)")
	theCmd.Flags().StringVarP(&cfgFileName, "cfg-file", "C", "", "config file of linked services")
	theCmd.Flags().BoolVarP(&batch, "batch", "B", false, "use batch APIs if possible")
	theCmd.Flags().DurationVar(&batchPollInterval, "poll-interval", 30*time.Second, "(used only with batch enabled and provider supporting it, it specifies how often to poll for batch completion; set to 0 to fire-and-forget (prints batch ID and exits)")
	// StringToString map flag. Behavioral notes:
	//   - Repeats merge: --var a=1 --var b=2 -> {a:1, b:2} (first use replaces the
	//     default, later uses add).
	//   - Comma is the separator (CSV-parsed): a value containing a comma must be
	//     CSV-quoted, e.g. -V 'k="a,b"'. The split on '=' is on the first '=', so
	//     values may contain '='.
	//   - Default nil means promptVars stays a nil map when the flag is never passed
	//     (reading is safe; only writing to it would panic). Pass map[string]string{}
	//     as the default instead if a non-nil map is preferred.
	theCmd.Flags().StringToStringVarP(&promptVars, "var", "V", nil, "prompt template variables as key=value (repeatable, or comma-separated: -V a=1,b=2)")

	// Presence is enforced by cobra before Run executes; value validity (paths
	// exist, are files/folders) is still checked in validateArgs.
	for _, name := range []string{"def", "out-folder", "llm", "cfg-file", "model"} {
		_ = theCmd.MarkFlagRequired(name)
	}
}

// enumValue is a pflag.Value that accepts only one of a fixed set of strings and
// writes the chosen value into an existing *string, so callers keep using the
// plain string. An invalid value is rejected by cobra at parse time.
type enumValue struct {
	allowed []string
	value   *string
}

func newEnumValue(target *string, def string, allowed ...string) *enumValue {
	*target = def
	return &enumValue{allowed: allowed, value: target}
}

func (e *enumValue) String() string { return *e.value }

func (e *enumValue) Set(v string) error {
	if slices.Contains(e.allowed, v) {
		*e.value = v
		return nil
	}
	return fmt.Errorf("must be one of %s", strings.Join(e.allowed, ", "))
}

func (e *enumValue) Type() string { return "string" }

func validateArgs() error {
	const semLogContext = semLogContextCmd + "validate-args"

	var err error

	cfgFileName, err = agentutil.ResolveFilename("lks-file", cfgFileName, true)
	if err != nil {
		return err
	}

	return nil
}
