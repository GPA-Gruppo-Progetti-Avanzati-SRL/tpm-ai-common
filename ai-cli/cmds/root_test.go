package cmds_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/ai-cli/cmds"
	_ "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/ai-cli/cmds/prompt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

const (
	TPMGoNgSchematicsVersion = "v0.0.1-SNAPSHOT"
)

type cmdTestCase struct {
	cmdArgs []string
	enabled bool
}

// cmdTestGroup bundles the test cases for one command. A case runs only when
// BOTH its group is enabled AND the case itself is enabled — so you can flip a
// whole command's tests on/off with the group toggle, then pick individual
// cases inside it. Groups also become t.Run subtests, so `go test -run
// TestCmd/j-conv` (or .../j-conv/compile) selects a group or a single case.
type cmdTestGroup struct {
	name    string
	enabled bool
	cases   []cmdTestCase
}

const (
	XsxdatJConvFolder   = "${COB_GAME_COBOLS_PATH}/condizioni/xsxdat-cob-rev5"
	C6sp01axJConvFolder = "${COB_GAME_COBOLS_PATH}/condizioni/c6sp01ax-cob-rev6"

	JConvFolder = C6sp01axJConvFolder
)

var cmdTestGroups = []cmdTestGroup{
	{
		name:    "help",
		enabled: false,
		cases: []cmdTestCase{
			{
				cmdArgs: []string{"prompt", "execute", "--help"},
				enabled: true,
			},
		},
	},
	{
		name:    "prompt",
		enabled: true,
		cases: []cmdTestCase{
			{
				cmdArgs: []string{"prompt", "execute",
					"--def", "./prompt/testdata/node-summary.yml",
					"--model", "claude-sonnet-4-6",
					"--llm", "anthropic",
					"--var", "COBOL_SOURCE=./node-summary-source-code.cob",
					"--cfg-file", "config.yml",
					"--out-folder", "/tmp",
					"--verbose",
					"--mode", "batch",
					// "--resume-batch-id", "msgbatch_01DvGhFKzbUWv5sRuNNwkCs6",
				},
				enabled: false,
			},
			{
				cmdArgs: []string{"prompt", "execute",
					"--def", "./prompt/testdata/node-summary-json.yml",
					"--model", "claude-sonnet-4-6",
					"--llm", "anthropic",
					"--var", "COBOL_SOURCE=./node-summary-source-code.cob",
					"--cfg-file", "config.yml",
					"--out-folder", "/tmp",
					"--verbose",
					"--mode", "batch",
					// "--resume-batch-id", "msgbatch_01T2UmzJXqmskSR1jmmxUDn6",
				},
				enabled: false,
			},
			{
				cmdArgs: []string{"prompt", "execute",
					"--def", "./prompt/testdata/program-info.yml",
					"--model", "claude-sonnet-4-6",
					"--llm", "anthropic",
					"--var", "COBOL_SOURCE=./program-info-source-code.cob",
					"--cfg-file", "config.yml",
					"--out-folder", "/tmp",
					"--verbose",
					"--mode", "batch",
					// "--resume-batch-id", "msgbatch_01T2UmzJXqmskSR1jmmxUDn6",
				},
				enabled: true,
			},
		},
	},
}

// caseLabel derives a short, readable subtest name from a case's arguments:
// the --step value (j-conv), else the --cob-file / --work-folder basename,
// else a positional fallback.
func caseLabel(args []string, idx int) string {
	for i, a := range args {
		if a == "--step" && i+1 < len(args) {
			return args[i+1]
		}
	}
	for _, key := range []string{"--cob-file", "--work-folder"} {
		for i, a := range args {
			if a == key && i+1 < len(args) {
				return filepath.Base(args[i+1])
			}
		}
	}
	return fmt.Sprintf("case-%d", idx)
}

// resetHelpFlag clears cobra's sticky --help bool across the whole command tree.
// cobra adds it lazily to the executed command, so reset every command's copy.
func resetHelpFlag(c *cobra.Command) {
	if f := c.Flags().Lookup("help"); f != nil {
		_ = f.Value.Set("false")
		f.Changed = false
	}
	for _, sub := range c.Commands() {
		resetHelpFlag(sub)
	}
}

func TestCmd(t *testing.T) {
	const semLogContext = "test-commands"
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	for _, group := range cmdTestGroups {
		if !group.enabled {
			continue
		}
		t.Run(group.name, func(t *testing.T) {
			for i, cmdTest := range group.cases {
				if !cmdTest.enabled {
					continue
				}
				t.Run(caseLabel(cmdTest.cmdArgs, i), func(t *testing.T) {
					var tArgs []string
					for _, arg := range cmdTest.cmdArgs {
						tArgs = append(tArgs, arg)
					}

					// The whole suite runs in one process on the shared global
					// cmds.RootCmd, and cobra/pflag do NOT reset flag state between
					// Execute() calls. In particular, cobra's --help is a sticky bool:
					// once an earlier case passes --help it stays true, and a later
					// Execute() short-circuits to help output (printing usage and
					// skipping the command's Run) — so a breakpoint in a later case is
					// never reached and require.NoError still passes. Clear it first.
					resetHelpFlag(cmds.RootCmd)

					cmds.RootCmd.SetArgs(tArgs)
					cmds.Version = TPMGoNgSchematicsVersion

					err := cmds.RootCmd.Execute()
					require.NoError(t, err)
				})
			}
		})
	}
}
