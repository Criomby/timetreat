package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	startProject     string
	startPrev        bool
	startDescription string
	startTime        string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new tracking entry",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		ensureLogFile(true)
		checkArgsProjectDescription(startProject, startDescription)

		// TODO check if curent task is running and end it gracefully if it is
		prevTask, _, err := getLogEntryFromEnd(0)
		checkErr(err)
		if !prevTask.isZero() && prevTask.Stop == "" {
			fmt.Println("previous task is still running, stop it first")
			os.Exit(1)
		}

		var project string
		projectFlagValue := startProject
		if startPrev == true && projectFlagValue != "" {
			fmt.Println("conflicting options: --prev and --project name")
			os.Exit(1)
		}
		if startPrev == true {
			if prevTask.isZero() {
				fmt.Println("no previous task exists")
				os.Exit(1)
			}
			project = prevTask.Project
		} else {
			project = projectFlagValue
		}

		ts := time.Now().Local()
		if startTime != "" {
			parsedTime, err := time.Parse(timeArgLayout, startTime)
			if err != nil {
				fmt.Printf("%s: %s\n", formattedStringsStyled.Error, err)
				os.Exit(1)
			}
			ts = time.Date(ts.Year(), ts.Month(), ts.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, ts.Location())
		}

		writeLogEntry(&entry{
			Start:       ts.Format(time.RFC3339),
			Project:     project,
			Description: startDescription,
		})

		output := ts.Format(time.TimeOnly)
		if project != "" {
			output += " - " + project
		}
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&startProject, "project", "p", "", "project name")
	startCmd.Flags().BoolVar(&startPrev, "prev", false, "copy the last used project name for this entry")
	startCmd.Flags().StringVarP(&startDescription, "description", "d", "", "activity description")
	startCmd.Flags().StringVarP(&startTime, "time", "t", "", "from a specific time (HH:MM)")
}

// verify the length of project name and description
func checkArgsProjectDescription(project string, description string) {
	if len(project)+len(description) > maxLenProDesc {
		fmt.Printf("project name and description too long (max %d chars, is %d chars)\n", maxLenProDesc, len(project)+len(description))
		os.Exit(1)
	}
}
