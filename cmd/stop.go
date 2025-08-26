package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	stopProject     string
	stopDescription string
	stopTime        string
	stopRound       string
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the currently running task",
	Long:  "TODO",
	Run: func(cmd *cobra.Command, args []string) {
		EnsureLogFile(false)
		checkArgsProjectDescription(stopProject, stopDescription)

		curTask, offset, err := GetLogEntryFromEnd(0)
		CheckErr(err)
		CheckTaskIsNotZero(&curTask)
		CheckTaskIsRunning(&curTask)

		// stop time
		ts := time.Now().Local()
		if stopTime != "" {
			parsedTime, err := time.Parse(timeArgLayout, stopTime)
			CheckErr(err)
			ts = time.Date(ts.Year(), ts.Month(), ts.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, ts.Location())
		}
		if stopRound != "" {
			rd, err := time.ParseDuration(stopRound)
			CheckErr(err)
			ts = ts.Round(rd)
		}

		start, err := time.Parse(time.RFC3339, curTask.Start)
		CheckErr(err)
		if ts.Before(start) {
			formattedStringsStyled.PrintfError("stop time before start time")
			os.Exit(1)
		}

		curTask.Stop = ts.Format(time.RFC3339)

		// project name
		if curTask.Project == "" && stopProject == "" {
			formattedStringsStyled.PrintfWarning("empty project name")
		} else if curTask.Project == "" && stopProject != "" {
			curTask.Project = stopProject
		} else {
			formattedStringsStyled.PrintfWarning("project name exists: %s", curTask.Project)
			answer, err := AskForInputInOptions("[i] ignore or [o] override?", []string{"i", "o"})
			CheckErr(err)
			if answer == "o" {
				curTask.Project = stopProject
			}
		}

		// description
		if curTask.Description == "" && stopDescription == "" {
			// TODO ask for description
		} else if stopDescription != "" {
			totalLen := len(curTask.Project) + len(curTask.Description) + len(stopDescription)
			if totalLen > maxLenProDesc {
				fmt.Printf("description too long (max %d chars, is %d chars)\n", maxLenProDesc, totalLen)
				os.Exit(1)
			}
			if curTask.Description != "" {
				curTask.Description += ", " + stopDescription
			} else {
				curTask.Description = stopDescription
			}
		}

		// adjust for first line in log file and newline char in subsequent entries
		if offset != 0 {
			offset++
		}
		err = RemoveLastLogEntry(offset)
		CheckErr(err)

		err = WriteLogEntry(&curTask)
		CheckErr(err)

		fmt.Printf("stopped %s at %s\ntook %s\n", curTask.Project, ts.Format(time.TimeOnly), ts.Sub(start).Truncate(time.Second).String())
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
	// stopCmd.Flags().BoolVarP(&noAsk, "no-ask", "n", false, "do not task for missing info")
	stopCmd.Flags().StringVarP(&stopProject, "project", "p", "", "project name")
	stopCmd.Flags().StringVarP(&stopDescription, "description", "d", "", "activity description")
	stopCmd.Flags().StringVarP(&stopTime, "time", "t", "", "from a specific time (HH:MM)")
	stopCmd.Flags().StringVarP(&stopRound, "round", "r", "", "round time by e.g. 15m, 1m, etc.")
}
