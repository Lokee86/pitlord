package calibration

import (
	"path"
	"strings"

	"github.com/Lokee86/pitlord/internal/scan"
)

func matchFindings(reference Reference, expectation Expectation, findings []scan.Finding) ([]scan.Finding, []int) {
	matched := make([]scan.Finding, 0)
	indexes := make([]int, 0)
	for index, finding := range findings {
		if !matchesReference(reference, finding) || !matchesPath(expectation, finding.Scope.Path) || !matchesLocation(expectation, finding.Location) {
			continue
		}
		matched = append(matched, finding)
		indexes = append(indexes, index)
	}
	return matched, indexes
}

func matchesReference(reference Reference, finding scan.Finding) bool {
	if detector := strings.TrimSpace(reference.Detector); detector != "" && finding.Detector != detector {
		return false
	}
	if analyzer := strings.TrimSpace(reference.Analyzer); analyzer != "" && finding.Analyzer != analyzer {
		return false
	}
	if ruleID := strings.TrimSpace(reference.RuleID); ruleID != "" && finding.RuleID != ruleID {
		return false
	}
	if language := strings.TrimSpace(reference.Language); language != "" && finding.Language != language {
		return false
	}
	return true
}

func matchesLocation(expectation Expectation, actual *scan.SourceSpan) bool {
	expected := expectation.Location
	if expected == nil {
		return true
	}
	if actual == nil || normalizePath(actual.Path) != normalizePath(expectation.Path) || actual.StartLine != expected.StartLine {
		return false
	}
	if expected.StartColumn > 0 && actual.StartColumn != expected.StartColumn {
		return false
	}
	if expected.EndLine > 0 && actual.EndLine != expected.EndLine {
		return false
	}
	if expected.EndColumn > 0 && actual.EndColumn != expected.EndColumn {
		return false
	}
	return true
}

func matchesPath(expectation Expectation, value string) bool {
	value = normalizePath(value)
	if expectation.Path != "" {
		return value == normalizePath(expectation.Path)
	}
	prefix := strings.TrimSuffix(normalizePath(expectation.PathPrefix), "/")
	if prefix == "." {
		return true
	}
	return value == prefix || strings.HasPrefix(value, prefix+"/")
}

func normalizePath(value string) string {
	value = strings.ReplaceAll(strings.TrimSpace(value), "\\", "/")
	return strings.TrimPrefix(path.Clean(value), "./")
}
