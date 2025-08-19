package cmd

import (
	"bufio"
	"cmp"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type exportFormatsRegistry struct {
	Csv string
}

var supportedExportFormats = &exportFormatsRegistry{
	Csv: "csv",
}

var (
	exportFormat   string
	exportFilePath string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export log in various formats to a separate file",
	Long: `TODO

Supported formats:
    - csv
	`,
	Run: func(cmd *cobra.Command, args []string) {
		ensureLogFile(false)

		exportFormat = strings.ToLower(exportFormat)
		isValidFormat := false
		formatValue := reflect.ValueOf(*supportedExportFormats)
		for i := 0; i < formatValue.NumField(); i++ {
			if exportFormat == formatValue.Field(i).Interface() {
				isValidFormat = true
			}
		}
		if !isValidFormat {
			fmt.Printf("unsupported format specified: %s\n", exportFormat)
			os.Exit(1)
		}

		exportFilePath = cmp.Or(exportFilePath, filepath.Dir(globalConfig.logFile))

		filename := fmt.Sprintf("timetreat_export_%s.%s", strings.ReplaceAll(time.Now().Format(time.DateTime), " ", "_"), exportFormat)
		fullPath := filepath.Join(exportFilePath, filename)
		if askForConfirmation(fmt.Sprintf("export to '%s'?", fullPath)) {
			file, err := os.Create(fullPath)
			checkErr(err)
			defer file.Close()
			if err := exportLogFile(file, exportFormat); err != nil {
				fmt.Println("error writing export file:", err)
			} else {
				fmt.Println("DONE")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVar(&exportFormat, "format", supportedExportFormats.Csv, "export file format")
	exportCmd.Flags().StringVarP(&exportFilePath, "dir", "d", "", "export file to this dir (default same as log file)")
}

func exportLogFile(file *os.File, format string) error {
	if format == supportedExportFormats.Csv {
		var fieldNames []string
		entryType := reflect.TypeOf(entry{})
		for i := 0; i < entryType.NumField(); i++ {
			fieldNames = append(fieldNames, entryType.Field(i).Name)
		}
		writer := csv.NewWriter(file)
		writer.Comma = ';'
		if err := writer.Write(fieldNames); err != nil {
			return err
		}

		file, err := os.Open(globalConfig.logFile)
		if err != nil {
			return err
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			return err
		}
		fileSize := stat.Size()
		if fileSize == 0 {
			fmt.Println("log file is empty")
			return nil
		}

		scanner := bufio.NewScanner(file)
		var processError error
		var entryBuffer entry
		for scanner.Scan() {
			line := scanner.Text()
			if err = json.Unmarshal([]byte(line), &entryBuffer); err != nil {
				processError = err
				break
			}
			if err := writer.Write([]string{
				entryBuffer.Start,
				entryBuffer.Stop,
				entryBuffer.Project,
				entryBuffer.Description,
			}); err != nil {
				processError = err
				break
			}
		}

		writer.Flush()
		if err := writer.Error(); err != nil {
			fmt.Println("error flushing csv writer:", err)
			return err
		}
		if processError != nil {
			return processError
		}
	} else {
		return fmt.Errorf("export format function not found: %s\n", exportFormat)
	}
	return nil
}
