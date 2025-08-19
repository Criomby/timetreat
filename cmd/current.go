package cmd

import (
	"cmp"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	curProject bool
	curStart   bool
	curDelta   bool
)

var curCmd = &cobra.Command{
	Use:     "current",
	Aliases: []string{"cur"},
	Short:   "Show current task information",
	Run: func(cmd *cobra.Command, args []string) {
		ensureLogFile(false)

		curTask, _, err := getLogEntryFromEnd(0)
		checkErr(err)
		checkTaskIsNotZero(&curTask)
		checkTaskIsRunning(&curTask)

		/* style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4"))
		fmt.Println(style.Render(fmt.Sprintf("%+v", curTask))) */
		start, err := time.Parse(time.RFC3339, curTask.Start)
		checkErr(err)

		var output string
		if !curProject && !curStart && !curDelta {
			if curTask.Project != "" {
				output += fmt.Sprintf("project: %s\n", curTask.Project)
			}
			output += fmt.Sprintf("start: %s\nfor: %s\n", start.Format(time.DateTime), time.Since(start).Truncate(time.Second).String())
			if curTask.Description != "" {
				output += fmt.Sprintf("description: %s\n", curTask.Description)
			}
			fmt.Print(output)
		} else {
			if curProject {
				output += cmp.Or(curTask.Project, "none")
			}
			if curStart {
				start, err := time.Parse(time.RFC3339, curTask.Start)
				checkErr(err)
				if curProject {
					output += " - "
				}
				output += start.Format(time.DateTime)
			}
			if curDelta {
				start, err := time.Parse(time.RFC3339, curTask.Start)
				checkErr(err)
				if curProject || curStart {
					output += " - "
				}
				output += time.Since(start).Truncate(time.Second).String()
			}
			fmt.Println(output)
		}
	},
}

func init() {
	rootCmd.AddCommand(curCmd)
	curCmd.Flags().BoolVarP(&curProject, "project", "p", false, "only list current project name")
	curCmd.Flags().BoolVarP(&curStart, "start", "s", false, "only list current start time")
	curCmd.Flags().BoolVarP(&curDelta, "duration", "d", false, "only list current duration")
}
