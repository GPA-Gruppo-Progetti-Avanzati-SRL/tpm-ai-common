package prompt

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/ai-cli/cmds"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/lksregistry"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fileutil"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	cobFileName       string
	cobFileNames      []string
	outFolder         string
	promptsFolder     string
	category          string
	cfgFileName       string
	batch             bool
	batchPollInterval time.Duration
)

const semLogContextCmd = "prompt::"

var theCmd = &cobra.Command{
	Use:   "prompt",
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

		log.Info().Strs("cob-files", cobFileNames).Bool("batch", batch).Msg(semLogContext)

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

	return nil
}

func init() {
	cmds.RootCmd.AddCommand(theCmd)
	theCmd.Flags().StringVarP(&cobFileName, "cob-file", "f", "", "comma-separated list of cob files to process")
	theCmd.Flags().StringVarP(&outFolder, "out-folder", "o", "", "the output folder where to save the analyzed data")
	theCmd.Flags().StringVarP(&promptsFolder, "prompts-repo", "p", "", "folder where prompts are located")
	theCmd.Flags().StringVarP(&promptsFolder, "category", "c", "", "category to use")
	theCmd.Flags().StringVarP(&cfgFileName, "lks-file", "l", "", "linked services config file name")
	theCmd.Flags().BoolVarP(&batch, "batch", "b", false, "submit all files as a single batch request")
	theCmd.Flags().DurationVar(&batchPollInterval, "poll-interval", 30*time.Second, "how often to poll for batch completion; set to 0 to fire-and-forget (prints batch ID and exits)")
}

func validateArgs() error {
	const semLogContext = semLogContextCmd + "validate-args"

	var err error

	cfgFileName, err = resolveFilename("lks-file", cfgFileName, true)
	if err != nil {
		return err
	}

	if cobFileName == "" {
		err = errors.New("error: please specify cob file name(s)")
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	for _, part := range strings.Split(cobFileName, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		resolved, _ := fileutil.ResolvePath(part)
		if !fileutil.FileExists(resolved) {
			err = fmt.Errorf("error: cob file doesn't exist: %s", resolved)
			log.Error().Err(err).Str("cob-file", resolved).Msg(semLogContext)
			return err
		}
		cobFileNames = append(cobFileNames, resolved)
	}

	if len(cobFileNames) == 0 {
		err = errors.New("error: no valid cob files specified")
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	outFolder, err = resolveFolder("out-folder", outFolder, true)
	if err != nil {
		return err
	}

	promptsFolder, err = resolveFolder("prompts-repo", promptsFolder, true)
	if err != nil {
		return err
	}

	return nil
}

func resolveFolder(pn, folder string, required bool) (string, error) {
	const semLogContext = semLogContextCmd + "check-folder"
	var err error

	if folder == "" && required {
		err = fmt.Errorf("error: unspecified folder %s param", pn)
		log.Error().Err(err).Msg(semLogContext)
		return folder, err
	}

	if folder != "" {
		folder, _ = fileutil.ResolvePath(folder)
		if !fileutil.FileExists(filepath.Dir(folder)) {
			err = fmt.Errorf("error: folder %s specifies an invalid path", folder)
			log.Error().Err(err).Msg(semLogContext)
			return folder, err
		}

		if fi, err := os.Stat(folder); err != nil || !fi.IsDir() {
			err = fmt.Errorf("error: folder %s is not a valid folder", folder)
			log.Error().Err(err).Msg(semLogContext)
			return folder, err
		}
	}

	return folder, nil
}

func resolveFilename(pn, fn string, required bool) (string, error) {
	const semLogContext = semLogContextCmd + "resolve-filename"
	var err error

	if fn == "" && required {
		err = fmt.Errorf("error: unspecified file %s param", pn)
		log.Error().Err(err).Msg(semLogContext)
		return fn, err
	}

	if fn != "" {
		fn, _ = fileutil.ResolvePath(fn)
		if !fileutil.FileExists(filepath.Dir(fn)) {
			err = fmt.Errorf("error: file %s specifies an invalid path", fn)
			log.Error().Err(err).Msg(semLogContext)
			return fn, err
		}

		if fi, err := os.Stat(fn); err != nil || fi.IsDir() {
			err = fmt.Errorf("error: file %s is not a valid file", fn)
			log.Error().Err(err).Msg(semLogContext)
			return fn, err
		}
	}

	return fn, nil
}
