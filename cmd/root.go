/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "timetreat",
	Version: "0.1.0-alpha",
	Short:   "Treat yourself to some tasty time tracking.",
	Long: `  _   _                _                  _
 | | (_)              | |                | |
 | |_ _ _ __ ___   ___| |_ _ __ ___  __ _| |_
 | __| | '_ ' _ \ / _ \ __| '__/ _ \/ _' | __|
 | |_| | | | | | |  __/ |_| | |  __/ (_| | |_
  \__|_|_| |_| |_|\___|\__|_|  \___|\__,_|\__|

Time tracking made easy and convenient.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&globalConfig.logFile, "log-file", "timetreat.log", "log file")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	globalConfig.logFile = getLogFilePath()
}

func getLogFilePath() string {
	if env := os.Getenv("TIMETREAT_LOG"); env != "" {
		if filepath.IsAbs(env) {
			return env
		} else {
			fmt.Printf("log file path must be absolute: %s\n", env)
			os.Exit(1)
		}
	}
	home, err := os.UserHomeDir()
	checkErr(err)
	return home + "/timetreat.log"
}
