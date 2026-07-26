package policy

import "github.com/Lokee86/pitlord/internal/arcana"

const ExternalArea = "<external>"

type AreaAnalysisOptions struct {
	ScopePaths        []string
	ScopeExcludePaths []string
	OwnershipKinds    []string
	Relations         []string
	ExampleLimit      int
}

type AreaAnalysis struct {
	Schema       string              `json:"schema"`
	ScopePaths   []string            `json:"scope_paths"`
	Areas        []AreaStatistics    `json:"areas"`
	Ownership    OwnershipStatistics `json:"ownership"`
	Dependencies []AreaDependency    `json:"dependencies"`
	Cycles       [][]string          `json:"cycles"`
	Summary      AreaAnalysisSummary `json:"summary"`
}

type AreaStatistics struct {
	ID                    string         `json:"id"`
	NodeCount             int            `json:"node_count"`
	KindCounts            map[string]int `json:"kind_counts"`
	InternalRelationships int            `json:"internal_relationships"`
}

type OwnershipStatistics struct {
	CheckedNodes     int                `json:"checked_nodes"`
	OwnedNodes       int                `json:"owned_nodes"`
	UnownedNodes     int                `json:"unowned_nodes"`
	OverlappingNodes int                `json:"overlapping_nodes"`
	UnownedExamples  []arcana.Node      `json:"unowned_examples,omitempty"`
	OverlapExamples  []OwnershipOverlap `json:"overlap_examples,omitempty"`
}

type OwnershipOverlap struct {
	Node  arcana.Node `json:"node"`
	Areas []string    `json:"areas"`
}

type AreaDependency struct {
	SourceArea     string         `json:"source_area"`
	TargetArea     string         `json:"target_area"`
	EdgeCount      int            `json:"edge_count"`
	RelationCounts map[string]int `json:"relation_counts"`
	Representative Evidence       `json:"representative"`
}

type AreaAnalysisSummary struct {
	SourceNodesScanned    int `json:"source_nodes_scanned"`
	RelationshipsScanned  int `json:"relationships_scanned"`
	RelationshipsSelected int `json:"relationships_selected"`
	AreaCount             int `json:"area_count"`
	DependencyCount       int `json:"dependency_count"`
	CycleCount            int `json:"cycle_count"`
}
