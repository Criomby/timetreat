package tests

import (
	"testing"
	"timetreat/cmd"

	"github.com/stretchr/testify/assert"
)

func TestCheckCmd(t *testing.T) {
	assert := assert.New(t)

	logFile := "9_regular.log"
	setupTestdataLogFile(t, logFile)
	assert.Equal(0, cmd.RunCheck(), "check cmd did not exit with exit code 0")

	logFile = "1_invalid.log"
	setupTestdataLogFile(t, logFile)
	assert.Equal(1, cmd.RunCheck(), "check cmd did not exit with exit code 1")

	logFile = "1_invalid_start_time.log"
	setupTestdataLogFile(t, logFile)
	assert.Equal(1, cmd.RunCheck(), "check cmd did not exit with exit code 1")
}
