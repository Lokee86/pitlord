package policy

import "github.com/Lokee86/pitlord/internal/arcana"

const Version = 1

const (
	RuleForbidDependency = "forbid_dependency"
	RuleRequireOwnership = "require_ownership"
	RuleForbidAreaCycles = "forbid_area_cycles"
	RuleForbidContent    = "forbid_content"
	RuleRequireContent   = "require_content"
	RuleRequirePath      = "require_path"
	RuleForbidPath       = "forbid_path"
)

var defaultRelations = []string{
	"imports",
	"depends-on",
	"calls",
	"references",
	"extends",
	"implements",
	"includes",
	"routes-to",
	"communicates-with",
	"observed-calls",
}

type Document struct {
	Version  int      `json:"version"`
	Includes []string `json:"includes,omitempty"`
	Areas    []Area   `json:"areas,omitempty"`
	Rules    []Rule   `json:"rules"`
}

type Area struct {
	ID           string   `json:"id"`
	Description  string   `json:"description,omitempty"`
	Paths        []string `json:"paths"`
	ExcludePaths []string `json:"exclude_paths,omitempty"`
	Kinds        []string `json:"kinds,omitempty"`
}

type Rule struct {
	ID       string `json:"id"`
	Type     string `json:"type,omitempty"`
	Message  string `json:"message,omitempty"`
	Severity string `json:"severity,omitempty"`

	FromAreas        []string `json:"from_areas,omitempty"`
	ToAreas          []string `json:"to_areas,omitempty"`
	FromPaths        []string `json:"from_paths,omitempty"`
	ToPaths          []string `json:"to_paths,omitempty"`
	FromExcludePaths []string `json:"from_exclude_paths,omitempty"`
	ToExcludePaths   []string `json:"to_exclude_paths,omitempty"`
	SourceKinds      []string `json:"source_kinds,omitempty"`
	TargetKinds      []string `json:"target_kinds,omitempty"`
	Relations        []string `json:"relations,omitempty"`
	Areas            []string `json:"areas,omitempty"`

	ScopePaths        []string `json:"scope_paths,omitempty"`
	ScopeExcludePaths []string `json:"scope_exclude_paths,omitempty"`
	AllowUnowned      bool     `json:"allow_unowned,omitempty"`
	AllowOverlaps     bool     `json:"allow_overlaps,omitempty"`

	IncludePaths []string `json:"include_paths,omitempty"`
	ExcludePaths []string `json:"exclude_paths,omitempty"`
	Literal      string   `json:"literal,omitempty"`
	Regex        string   `json:"regex,omitempty"`
	Pattern      string   `json:"pattern,omitempty"`
	PatternType  string   `json:"pattern_type,omitempty"`
	Path         string   `json:"path,omitempty"`
}

type Evidence struct {
	Issue      string       `json:"issue"`
	Source     arcana.Node  `json:"source"`
	Relation   string       `json:"relation,omitempty"`
	Target     *arcana.Node `json:"target,omitempty"`
	Areas      []string     `json:"areas,omitempty"`
	SourceArea string       `json:"source_area,omitempty"`
	TargetArea string       `json:"target_area,omitempty"`
}

type Diagnostic struct {
	RuleID   string     `json:"rule_id"`
	Message  string     `json:"message"`
	Severity string     `json:"severity"`
	Evidence []Evidence `json:"evidence"`
}

type Summary struct {
	RulesChecked         int `json:"rules_checked"`
	SourceNodesScanned   int `json:"source_nodes_scanned"`
	RelationshipsScanned int `json:"relationships_scanned"`
	Errors               int `json:"errors"`
	Warnings             int `json:"warnings"`
	SuppressedEvidence   int `json:"suppressed_evidence,omitempty"`
}

type Expectation struct {
	ExpectedIDs []string `json:"expected_ids,omitempty"`
	ActualIDs   []string `json:"actual_ids,omitempty"`
	MissingIDs  []string `json:"missing_ids,omitempty"`
	Unexpected  []string `json:"unexpected_ids,omitempty"`
	Matched     bool     `json:"matched"`
}

type Report struct {
	Schema       string       `json:"schema"`
	Snapshot     string       `json:"snapshot"`
	PolicySource string       `json:"policy_source"`
	Diagnostics  []Diagnostic `json:"diagnostics"`
	Summary      Summary      `json:"summary"`
	Expectation  *Expectation `json:"expectation,omitempty"`
}
