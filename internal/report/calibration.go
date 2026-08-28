package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/Lokee86/pitlord/internal/calibration"
)

func WriteCalibrationJSON(writer io.Writer, result calibration.Result) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func WriteCalibrationText(writer io.Writer, result calibration.Result) error {
	for _, item := range result.Items {
		if item.Outcome != calibration.OutcomeFalsePositive &&
			item.Outcome != calibration.OutcomeFalseNegative &&
			!item.SeverityMismatch {
			continue
		}
		if _, err := fmt.Fprintf(writer, "%s %s [%s]", strings.ToUpper(string(item.Outcome)), item.ID, item.Class); err != nil {
			return err
		}
		if item.SeverityMismatch {
			if _, err := fmt.Fprint(writer, " severity-mismatch"); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(writer); err != nil {
			return err
		}
		for _, finding := range item.MatchedFindings {
			if _, err := fmt.Fprintf(writer, "  %s %s\n", finding.Severity, finding.Scope.Path); err != nil {
				return err
			}
		}
	}
	if len(result.UnlabelledFindings) > 0 {
		if _, err := fmt.Fprintf(writer, "Unlabelled findings: %d\n", len(result.UnlabelledFindings)); err != nil {
			return err
		}
		for _, finding := range result.UnlabelledFindings {
			if _, err := fmt.Fprintf(writer, "  %s %s\n", finding.Severity, finding.Scope.Path); err != nil {
				return err
			}
		}
	}
	_, err := fmt.Fprintf(
		writer,
		"Calibration %s: TP=%d TN=%d FP=%d FN=%d severity=%d unlabelled=%d precision=%.3f recall=%.3f\n",
		result.Corpus,
		result.Summary.TruePositive,
		result.Summary.TrueNegative,
		result.Summary.FalsePositive,
		result.Summary.FalseNegative,
		result.Summary.SeverityMismatches,
		result.Summary.UnlabelledFindings,
		result.Summary.LabeledPrecision,
		result.Summary.LabeledRecall,
	)
	return err
}
