//go:build e2e

package e2e

import (
	"encoding/json"
	"os"
	"testing"
)

func TestParityReportRegression(t *testing.T) {
	// Ensure a committed regression fixture exists and is parseable
	fixturePath := "testdata/parity_regression.json"
	raw, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("failed to read parity regression fixture: %v", err)
	}

	var report parityReport
	if err := json.Unmarshal(raw, &report); err != nil {
		t.Fatalf("failed to unmarshal parity regression fixture: %v", err)
	}

	if report.TotalScenarios == 0 {
		t.Fatal("parity regression fixture contains zero scenarios")
	}
}
