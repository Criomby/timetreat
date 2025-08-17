/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
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
		curTask.Stop = ts.Format(time.RFC3339)
		// project name
		if curTask.Project == "" && stopProject == "" {
			// TODO ask for project name
			fmt.Println("WARNING: empty project name")
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

		start, err := time.Parse(time.RFC3339, curTask.Start)
		checkErr(err)
		fmt.Printf("stopped %s at %s\ntook %s\n", curTask.Project, ts.Format(time.TimeOnly), time.Since(start).Truncate(time.Second).String())
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
	// stopCmd.Flags().BoolVarP(&noAsk, "no-ask", "n", false, "do not task for missing info")
	stopCmd.Flags().StringVarP(&stopProject, "project", "p", "", "project name")
	stopCmd.Flags().StringVarP(&stopDescription, "description", "d", "", "activity description")
}
