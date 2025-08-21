package tests

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"timetreat/cmd"
)

// Takes the filename of a file in testdata/ and
// sets the global var cmd.GlobalConfig.LogFile to the abspath of that file.
func setupTestdataLogFile(t *testing.T, filename string) {
	if filepath.IsAbs(filename) {
		cmd.GlobalConfig.LogFile = filename
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			t.FailNow()
		}
		cmd.GlobalConfig.LogFile = filepath.Join(cwd, "..", "testdata", filename)
	}
	t.Logf("log file: %s\n", cmd.GlobalConfig.LogFile)
}

// Copies the file from cmd.GlobalConfig.LogFile to a temp dir
// and returns the new path.
func tempCopyLogFile(t *testing.T) (string, error) {
	tempdir := t.TempDir()
	filename := filepath.Base(cmd.GlobalConfig.LogFile)
	newPath := filepath.Join(tempdir, filename)

	source, err := os.Open(cmd.GlobalConfig.LogFile)
	if err != nil {
		return "", err
	}
	defer source.Close()

	destination, err := os.Create(newPath)
	if err != nil {
		return "", err
	}
	defer destination.Close()

	copiedBytes, err := io.Copy(destination, source)
	t.Logf("copied %d bytes\n", copiedBytes)
	if err != nil {
		return "", err
	}

	return newPath, nil
}
