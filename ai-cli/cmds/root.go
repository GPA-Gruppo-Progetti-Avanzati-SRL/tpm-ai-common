/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

package cmds

import (
	"os"

	"github.com/spf13/cobra"
)

var Version = ""
var Verbose bool

// rootCmd represents the base command when called without any subcommands

var RootCmd = &cobra.Command{
	Use:   "ai-cli",
	Short: "command line tool to handle tasks against LLMs",
	Long:  `Preliminary commands available`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var PromptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "command for doing some prompt related actions",
	Long:  `command for doing some prompt related actions`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tpm-iso20022-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.AddCommand(PromptCmd)
	RootCmd.PersistentFlags().BoolVar(&Verbose, "verbose" /*"v",*/, false, "verbose output")
}
