package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/ai-cli/cmds"
	_ "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/ai-cli/cmds/prompt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed sha.txt
var sha string

//go:embed VERSION
var version string

// appLogo contains the ASCII splash screen
//
//go:embed app-logo.txt
var appLogo []byte

func main() {
	fmt.Println(string(appLogo))
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Sha: %s\n", sha)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	cmds.Execute()
}
