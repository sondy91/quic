package cmd_test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONCmd_Encode(t *testing.T) {
	cmd := exec.Command("./quic", "json", "encode", "test")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)

	// Check the output
	expected := `"test"`
	assert.Equal(t, expected, string(output))
}

func TestJSONCmd_Decode(t *testing.T) {
	// Run the command
	cmd := exec.Command("./quic", "json", "decode", `"test"`)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)

	// Check the output
	expected := "Decoded JSON: test"
	assert.Equal(t, expected, strings.TrimSpace(string(output)))
}

func TestJSONCmd_InvalidAction(t *testing.T) {
	// Run the command
	cmd := exec.Command("./quic", "json", "invalid", "test")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)

	// Check the output
	expected := "Invalid action. Please use `encode` or `decode`"
	assert.Equal(t, expected, strings.TrimSpace(string(output)))
}

func TestJSONCmd_DecodeInvalidJSON(t *testing.T) {
	// Run the command
	cmd := exec.Command("./quic", "json", "decode", "invalid")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)

	// Check the output
	expected := "invalid character 'i' looking for beginning of value"
	assert.Contains(t, string(output), expected)
}
