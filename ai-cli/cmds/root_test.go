package cmds_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/ai-cli/cmds"
	_ "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/ai-cli/cmds/prompt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
		name:    "prompt",
		enabled: true,
		cases: []cmdTestCase{
			{
				cmdArgs: []string{"prompt",
					"--format", "vmember",
					"--file", "/Users/marioa.imperato/projects/tpm/game/cob-game-cobol-samples/cob-sources/rpol/sispar/FIRACC.CDBP.CO2P.txt",
					"--out", "/Users/marioa.imperato/projects/tpm/game/cob-game-cobol-samples/cob-sources/rpol/sispar/pgms",
					"--verbose",
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

func TestCmd(t *testing.T) {
	const semLogContext = "test-commands"
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cobol_paths := os.Getenv("COB_GAME_COBOLS_PATH")
	if cobol_paths == "" {
		cobol_paths = "/Users/marioa.imperato/projects/tpm/game/cob-game-cobol-samples-exports"
	}
	require.NotEmpty(t, cobol_paths)

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
						arg = strings.ReplaceAll(arg, "${COB_GAME_COBOLS_PATH}", cobol_paths)
						tArgs = append(tArgs, arg)
					}

					cmds.RootCmd.SetArgs(tArgs)
					cmds.Version = TPMGoNgSchematicsVersion

					err := cmds.RootCmd.Execute()
					require.NoError(t, err)
				})
			}
		})
	}
}
