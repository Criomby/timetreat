package cmd

import (
	"bufio"
	"cmp"
	"encoding/csv"
	"encoding/json"
	"errors"
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
		EnsureLogFile(false)

		exportFormat = strings.ToLower(exportFormat)
		isValidFormat := false
		formatValue := reflect.ValueOf(*supportedExportFormats)
		for i := 0; i < formatValue.NumField(); i++ {
			if exportFormat == formatValue.Field(i).Interface() {
				isValidFormat = true
			}
		}
		if !isValidFormat {
			formattedStringsStyled.PrintfError("unsupported format specified: %s", exportFormat)
			os.Exit(1)
		}

		exportFilePath = cmp.Or(exportFilePath, filepath.Dir(GlobalConfig.LogFile))
		if !filepath.IsAbs(exportFilePath) {
			formattedStringsStyled.PrintfError("export filepath must be absolute")
			os.Exit(1)
		}

		filename := fmt.Sprintf("timetreat_export_%s.%s", strings.ReplaceAll(time.Now().Format(time.DateTime), " ", "_"), exportFormat)
		fullPath := filepath.Join(exportFilePath, filename)
		if AskForConfirmation(fmt.Sprintf("export to '%s'?", fullPath)) {
			file, err := os.Create(fullPath)
			CheckErr(err)
			defer file.Close()
			if err := ExportLogFile(file, exportFormat); err != nil {
				formattedStringsStyled.PrintfError("writing export file\n%s", err)
			} else {
				fmt.Println(formattedStringsStyled.Ok)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVar(&exportFormat, "format", supportedExportFormats.Csv, "export file format")
	exportCmd.Flags().StringVarP(&exportFilePath, "dir", "d", "", "export file to this dir (default same as log file)")
}

func ExportLogFile(file *os.File, format string) error {
	if format == supportedExportFormats.Csv {
		var fieldNames []string
		entryType := reflect.TypeOf(Entry{})
		for i := 0; i < entryType.NumField(); i++ {
			fieldNames = append(fieldNames, entryType.Field(i).Name)
		}
		writer := csv.NewWriter(file)
		writer.Comma = ';'
		if err := writer.Write(fieldNames); err != nil {
			return err
		}

		file, err := os.Open(GlobalConfig.LogFile)
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
			return errors.New("log file is empty")
		}

		scanner := bufio.NewScanner(file)
		var processError error
		var entryBuffer Entry
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
			return fmt.Errorf("error flushing csv writer: %s", err)
		}
		if processError != nil {
			return processError
		}
	} else {
		return fmt.Errorf("export format function not found: %s", exportFormat)
	}
	return nil
}
