package mutation

import "time"

const ContractVersion = 1

type Manifest struct {
	Version                      int                  `json:"version"`
	MutationID                   string               `json:"mutation_id"`
	Kind                         string               `json:"kind"`
	Language                     string               `json:"language"`
	Summary                      string               `json:"summary"`
	SourceRepository             string               `json:"source_repository"`
	SourceWasDirty               bool                 `json:"source_was_dirty"`
	BaseCommit                   string               `json:"base_commit"`
	Worktree                     string               `json:"worktree"`
	Edits                        []Edit               `json:"edits"`
	ExpectedAddedRelationships   []Relationship       `json:"expected_added_relationships"`
	ExpectedRemovedRelationships []Relationship       `json:"expected_removed_relationships"`
	Selection                    *SelectionEvidence   `json:"selection,omitempty"`
	Verification                 []VerificationResult `json:"verification"`
	CreatedAt                    time.Time            `json:"created_at"`
}

type Relationship struct {
	Source     string `json:"source"`
	Relation   string `json:"relation"`
	Target     string `json:"target"`
	SourceArea string `json:"source_area,omitempty"`
	TargetArea string `json:"target_area,omitempty"`
}

type Edit struct {
	Path         string `json:"path"`
	BeforeSHA256 string `json:"before_sha256,omitempty"`
	AfterSHA256  string `json:"after_sha256"`
	Created      bool   `json:"created,omitempty"`
}

type SelectionEvidence struct {
	Mode                 string   `json:"mode"`
	Seed                 *int64   `json:"seed,omitempty"`
	GroupID              string   `json:"group_id,omitempty"`
	BoundaryCandidates   int      `json:"boundary_candidates"`
	EligibleCandidates   int      `json:"eligible_candidates"`
	SelectedCandidateIDs []string `json:"selected_candidate_ids"`
	Depth                int      `json:"depth,omitempty"`
	GeneratedSymbols     []string `json:"generated_symbols,omitempty"`
}

type VerificationResult struct {
	Name   string   `json:"name"`
	Args   []string `json:"args,omitempty"`
	Dir    string   `json:"dir,omitempty"`
	Output string   `json:"output,omitempty"`
}

type Check struct {
	Expectation     string       `json:"expectation"`
	Relationship    Relationship `json:"relationship"`
	BaselinePresent bool         `json:"baseline_present"`
	MutatedPresent  bool         `json:"mutated_present"`
	Matched         bool         `json:"matched"`
}

type Result struct {
	Schema           string  `json:"schema"`
	ManifestPath     string  `json:"manifest_path"`
	MutationID       string  `json:"mutation_id"`
	Kind             string  `json:"kind"`
	Language         string  `json:"language"`
	BaseCommit       string  `json:"base_commit"`
	BaselineSnapshot string  `json:"baseline_snapshot"`
	MutatedSnapshot  string  `json:"mutated_snapshot"`
	Checks           []Check `json:"checks"`
	Matched          bool    `json:"matched"`
}
