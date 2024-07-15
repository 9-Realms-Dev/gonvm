package integration

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

// Testing environment would have the CLI installed already. There will be a test.dockerfile with the an environment setup

func TestInstallLts(t *testing.T) {
	version := "lts"
	cmd := exec.Command("gonvm", "install", "-y", version)
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	assert.Contains(t, string(output), "as the latest lts version")
}
