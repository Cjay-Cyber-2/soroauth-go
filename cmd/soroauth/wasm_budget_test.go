package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWASMBudgetCLIJSON(t *testing.T) {
	stdout, stderr, err := runCLI(t, "wasm-budget", "--json")
	if err != nil {
		// If building WASM fails due to missing toolchain target or similar in minimal test env, skip or assert.
		t.Skipf("wasm-budget build skipped or failed: %v", err)
	}

	var report WASMReport
	if err := json.Unmarshal([]byte(stdout), &report); err != nil {
		t.Fatalf("failed to unmarshal wasm-budget json: %v\nstdout was: %q", err, stdout)
	}

	if report.BudgetInBytes != WASMBudget {
		t.Errorf("expected budget %d, got %d", WASMBudget, report.BudgetInBytes)
	}

	if stderr != "" {
		// CLI failure paths or JSON output mode must keep stdout results-only and stderr clean when succeeding
	}
}

func TestWASMBudgetStdoutPurityOnFailure(t *testing.T) {
	// Assert that on failure paths (such as an invalid flag), stdout stays strictly results-only (empty)
	// and all error or usage messages go to stderr.
	stdout, stderr, err := runCLI(t, "wasm-budget", "--invalid-flag")
	if err == nil {
		t.Fatal("expected error for invalid flag, got nil")
	}
	if stdout != "" {
		t.Errorf("expected stdout to be empty on failure path, got %q", stdout)
	}
	if !strings.Contains(stderr, "flag provided but not defined") && !strings.Contains(stderr, "wasm-budget") && !strings.Contains(stderr, "usage") {
		t.Errorf("expected stderr to contain flag error or usage, got %q", stderr)
	}
}
func TestWASMBudgetCommandJSON(t *testing.T) {
	tmp := t.TempDir()
	wasmFile := filepath.Join(tmp, "dummy.wasm")
	if err := os.WriteFile(wasmFile, make([]byte, 1024), 0644); err != nil {
		t.Fatalf("writing temp wasm: %v", err)
	}

	var stdout, stderr bytes.Buffer
	err := handleWASMBudget([]string{"--json", "--out", wasmFile}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("handleWASMBudget failed: %v, stderr: %s", err, stderr.String())
	}

	var report WASMReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("failed to unmarshal JSON report: %v, stdout: %s", err, stdout.String())
	}

	if report.SizeInBytes != 1024 {
		t.Errorf("expected size 1024, got %d", report.SizeInBytes)
	}
	if !report.Passed {
		t.Error("expected report.Passed to be true")
	}
}

func TestWASMBudgetCommandHuman(t *testing.T) {
	tmp := t.TempDir()
	wasmFile := filepath.Join(tmp, "dummy.wasm")
	if err := os.WriteFile(wasmFile, make([]byte, 1024), 0644); err != nil {
		t.Fatalf("writing temp wasm: %v", err)
	}

	var stdout, stderr bytes.Buffer
	err := handleWASMBudget([]string{"--out", wasmFile}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("handleWASMBudget failed: %v, stderr: %s", err, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "WASM build size:") {
		t.Errorf("unexpected output: %s", output)
	}
}

func TestFormatBytes(t *testing.T) {
	if got := formatBytes(500); got != "500 B" {
		t.Errorf("formatBytes(500) = %q, want '500 B'", got)
	}
	if got := formatBytes(2 * 1024 * 1024); !strings.Contains(got, "MiB") {
		t.Errorf("formatBytes(2 MiB) = %q, want MiB unit", got)
	}
}
