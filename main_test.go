package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMainCLI(t *testing.T) {
	// Build the CLI binary
	binPath := filepath.Join(t.TempDir(), "test-file-mapper")
	cmdBuild := exec.Command("go", "build", "-o", binPath, "main.go")
	if out, err := cmdBuild.CombinedOutput(); err != nil {
		t.Fatalf("Build failed: %v\n%s", err, string(out))
	}

	// Now run the binary with --help
	cmdRun := exec.Command(binPath, "--help")
	var stdout, stderr bytes.Buffer
	cmdRun.Stdout = &stdout
	cmdRun.Stderr = &stderr

	if err := cmdRun.Run(); err != nil {
		t.Fatalf("Command failed: %v\nstderr: %s", err, stderr.String())
	}

	outStr := stdout.String()
	// The help output contains "USAGE:" in uppercase.
	if !strings.Contains(outStr, "USAGE:") {
		t.Errorf("Expected usage text in --help output, got:\n%s", outStr)
	}
}

func TestMainCLI_FlatList(t *testing.T) {
	// Build CLI
	binPath := filepath.Join(t.TempDir(), "test-file-mapper")
	cmdBuild := exec.Command("go", "build", "-o", binPath, "main.go")
	if out, err := cmdBuild.CombinedOutput(); err != nil {
		t.Fatalf("Build failed: %v\n%s", err, string(out))
	}

	// Create a test directory structure
	tmpDir := t.TempDir()
	sampleFile := filepath.Join(tmpDir, "hello.txt")
	if err := os.WriteFile(sampleFile, []byte("hello content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Run the tool with --flat
	cmdRun := exec.Command(binPath, "--path", tmpDir, "--flat")
	var stdout, stderr bytes.Buffer
	cmdRun.Stdout = &stdout
	cmdRun.Stderr = &stderr

	if err := cmdRun.Run(); err != nil {
		t.Fatalf("Command failed: %v\nstderr: %s", err, stderr.String())
	}

	outStr := stdout.String()
	// Just check that "hello.txt" is somewhere in the output
	if !strings.Contains(outStr, "hello.txt") {
		t.Errorf("Expected 'hello.txt' in the flat listing output, got:\n%s", outStr)
	}
}
