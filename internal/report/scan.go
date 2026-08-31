package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/Lokee86/pitlord/internal/scan"
)

func WriteScanJSON(writer io.Writer, result scan.Result) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func WriteScanText(writer io.Writer, result scan.Result) error {
	if len(result.Findings) == 0 {
		_, err := fmt.Fprintf(
			writer,
			"Pitlord scan: no findings in %s.\n",
			scanScope(result.Scope),
		)
		return err
	}
	for _, finding := range result.Findings {
		if _, err := fmt.Fprintf(
			writer,
			"%s %s %s: %s\n",
			strings.ToUpper(string(finding.Severity)),
			strings.ToUpper(string(finding.Disposition)),
			scanRuleLabel(finding),
			finding.Summary,
		); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(writer, "  Scope: %s\n", scanScope(finding.Scope)); err != nil {
			return err
		}
		if finding.Location != nil {
			if _, err := fmt.Fprintf(writer, "  Location: %s\n", scanLocation(*finding.Location)); err != nil {
				return err
			}
		}
		if finding.Rationale != "" {
			if _, err := fmt.Fprintf(writer, "  Why: %s\n", finding.Rationale); err != nil {
				return err
			}
		}
		for _, evidence := range finding.Evidence {
			if _, err := fmt.Fprintf(writer, "  Evidence: %s\n", evidence.Message); err != nil {
				return err
			}
		}
		if finding.RequiredOutcome != "" {
			if _, err := fmt.Fprintf(writer, "  Required outcome: %s\n", finding.RequiredOutcome); err != nil {
				return err
			}
		}
		if finding.RecommendedAction != "" {
			if _, err := fmt.Fprintf(writer, "  Action: %s\n", finding.RecommendedAction); err != nil {
				return err
			}
		}
		if finding.SuggestedFix != nil {
			if _, err := fmt.Fprintf(writer, "  Fix%s: %s\n", scanApplicability(finding.SuggestedFix.Applicability), finding.SuggestedFix.Message); err != nil {
				return err
			}
		}
	}
	_, err := fmt.Fprintf(
		writer,
		"Pitlord scan: %d findings (%d guard, %d advisory; %d critical, %d high, %d warnings, %d info).\n",
		result.Summary.FindingCount,
		result.Summary.Guard,
		result.Summary.Advisory,
		result.Summary.Critical,
		result.Summary.High,
		result.Summary.Warnings,
		result.Summary.Info,
	)
	return err
}

func scanRuleLabel(finding scan.Finding) string {
	if finding.RuleID != "" {
		return finding.RuleID
	}
	return finding.Detector
}

func scanLocation(span scan.SourceSpan) string {
	location := span.Path
	if span.StartLine > 0 {
		location = fmt.Sprintf("%s:%d", location, span.StartLine)
		if span.StartColumn > 0 {
			location = fmt.Sprintf("%s:%d", location, span.StartColumn)
		}
	}
	return location
}

func scanApplicability(applicability scan.Applicability) string {
	if applicability == "" {
		return ""
	}
	return fmt.Sprintf(" (%s)", applicability)
}

func scanScope(scope scan.Scope) string {
	if scope.Path != "" {
		return scope.Path
	}
	if scope.Name != "" {
		return scope.Name
	}
	if scope.ID != "" {
		return scope.ID
	}
	if scope.Kind != "" {
		return scope.Kind
	}
	return "."
}
