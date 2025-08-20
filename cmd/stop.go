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
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the currently running task",
	Long:  "TODO",
	Run: func(cmd *cobra.Command, args []string) {
		ensureLogFile(false)
		checkArgsProjectDescription(stopProject, stopDescription)

		curTask, offset, err := getLogEntryFromEnd(0)
		checkErr(err)
		checkTaskIsNotZero(&curTask)
		checkTaskIsRunning(&curTask)

		ts := time.Now().Local()
		if stopTime != "" {
			parsedTime, err := time.Parse(timeArgLayout, stopTime)
			if err != nil {
				fmt.Printf("%s: %s\n", formattedStringsStyled.Error, err)
				os.Exit(1)
			}
			ts = time.Date(ts.Year(), ts.Month(), ts.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, ts.Location())
		}

		start, err := time.Parse(time.RFC3339, curTask.Start)
		if ts.Before(start) {
			fmt.Printf("%s: Stop time before start time\n", formattedStringsStyled.Error)
			os.Exit(1)
		}

		curTask.Stop = ts.Format(time.RFC3339)

		// project name
		if curTask.Project == "" && stopProject == "" {
			fmt.Printf("%s: empty project name\n", formattedStringsStyled.Warning)
		} else if curTask.Project == "" && stopProject != "" {
			curTask.Project = stopProject
		}
		// description
		if curTask.Description == "" && stopDescription == "" {
			// TODO ask for description
		} else if stopDescription != "" {
			// append description passed
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
		removeLastLogEntry(offset)
		writeLogEntry(&curTask)

		checkErr(err)
		fmt.Printf("stopped %s at %s\ntook %s\n", curTask.Project, ts.Format(time.TimeOnly), ts.Sub(start).Truncate(time.Second).String())
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
	// stopCmd.Flags().BoolVarP(&noAsk, "no-ask", "n", false, "do not task for missing info")
	stopCmd.Flags().StringVarP(&stopProject, "project", "p", "", "project name")
	stopCmd.Flags().StringVarP(&stopDescription, "description", "d", "", "activity description")
	stopCmd.Flags().StringVarP(&stopTime, "time", "t", "", "from a specific time (HH:MM)")
}
