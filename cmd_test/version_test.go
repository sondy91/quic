package cmd_test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionCmd(t *testing.T) {
	cmd := exec.Command("../quic", "version")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)

	// Check the output
	expected := "v0.0.1"
	assert.Equal(t, expected, strings.TrimSpace(string(output)))
}
