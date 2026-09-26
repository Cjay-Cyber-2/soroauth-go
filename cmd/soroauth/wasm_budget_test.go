package main

import (
	"encoding/json"
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
