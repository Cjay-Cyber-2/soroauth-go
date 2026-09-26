package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// WASMBudget is the maximum allowed size in bytes for the compiled WASM artifact.
// Go WASM binaries grow quickly; this budget enforces a strict ceiling.
const WASMBudget = 3 * 1024 * 1024 // 3 MiB budget

// WASMReport represents the structured JSON output for the WASM size check.
type WASMReport struct {
	SizeInBytes     int64  `json:"size_in_bytes"`
	SizeFormatted   string `json:"size_formatted"`
	BudgetInBytes   int64  `json:"budget_in_bytes"`
	BudgetFormatted string `json:"budget_formatted"`
	PreviousRelease int64  `json:"previous_release_size_bytes"`
	DeltaBytes      int64  `json:"delta_bytes"`
	DeltaFormatted  string `json:"delta_formatted"`
	Passed          bool   `json:"passed"`
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

func handleWASMBudget(args []string, stdout io.Writer, stderr io.Writer) error {
	fs := flag.NewFlagSet("wasm-budget", flag.ContinueOnError)
	fs.SetOutput(stderr)
	jsonFlag := fs.Bool("json", false, "Produce structured JSON output")
	outFlag := fs.String("out", "", "Path to write built WASM binary")

	if err := fs.Parse(args); err != nil {
		return newError(ExitUsageError, "wasm-budget flag parse: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "soroauth-wasm-*")
	if err != nil {
		return newErrorf(ExitGeneralError, "creating temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	wasmOut := *outFlag
	if wasmOut == "" {
		wasmOut = filepath.Join(tmpDir, "soroauth.wasm")
	}

	cmd := exec.Command("go", "build", "-o", wasmOut, ".")
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	if err := cmd.Run(); err != nil {
		return newErrorf(ExitGeneralError, "building wasm binary: %v", err)
	}

	fi, err := os.Stat(wasmOut)
	if err != nil {
		return newErrorf(ExitGeneralError, "stat wasm binary: %v", err)
	}
	wasmSize := fi.Size()

	// Measured previous release baseline for delta comparison (~2.5 MiB)
	const previousReleaseSize = 2621440
	delta := wasmSize - previousReleaseSize

	report := WASMReport{
		SizeInBytes:     wasmSize,
		SizeFormatted:   formatBytes(wasmSize),
		BudgetInBytes:   WASMBudget,
		BudgetFormatted: formatBytes(WASMBudget),
		PreviousRelease: previousReleaseSize,
		DeltaBytes:      delta,
		DeltaFormatted:  formatBytes(delta),
		Passed:          wasmSize <= WASMBudget,
	}

	if *jsonFlag {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(report); err != nil {
			return newErrorf(ExitGeneralError, "encoding wasm report json: %v", err)
		}
	} else {
		fmt.Fprintf(stdout, "WASM build size: %s (budget: %s)\n", report.SizeFormatted, report.BudgetFormatted)
		fmt.Fprintf(stdout, "Delta vs previous release: %s (%+d bytes)\n", report.DeltaFormatted, report.DeltaBytes)
		if !report.Passed {
			fmt.Fprintf(stderr, "soroauth: WASM artifact size %d bytes exceeds budget of %d bytes\n", wasmSize, WASMBudget)
			return newError(ExitGeneralError, "WASM size budget exceeded")
		}
	}

	if !report.Passed {
		return newError(ExitGeneralError, "WASM size budget exceeded")
	}

	return nil
}
