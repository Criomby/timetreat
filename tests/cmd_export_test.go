package tests

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"timetreat/cmd"

	"github.com/stretchr/testify/assert"
)

func TestExportLogFileCsv(t *testing.T) {
	logFile := "9_regular.log"
	setupTestdataLogFile(t, logFile)
	assert := assert.New(t)

	newLogfile, err := tempCopyLogFile(t)
	if err != nil {
		t.Fatal(err)
	}
	setupTestdataLogFile(t, newLogfile)

	exportFilename := fmt.Sprintf("timetreat_export_%s.%s", strings.ReplaceAll(time.Now().Format(time.DateTime), " ", "_"), "csv")
	exportFileAbsPath := filepath.Join(filepath.Dir(newLogfile), exportFilename)
	exportFile, err := os.Create(exportFileAbsPath)
	if err != nil {
		t.Fatal(err)
	}
	defer exportFile.Close()

	if err := cmd.ExportLogFile(exportFile, "csv"); err != nil {
		t.Fatal(err)
	} else {
		t.Log("export written:", exportFileAbsPath)
	}

	stat, err := exportFile.Stat()
	if err != nil {
		t.Fatal(err)
	}
	fileSize := stat.Size()
	if fileSize == 0 {
		t.Fatal("export file is empty")
	}
	t.Log("export size:", fileSize)

	_, err = exportFile.Seek(0, 0)
	if err != nil {
		t.Fatal("Error seeking file:", err)
	}
	csvReader := csv.NewReader(exportFile)
	csvReader.Comma = ';'
	records, err := csvReader.ReadAll()
	if err != nil {
		t.Fatal("unable to parse csv export file:", err)
	}
	nRecords := len(records) - 1
	assert.Truef(nRecords == 9, "expected 9 lines, export contains %d", nRecords)
	t.Log(nRecords, "records written")
}
