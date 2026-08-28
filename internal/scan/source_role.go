package scan

import "strings"

type sourceRole string

const (
	sourceRoleProduction    sourceRole = "production"
	sourceRoleNonProduction sourceRole = "non-production"
)

func classifySourceRole(filePath string) sourceRole {
	normalized := strings.ToLower(normalizedUnitPath(filePath))
	segments := strings.Split(normalized, "/")
	for _, segment := range segments {
		if nonProductionSegment(segment) {
			return sourceRoleNonProduction
		}
	}
	if strings.Contains(normalized, "/snippets/docs/") || strings.HasPrefix(normalized, "snippets/docs/") {
		return sourceRoleNonProduction
	}
	return sourceRoleProduction
}

func nonProductionSegment(segment string) bool {
	switch segment {
	case "test", "tests", "testing", "testdata", "fixtures", "bench", "benchmark", "benchmarks", "generated", "vendor", "examples", "samples", "tooling", "devtools":
		return true
	}
	return strings.HasSuffix(segment, ".tests") ||
		strings.HasSuffix(segment, "-tests") ||
		strings.HasPrefix(segment, "test-") ||
		strings.HasPrefix(segment, "test_") ||
		strings.Contains(segment, ".test.") ||
		strings.Contains(segment, ".spec.") ||
		strings.Contains(segment, "_test.") ||
		strings.HasSuffix(segment, "-testing") ||
		strings.HasSuffix(segment, "-samples") ||
		strings.HasSuffix(segment, "-examples") ||
		strings.HasSuffix(segment, "-benchmarks")
}
