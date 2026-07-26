package snapshotdiff

import "github.com/Lokee86/pitlord/internal/policy"

type Finding struct {
	Fingerprint string          `json:"fingerprint"`
	RuleID      string          `json:"rule_id"`
	Message     string          `json:"message"`
	Severity    string          `json:"severity"`
	Evidence    policy.Evidence `json:"evidence"`
}

type Summary struct {
	BeforeFindings int `json:"before_findings"`
	AfterFindings  int `json:"after_findings"`
	Introduced     int `json:"introduced"`
	Resolved       int `json:"resolved"`
	Persistent     int `json:"persistent"`
}

type Result struct {
	Schema         string    `json:"schema"`
	PolicySource   string    `json:"policy_source"`
	BeforeSnapshot string    `json:"before_snapshot"`
	AfterSnapshot  string    `json:"after_snapshot"`
	Introduced     []Finding `json:"introduced"`
	Resolved       []Finding `json:"resolved"`
	Persistent     []Finding `json:"persistent"`
	Summary        Summary   `json:"summary"`
}
