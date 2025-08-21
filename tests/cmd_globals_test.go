package tests

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"timetreat/cmd"

	"github.com/stretchr/testify/assert"
)

// ╭────────────────────────────╮
// │           Entry            │
// ╰────────────────────────────╯

func TestEntryIsEmpty(t *testing.T) {
	assert := assert.New(t)
	emptyEntry := cmd.Entry{}
	assert.True(emptyEntry.IsEmpty(), "emptyEntry.IsEmpty() returns false")
}

func TestLogHasEntryFields(t *testing.T) {
	logFile := "1_extra_field.log"
	setupTestdataLogFile(t, logFile)
	assert := assert.New(t)

	file, err := os.Open(cmd.GlobalConfig.LogFile)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}

	decoder := json.NewDecoder(strings.NewReader(line))
	decoder.DisallowUnknownFields()

	var entry cmd.Entry
	err = decoder.Decode(&entry)
	assert.EqualErrorf(err, `json: unknown field "extraField"`, "no extra field(s) detected")
}

// ╭────────────────────────────╮
// │          log file          │
// ╰────────────────────────────╯

func TestSetupLogFilePath(t *testing.T) {
	logFile := "9_regular.log"
	setupTestdataLogFile(t, logFile)
	assert := assert.New(t)

	isMatch := strings.HasSuffix(cmd.GlobalConfig.LogFile, "/timetreat/testdata/"+logFile)
	assert.Equal(true, isMatch, "The global config log file should point to */timetreat/testdata/*")
	assert.True(cmd.LogFileExists(), "could not find log file")
}

func TestGetAllLogEntriesFromEnd(t *testing.T) {
	logFile := "9_regular.log"
	setupTestdataLogFile(t, logFile)
	assert := assert.New(t)

	var offset int64 = 0
	var nEntries int = 0
	for {
		entry, newOffset, err := cmd.GetLogEntryFromEnd(offset)
		if err != nil && err != io.EOF {
			t.Fatalf("error reading log file entries: %s", err)
		}
		if !assert.False(entry.IsEmpty(), "log file is unexpectedly empty") {
			t.FailNow()
		}
		nEntries++
		if newOffset == 0 && err == io.EOF {
			break
		}
		offset = newOffset
	}
	assert.Equal(9, nEntries, fmt.Sprintf("read fewer lines than expected: %d/%d", 9, nEntries))
	t.Logf("read %d entries from %s", nEntries, logFile)
}

func TestReadEmptyLog(t *testing.T) {
	logFile := "empty.log"
	setupTestdataLogFile(t, logFile)
	assert := assert.New(t)

	var offset int64 = 0
	entry, newOffset, err := cmd.GetLogEntryFromEnd(offset)
	if err != nil && err != io.EOF {
		t.Fatalf("error reading log file entries: %s", err)
	}
	assert.True(entry.IsEmpty(), "log file is not empty")
	if newOffset != 0 || err != io.EOF {
		t.Fatalf("returned wrong offset and err for empty log file: offset: %d, err: %s", newOffset, err)
	}
}

func TestRemoveLastLogEntry(t *testing.T) {
	logFile := "9_regular.log"
	setupTestdataLogFile(t, logFile)

	newLogfile, err := tempCopyLogFile(t)
	if err != nil {
		t.Fatal(err)
	}
	setupTestdataLogFile(t, newLogfile)

	t.Log("TODO")
}
