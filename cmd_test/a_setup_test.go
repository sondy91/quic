package cmd_test

import (
	"log"
	"os"
	"os/exec"
	"testing"
)

const binary = "quic"

func TestMain(m *testing.M) {
	// Build the binary before running tests
	buildCmd := exec.Command("go", "build", "-o", binary, "../")
	if err := buildCmd.Run(); err != nil {
		log.Fatalf("Failed to build binary: %v", err)
	}

	// Run the tests
	code := m.Run()

	// Clean up the binary after tests
	os.Remove(binary)
	os.Exit(code)
}
