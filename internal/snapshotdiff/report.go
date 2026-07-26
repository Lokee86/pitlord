package snapshotdiff

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/Lokee86/pitlord/internal/policy"
	pitlordreport "github.com/Lokee86/pitlord/internal/report"
)

func WriteJSON(writer io.Writer, result Result) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(result)
}

func WriteText(writer io.Writer, result Result) error {
	if _, err := fmt.Fprintf(
		writer,
		"Pitlord snapshot diff: %d introduced, %d resolved, %d persistent finding(s).\n",
		result.Summary.Introduced,
		result.Summary.Resolved,
		result.Summary.Persistent,
	); err != nil {
		return err
	}
	if len(result.Introduced) > 0 {
		if _, err := fmt.Fprintln(writer, "\nIntroduced:"); err != nil {
			return err
		}
		if err := writeFindings(writer, result.Introduced); err != nil {
			return err
		}
	}
	if len(result.Resolved) > 0 {
		if _, err := fmt.Fprintln(writer, "\nResolved:"); err != nil {
			return err
		}
		if err := writeFindings(writer, result.Resolved); err != nil {
			return err
		}
	}
	return nil
}

func WriteIntroducedSARIF(writer io.Writer, result Result, toolVersion string) error {
	diagnostics := Diagnostics(result.Introduced)
	report := policy.Report{
		Schema:       "pitlord.report.v1",
		Snapshot:     result.AfterSnapshot,
		PolicySource: result.PolicySource,
		Diagnostics:  diagnostics,
		Summary: policy.Summary{
			Errors:   countSeverity(diagnostics, "error"),
			Warnings: countSeverity(diagnostics, "warning"),
		},
	}
	return pitlordreport.WriteSARIF(writer, report, toolVersion)
}

func writeFindings(writer io.Writer, findings []Finding) error {
	for _, finding := range findings {
		diagnostic := policy.Diagnostic{
			RuleID:   finding.RuleID,
			Message:  finding.Message,
			Severity: finding.Severity,
			Evidence: []policy.Evidence{finding.Evidence},
		}
		if err := pitlordreport.WriteDiagnosticText(writer, diagnostic); err != nil {
			return err
		}
	}
	return nil
}

func countSeverity(diagnostics []policy.Diagnostic, severity string) int {
	count := 0
	for _, diagnostic := range diagnostics {
		if diagnostic.Severity == severity {
			count++
		}
	}
	return count
}
