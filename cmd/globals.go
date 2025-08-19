package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// max length of project name and description combined
const maxLenProDesc = 150

// an activity file entry json struct
type entry struct {
	Start       string `json:"start"`
	Stop        string `json:"stop"`
	Project     string `json:"project"`
	Description string `json:"description"`
	// Tags        []string `json:"tags"`
}

func (e *entry) isZero() bool {
	if e.Start == "" && e.Stop == "" && e.Project == "" && e.Description == "" {
		return true
	}
	return false
}

// ╭────────────────────────────╮
// │       global config        │
// ╰────────────────────────────╯
type appState struct {
	logFile string
}

var globalConfig appState

// ╭────────────────────────────╮
// │           utils            │
// ╰────────────────────────────╯

type formattedStrings struct {
	Ok      string
	Warning string
	Error   string
}

var formattedStringsStyled *formattedStrings = &formattedStrings{
	Ok:      lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("Ok"),
	Warning: lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("Warning"),
	Error:   lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("Error"),
}

// DEBUG convenience function to panic errors,
// TODO replace for production
func checkErr(err error) {
	if err != nil && err != io.EOF {
		panic(err)
	}
}

func askForConfirmation(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/n) ", prompt)
	input, err := reader.ReadString('\n')
	checkErr(err)
	input = strings.ToLower(strings.TrimSpace(input))
	if input == "y" || input == "yes" {
		return true
	} else if input == "n" || input == "no" {
		return false
	} else {
		fmt.Println("Invalid input. Please enter y/yes or n/no.")
		os.Exit(1)
		return false
	}
}

func logFileExists() bool {
	info, err := os.Stat(globalConfig.logFile)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return !info.IsDir()
}

func createLogFile() {
	file, err := os.Create(globalConfig.logFile)
	checkErr(err)
	file.Close()
}

// If create is true, creates the missing log file,
// else it exists with exit code 1 when missing.
func ensureLogFile(create bool) {
	if !logFileExists() {
		fmt.Println("log file does not exist")
		if !create {
			fmt.Println("run the start command to create it")
			os.Exit(1)
		}
		if askForConfirmation("create the log file?") {
			createLogFile()
			fmt.Printf("created log file %s\n", globalConfig.logFile)
		}
		os.Exit(0)
	}
}

// Reads a log entry from the end of the file starting from a given offset.
// 0 offset reads the last line.
func getLogEntryFromEnd(offset int64) (entry, int64, error) {
	file, err := os.Open(globalConfig.logFile)
	if err != nil {
		return entry{}, 0, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return entry{}, 0, err
	}
	fileSize := stat.Size()
	if fileSize == 0 {
		return entry{}, 0, nil
	}

	if offset == 0 {
		offset = fileSize
	}

	bufferSize := int64(1024)
	if offset < bufferSize {
		bufferSize = offset
	}

	buffer := make([]byte, bufferSize)
	readOffset := offset - bufferSize
	_, err = file.ReadAt(buffer, readOffset)
	if err != nil && err != io.EOF {
		return entry{}, 0, err
	}

	var lastLine string
	var newOffset int64
	for i := bufferSize - 1; i >= 0; i-- {
		if buffer[i] == '\n' {
			if i < bufferSize-1 {
				lastLine = string(buffer[i+1:])
				newOffset = readOffset + int64(i)
				break
			}
		}
	}

	if lastLine == "" {
		// TODO remove?
		lastLine = string(buffer)
		newOffset = 0
		// file empty or corrupt
		// return entry{}, 0, nil
	}

	// Handle the edge case where the file does not end with a newline.
	// If the last line is the entire buffer, the new offset should be 0.
	if offset == fileSize && lastLine == string(buffer) {
		newOffset = 0
	}

	var entryBuffer entry
	err = json.Unmarshal([]byte(lastLine), &entryBuffer)
	if err != nil {
		return entry{}, 0, err
	}

	if newOffset == 0 && offset > 0 {
		return entryBuffer, newOffset, io.EOF
	}
	return entryBuffer, newOffset, nil
}

// TODO remove in favor of getLogEntryFromEnd.
// Get last entry from log file.
// Returns the last line entry and the offset of the last line.
func getLastLogEntry() (entry, int64) {
	file, err := os.Open(globalConfig.logFile)
	checkErr(err)
	defer file.Close()

	stat, err := file.Stat()
	checkErr(err)
	fileSize := stat.Size()
	if fileSize == 0 {
		fmt.Println("log file is empty")
		os.Exit(1)
	}

	bufferSize := int64(512) // at maxLenProDesc == 50
	if fileSize < bufferSize {
		bufferSize = fileSize
	}

	buffer := make([]byte, bufferSize)
	_, err = file.ReadAt(buffer, fileSize-bufferSize)
	checkErr(err)

	var lastLine string
	var offset int64
	for i := bufferSize - 2; i >= 0; i-- {
		if buffer[i] == '\n' {
			offset = fileSize - bufferSize + int64(i) + 1
			lastLine = string(buffer[i+1:])
			break
		}
	}
	if lastLine == "" {
		lastLine = string(buffer)
		offset = 0
	}

	var entryBuffer entry
	err = json.Unmarshal([]byte(lastLine), &entryBuffer)
	if err != nil {
		fmt.Println("Cannot get log entry:", err)
		fmt.Println(lastLine)
		os.Exit(1)
	}
	return entryBuffer, offset
}

func removeLastLogEntry(offset int64) {
	file, err := os.OpenFile(globalConfig.logFile, os.O_RDWR, 0644)
	checkErr(err)
	defer file.Close()
	if err := file.Truncate(offset); err != nil {
		checkErr(err)
	}
}

func writeLogEntry(task *entry) {
	taskBytes, _ := json.Marshal(task)
	file, err := os.OpenFile(globalConfig.logFile, os.O_APPEND|os.O_WRONLY, 0644)
	checkErr(err)
	defer file.Close()
	if _, err := file.WriteString(string(taskBytes) + "\n"); err != nil {
		checkErr(err)
	}
}

func checkTaskIsRunning(task *entry) {
	if task.Stop != "" {
		fmt.Println("no task running")
		os.Exit(0)
	}
}

func isTaskRunning(task *entry) bool {
	if task.Stop == "" {
		return true
	}
	return false
}

// check if file is empty
func checkTaskIsNotZero(task *entry) {
	if task.isZero() == true {
		fmt.Println("log file is empty")
		os.Exit(1)
	}
}
