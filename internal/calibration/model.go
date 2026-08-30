package calibration

import "github.com/Lokee86/pitlord/internal/scan"

const Schema = "pitlord.calibration.v1"

type FindingExpectation string

const (
	FindingRequired FindingExpectation = "required"
	FindingAbsent   FindingExpectation = "absent"
	FindingAllowed  FindingExpectation = "allowed"
)

type Expectation struct {
	ID          string             `json:"id"`
	Path        string             `json:"path,omitempty"`
	PathPrefix  string             `json:"path_prefix,omitempty"`
	Class       string             `json:"class"`
	Finding     FindingExpectation `json:"finding"`
	MinSeverity scan.Severity      `json:"min_severity,omitempty"`
	MaxSeverity scan.Severity      `json:"max_severity,omitempty"`
}

type Reference struct {
	Schema             string        `json:"schema"`
	Corpus             string        `json:"corpus"`
	SourceRevision     string        `json:"source_revision"`
	WorktreeDiffSHA256 string        `json:"worktree_diff_sha256,omitempty"`
	Detector           string        `json:"detector"`
	PathPrefix         string        `json:"path_prefix,omitempty"`
	Expectations       []Expectation `json:"expectations"`
}

type Outcome string

const (
	OutcomeTruePositive  Outcome = "true-positive"
	OutcomeTrueNegative  Outcome = "true-negative"
	OutcomeFalsePositive Outcome = "false-positive"
	OutcomeFalseNegative Outcome = "false-negative"
	OutcomeAllowedFound  Outcome = "allowed-found"
	OutcomeAllowedQuiet  Outcome = "allowed-quiet"
)

type Item struct {
	ID               string             `json:"id"`
	Class            string             `json:"class"`
	Finding          FindingExpectation `json:"finding"`
	Outcome          Outcome            `json:"outcome"`
	MatchedFindings  []scan.Finding     `json:"matched_findings"`
	SeverityMismatch bool               `json:"severity_mismatch,omitempty"`
}

type Summary struct {
	Expectations       int     `json:"expectations"`
	Required           int     `json:"required"`
	Absent             int     `json:"absent"`
	Allowed            int     `json:"allowed"`
	TruePositive       int     `json:"true_positive"`
	TrueNegative       int     `json:"true_negative"`
	FalsePositive      int     `json:"false_positive"`
	FalseNegative      int     `json:"false_negative"`
	SeverityMismatches int     `json:"severity_mismatches"`
	UnlabelledFindings int     `json:"unlabelled_findings"`
	LabeledPrecision   float64 `json:"labeled_precision"`
	LabeledRecall      float64 `json:"labeled_recall"`
}

type Result struct {
	Schema             string         `json:"schema"`
	Corpus             string         `json:"corpus"`
	SourceRevision     string         `json:"source_revision"`
	Detector           string         `json:"detector"`
	Items              []Item         `json:"items"`
	UnlabelledFindings []scan.Finding `json:"unlabelled_findings"`
	Summary            Summary        `json:"summary"`
}

func (result Result) HasMismatch() bool {
	return result.Summary.FalsePositive > 0 ||
		result.Summary.FalseNegative > 0 ||
		result.Summary.SeverityMismatches > 0
}
