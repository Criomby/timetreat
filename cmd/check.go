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
	Long: `
Checks the integrity of the provided log file.
Not only checks for valid JSON-lines entries but
also verifies the format if the start/stop timestamps.

Verifies the plausability of entries by checking that the
stop time comes after the start time.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		EnsureLogFile(false)
		os.Exit(RunCheck())
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func RunCheck() int {
	file, err := os.Open(GlobalConfig.LogFile)
	if err != nil {
		formattedStringsStyled.PrintfError("failed to open log file\n%s", err)
		return 1
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		formattedStringsStyled.PrintfError("failed to get log file stats\n%s", err)
		return 1
	}
	fileSize := stat.Size()
	if fileSize == 0 {
		fmt.Println("log file is empty")
		return 0
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
		err := json.Unmarshal([]byte(line), &entryBuffer)
		if err != nil {
			hasError = true
			error = err
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
			formattedStringsStyled.PrintfWarning("empty stop value in line %d", nLines)
		}
	}

	if err := scanner.Err(); err != nil {
		formattedStringsStyled.PrintfError("reading log file\n%s", err)
	}

	if hasError {
		formattedStringsStyled.PrintfError("%s\nline %d: %s", error, nLines, line)
	} else {
		fmt.Println(formattedStringsStyled.Ok)
	}
	if hasError {
		return 1
	}
	fmt.Printf("read %d entries\n", nLines)
	return 0
}
