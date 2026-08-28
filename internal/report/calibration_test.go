package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Lokee86/pitlord/internal/calibration"
)

func TestWriteCalibrationTextSummarizesMismatches(t *testing.T) {
	result := calibration.Result{
		Corpus: "fixture",
		Items: []calibration.Item{{
			ID:      "clean",
			Class:   "clean",
			Outcome: calibration.OutcomeFalsePositive,
		}},
		Summary: calibration.Summary{FalsePositive: 1, LabeledPrecision: 0, LabeledRecall: 1},
	}
	var output bytes.Buffer
	if err := WriteCalibrationText(&output, result); err != nil {
		t.Fatal(err)
	}
	text := output.String()
	if !strings.Contains(text, "FALSE-POSITIVE clean") || !strings.Contains(text, "FP=1") {
		t.Fatalf("unexpected output: %s", text)
	}
}
