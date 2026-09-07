package calibration

import (
	"sort"

	"github.com/Lokee86/pitlord/internal/scan"
)

func Evaluate(reference Reference, scanResult scan.Result) Result {
	items := make([]Item, 0, len(reference.Expectations))
	matched := make(map[int]struct{})
	summary := Summary{Expectations: len(reference.Expectations)}

	expectations := append([]Expectation(nil), reference.Expectations...)
	sort.Slice(expectations, func(i, j int) bool { return expectations[i].ID < expectations[j].ID })
	for _, expectation := range expectations {
		findings, indexes := matchFindings(reference, expectation, scanResult.Findings)
		for _, index := range indexes {
			matched[index] = struct{}{}
		}
		item := Item{
			ID:              expectation.ID,
			Class:           expectation.Class,
			Finding:         expectation.Finding,
			MatchedFindings: findings,
		}
		item.Outcome = classifyOutcome(expectation.Finding, len(findings) > 0)
		item.SeverityMismatch = severityMismatch(expectation, findings)
		accumulate(&summary, item)
		items = append(items, item)
	}

	unlabelled := make([]scan.Finding, 0)
	for index, finding := range scanResult.Findings {
		if !matchesReference(reference, finding) {
			continue
		}
		if _, ok := matched[index]; !ok {
			unlabelled = append(unlabelled, finding)
		}
	}
	summary.UnlabelledFindings = len(unlabelled)
	summary.LabeledPrecision = ratio(summary.TruePositive, summary.TruePositive+summary.FalsePositive)
	summary.LabeledRecall = ratio(summary.TruePositive, summary.TruePositive+summary.FalseNegative)

	return Result{
		Schema:               Schema,
		Corpus:               reference.Corpus,
		SourceRevision:       reference.SourceRevision,
		Detector:             reference.Detector,
		Analyzer:             reference.Analyzer,
		RuleID:               reference.RuleID,
		Language:             reference.Language,
		RequireFullyLabelled: reference.RequireFullyLabelled,
		Items:                items,
		UnlabelledFindings:   unlabelled,
		Summary:              summary,
	}
}

func classifyOutcome(expectation FindingExpectation, found bool) Outcome {
	switch expectation {
	case FindingRequired:
		if found {
			return OutcomeTruePositive
		}
		return OutcomeFalseNegative
	case FindingAbsent:
		if found {
			return OutcomeFalsePositive
		}
		return OutcomeTrueNegative
	default:
		if found {
			return OutcomeAllowedFound
		}
		return OutcomeAllowedQuiet
	}
}

func severityMismatch(expectation Expectation, findings []scan.Finding) bool {
	if len(findings) == 0 || expectation.Finding == FindingAbsent {
		return false
	}
	for _, finding := range findings {
		if expectation.MinSeverity != "" && severityRank(finding.Severity) < severityRank(expectation.MinSeverity) {
			return true
		}
		if expectation.MaxSeverity != "" && severityRank(finding.Severity) > severityRank(expectation.MaxSeverity) {
			return true
		}
	}
	return false
}

func accumulate(summary *Summary, item Item) {
	switch item.Finding {
	case FindingRequired:
		summary.Required++
	case FindingAbsent:
		summary.Absent++
	case FindingAllowed:
		summary.Allowed++
	}
	switch item.Outcome {
	case OutcomeTruePositive:
		summary.TruePositive++
	case OutcomeTrueNegative:
		summary.TrueNegative++
	case OutcomeFalsePositive:
		summary.FalsePositive++
	case OutcomeFalseNegative:
		summary.FalseNegative++
	}
	if item.SeverityMismatch {
		summary.SeverityMismatches++
	}
}

func ratio(numerator, denominator int) float64 {
	if denominator == 0 {
		return 1
	}
	return float64(numerator) / float64(denominator)
}

func severityRank(severity scan.Severity) int {
	switch severity {
	case scan.SeverityCritical:
		return 4
	case scan.SeverityHigh:
		return 3
	case scan.SeverityWarning:
		return 2
	case scan.SeverityInfo:
		return 1
	default:
		return 0
	}
}
