package scan

import "sort"

const Schema = "pitlord.scan.v1"

type Severity string

type Disposition string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"

	DispositionAdvisory Disposition = "advisory"
	DispositionGuard    Disposition = "guard"
)

type Scope struct {
	Kind string `json:"kind"`
	ID   string `json:"id,omitempty"`
	Path string `json:"path,omitempty"`
	Name string `json:"name,omitempty"`
}

type Evidence struct {
	Kind    string `json:"kind"`
	Message string `json:"message"`
}

type Finding struct {
	ID                string      `json:"id"`
	Detector          string      `json:"detector"`
	Disposition       Disposition `json:"disposition"`
	Severity          Severity    `json:"severity"`
	Scope             Scope       `json:"scope"`
	Summary           string      `json:"summary"`
	Rationale         string      `json:"rationale"`
	Evidence          []Evidence  `json:"evidence"`
	RequiredOutcome   string      `json:"required_outcome"`
	RecommendedAction string      `json:"recommended_action"`
}

type Summary struct {
	FindingCount int `json:"finding_count"`
	Advisory     int `json:"advisory"`
	Guard        int `json:"guard"`
	Info         int `json:"info"`
	Warnings     int `json:"warnings"`
	High         int `json:"high"`
	Critical     int `json:"critical"`
}

type Result struct {
	Schema   string    `json:"schema"`
	Scope    Scope     `json:"scope"`
	Findings []Finding `json:"findings"`
	Summary  Summary   `json:"summary"`
}

func finalize(result Result) Result {
	if result.Findings == nil {
		result.Findings = []Finding{}
	}
	for index := range result.Findings {
		if result.Findings[index].Evidence == nil {
			result.Findings[index].Evidence = []Evidence{}
		}
	}
	sort.SliceStable(result.Findings, func(left, right int) bool {
		lhs := result.Findings[left]
		rhs := result.Findings[right]
		if dispositionRank(lhs.Disposition) != dispositionRank(rhs.Disposition) {
			return dispositionRank(lhs.Disposition) > dispositionRank(rhs.Disposition)
		}
		if severityRank(lhs.Severity) != severityRank(rhs.Severity) {
			return severityRank(lhs.Severity) > severityRank(rhs.Severity)
		}
		if lhs.Detector != rhs.Detector {
			return lhs.Detector < rhs.Detector
		}
		if lhs.Scope.Path != rhs.Scope.Path {
			return lhs.Scope.Path < rhs.Scope.Path
		}
		return lhs.ID < rhs.ID
	})
	result.Summary = summarize(result.Findings)
	return result
}

func summarize(findings []Finding) Summary {
	summary := Summary{FindingCount: len(findings)}
	for _, finding := range findings {
		switch finding.Disposition {
		case DispositionAdvisory:
			summary.Advisory++
		case DispositionGuard:
			summary.Guard++
		}
		switch finding.Severity {
		case SeverityInfo:
			summary.Info++
		case SeverityWarning:
			summary.Warnings++
		case SeverityHigh:
			summary.High++
		case SeverityCritical:
			summary.Critical++
		}
	}
	return summary
}

func dispositionRank(disposition Disposition) int {
	switch disposition {
	case DispositionGuard:
		return 2
	case DispositionAdvisory:
		return 1
	default:
		return 0
	}
}

func severityRank(severity Severity) int {
	switch severity {
	case SeverityCritical:
		return 4
	case SeverityHigh:
		return 3
	case SeverityWarning:
		return 2
	case SeverityInfo:
		return 1
	default:
		return 0
	}
}
