package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Verify the integrity of the log file",
	Long:  "TODO",
	Run: func(cmd *cobra.Command, args []string) {
		EnsureLogFile(false)

		file, err := os.Open(GlobalConfig.LogFile)
		if err != nil {
			fmt.Printf("%s: failed to open log file\n%s\n", formattedStringsStyled.Error, err)
			os.Exit(1)
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			fmt.Printf("%s: failed to get log file stats\n%s\n", formattedStringsStyled.Error, err)
			os.Exit(1)
		}
		fileSize := stat.Size()
		if fileSize == 0 {
			fmt.Println("log file is empty")
			os.Exit(0)
		}

		scanner := bufio.NewScanner(file)
		nLines := 0
		hasError := false
		var error error
		var line string
		var entryBuffer Entry
		for scanner.Scan() {
			line = scanner.Text()
			nLines++
			error = json.Unmarshal([]byte(line), &entryBuffer)
			if error != nil {
				hasError = true
				break
			}
			timeStart, err := time.Parse(time.RFC3339, entryBuffer.Start)
			if err != nil {
				hasError = true
				error = err
				break
			}
			if entryBuffer.Stop != "" {
				timeStop, err := time.Parse(time.RFC3339, entryBuffer.Stop)
				if err != nil {
					hasError = true
					error = err
					break
				}
				if !timeStart.Before(timeStop) {
					hasError = true
					error = errors.New("start time comes before stop time")
					break
				}
			} else {
				fmt.Printf("%s: empty stop value in line %d\n", formattedStringsStyled.Warning, nLines)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("%s: reading log file\n%s\n", formattedStringsStyled.Error, err)
			os.Exit(1)
		}

		if hasError {
			fmt.Printf(": %s\nline %d: %s\n", formattedStringsStyled.Error, nLines, line)
		} else {
			fmt.Println(formattedStringsStyled.Ok)
		}
		if hasError {
			os.Exit(1)
		}
		fmt.Printf("read %d entries\n", nLines)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
