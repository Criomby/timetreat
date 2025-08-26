package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

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
`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	GlobalConfig.LogFile = getLogFilePath()
}

func getLogFilePath() string {
	if env := os.Getenv("TIMETREAT_LOG"); env != "" {
		if filepath.IsAbs(env) {
			return env
		} else {
			formattedStringsStyled.PrintfError("log file path must be absolute: %s", env)
			os.Exit(1)
		}
	}
	home, err := os.UserHomeDir()
	CheckErr(err)
	return home + "/timetreat.log"
}
