package cmd

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"time"

	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	listNum   int
	listDelta bool
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List activities in log",
	Long:    `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		EnsureLogFile(false)

		var fieldNames []string
		entryType := reflect.TypeOf(Entry{})
		for i := 0; i < entryType.NumField(); i++ {
			fieldNames = append(fieldNames, entryType.Field(i).Name)
		}
		if listDelta {
			fieldNames = append(fieldNames, "Delta")
		}

		termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
		CheckErr(err)
		t := table.New().
			Headers(fieldNames...).
			Width(termWidth).
			Wrap(true)
		var offset int64 = 0
		var nEntries int = 0
		for nEntries < listNum {
			var rowValues []string
			entry, newOffset, readErr := GetLogEntryFromEnd(offset)
			CheckErr(readErr)
			CheckTaskIsNotZero(&entry)
			start, err := time.Parse(time.RFC3339, entry.Start)
			CheckErr(err)
			stopVal := ""
			var stop time.Time
			if entry.Stop != "" {
				stop, err = time.Parse(time.RFC3339, entry.Stop)
				CheckErr(err)
				stopVal = stop.Format(time.DateTime)
			}
			rowValues = append(rowValues, start.Format(time.DateTime), stopVal, entry.Project, entry.Description)

			if listDelta {
				if entry.Stop == "" {
					rowValues = append(rowValues, "")
				} else {
					rowValues = append(rowValues, stop.Sub(start).String())
				}
			}

			t.Row(rowValues...)
			nEntries++
			if newOffset == 0 && readErr == io.EOF {
				break
			}
			offset = newOffset
		}
		fmt.Println(t.Render())
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().IntVarP(&listNum, "num", "n", 10, "max number of entries to show")
	listCmd.Flags().BoolVarP(&listDelta, "delta", "d", false, "calculate duration of entries")
}
