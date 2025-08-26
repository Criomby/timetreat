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

// ╭────────────────────────────╮
// │       global config        │
// ╰────────────────────────────╯

// max length of project name and description combined
const maxLenProDesc = 150

// Layout of time option value for start/stop commands.
const timeArgLayout string = "15:04"

type appState struct {
	LogFile string
}

var GlobalConfig appState

type formattedStrings struct {
	Ok      string
	Warning string
	Error   string
}

func (f *formattedStrings) PrintfOk(format string, a ...any) {
	fmt.Printf("%s: %s\n", f.Ok, fmt.Sprintf(format, a...))
}

func (f *formattedStrings) PrintfWarning(format string, a ...any) {
	fmt.Printf("%s: %s\n", f.Warning, fmt.Sprintf(format, a...))
}

func (f *formattedStrings) PrintfError(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", f.Error, fmt.Sprintf(format, a...))
}

var formattedStringsStyled *formattedStrings = &formattedStrings{
	Ok:      lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("Ok"),
	Warning: lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("Warning"),
	Error:   lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("Error"),
}

// ╭────────────────────────────╮
// │           Entry            │
// ╰────────────────────────────╯

// an activity file entry json struct
type Entry struct {
	Start       string `json:"start"`
	Stop        string `json:"stop"`
	Project     string `json:"project"`
	Description string `json:"description"`
}

func (e *Entry) IsEmpty() bool {
	if e.Start == "" && e.Stop == "" && e.Project == "" && e.Description == "" {
		return true
	}
	return false
}

// Reads a log entry from the end of the file starting from a given offset.
// 0 offset reads the last line.
func GetLogEntryFromEnd(offset int64) (Entry, int64, error) {
	file, err := os.Open(GlobalConfig.LogFile)
	if err != nil {
		return Entry{}, 0, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return Entry{}, 0, err
	}
	fileSize := stat.Size()
	if fileSize == 0 {
		return Entry{}, 0, io.EOF
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
		return Entry{}, 0, err
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
		// return Entry{}, 0, nil
	}

	if offset == fileSize && lastLine == string(buffer) {
		newOffset = 0
	}

	var entryBuffer Entry
	err = json.Unmarshal([]byte(lastLine), &entryBuffer)
	if err != nil {
		return Entry{}, 0, err
	}

	if newOffset == 0 {
		return entryBuffer, newOffset, io.EOF
	}
	return entryBuffer, newOffset, nil
}

func RemoveLastLogEntry(offset int64) error {
	file, err := os.OpenFile(GlobalConfig.LogFile, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := file.Truncate(offset); err != nil {
		return err
	}
	return nil
}

func WriteLogEntry(task *Entry) error {
	taskBytes, _ := json.Marshal(task)
	file, err := os.OpenFile(GlobalConfig.LogFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.WriteString(string(taskBytes) + "\n"); err != nil {
		return err
	}
	return nil
}

// ╭────────────────────────────╮
// │           utils            │
// ╰────────────────────────────╯

func LogFileExists() bool {
	info, err := os.Stat(GlobalConfig.LogFile)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return !info.IsDir()
}

// If create is true, creates the missing log file,
// else it exists with exit code 1 when missing.
func EnsureLogFile(create bool) {
	if !LogFileExists() {
		fmt.Println("log file does not exist")
		if !create {
			fmt.Println("run the start command to create it")
			os.Exit(1)
		}
		if AskForConfirmation("create the log file?") {
			file, err := os.Create(GlobalConfig.LogFile)
			CheckErr(err)
			file.Close()
			fmt.Printf("created log file %s\n", GlobalConfig.LogFile)
		}
		os.Exit(0)
	}
}

// ╭────────────────────────────╮
// │     convenience funcs      │
// ╰────────────────────────────╯

// Convenience function to print errors to stderr with colored output
// indicating an error if err != nil && err != io.EOF.
func CheckErr(err error) {
	if err != nil && err != io.EOF {
		formattedStringsStyled.PrintfError("%s", err)
		os.Exit(1)
	}
}

func CheckTaskIsRunning(task *Entry) {
	if task.Stop != "" {
		fmt.Println("no task running")
		os.Exit(0)
	}
}

// check if file is empty
func CheckTaskIsNotZero(task *Entry) {
	if task.IsEmpty() == true {
		fmt.Println("log file is empty")
		os.Exit(1)
	}
}

func AskForConfirmation(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/n) ", prompt)
	input, err := reader.ReadString('\n')
	CheckErr(err)
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

func AskForInputInOptions(prompt string, options []string) (string, error) {
	inputError := fmt.Errorf("%s %v\n", "input must be one of", options)
	fmt.Printf("%s ", prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	CheckErr(err)
	input = strings.ToLower(strings.TrimSpace(input))
	if len(input) != 1 {
		return "", inputError
	}

	for _, option := range options {
		if strings.Contains(input, option) {
			return string(option), nil
		}
	}
	return "", inputError
}
