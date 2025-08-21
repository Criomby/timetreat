package cmd

import (
	"cmp"
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/spf13/cobra"
)

var lastCmd = &cobra.Command{
	Use:   "last",
	Short: "Get last used projects",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		EnsureLogFile(false)

		var offset int64 = 0
		var projects []string
		var lastUsedAt []string
		for {
			entry, newOffset, err := GetLogEntryFromEnd(offset)
			CheckErr(err)
			CheckTaskIsNotZero(&entry)
			if !slices.Contains(projects, entry.Project) {
				projects = append(projects, entry.Project)
				lastUsed, err := time.Parse(time.RFC3339, entry.Start)
				CheckErr(err)
				lastUsedAt = append(lastUsedAt, lastUsed.Format(time.DateTime))
			}
			if newOffset == 0 && err == io.EOF {
				break
			}
			offset = newOffset
		}

		for i := 0; i < len(projects); i++ {
			projects[i] = cmp.Or(projects[i], "none")
			fmt.Printf("[%d] %s - %s\n", i, projects[i], lastUsedAt[i])
		}
	},
}

func init() {
	rootCmd.AddCommand(lastCmd)
}
