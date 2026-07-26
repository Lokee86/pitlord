package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
	"github.com/Lokee86/pitlord/internal/policy"
)

func TestWriteAreaAnalysisText(t *testing.T) {
	target := arcana.Node{Name: "Write", Path: "storage/store.go"}
	analysis := policy.AreaAnalysis{
		Areas: []policy.AreaStatistics{
			{ID: "api", NodeCount: 4, KindCounts: map[string]int{"file": 1, "function": 3}},
			{ID: "storage", NodeCount: 2, KindCounts: map[string]int{"file": 1, "function": 1}},
		},
		Ownership: policy.OwnershipStatistics{CheckedNodes: 3, OwnedNodes: 2, UnownedNodes: 1},
		Dependencies: []policy.AreaDependency{
			{
				SourceArea:     "api",
				TargetArea:     "storage",
				EdgeCount:      2,
				RelationCounts: map[string]int{"calls": 2},
				Representative: policy.Evidence{
					Issue:      "area_dependency",
					Source:     arcana.Node{Name: "Create", Path: "api/create.go"},
					Relation:   "calls",
					Target:     &target,
					SourceArea: "api",
					TargetArea: "storage",
				},
			},
		},
		Cycles:  [][]string{{"api", "storage"}},
		Summary: policy.AreaAnalysisSummary{AreaCount: 2, DependencyCount: 1, CycleCount: 1},
	}
	var output bytes.Buffer
	if err := WriteAreaAnalysisText(&output, analysis); err != nil {
		t.Fatal(err)
	}
	text := output.String()
	for _, expected := range []string{
		"2 areas, 1 cross-area/external dependencies, 1 cycle(s)",
		"api -> storage: 2 edges (calls=2)",
		"api <-> storage",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("expected %q in output:\n%s", expected, text)
		}
	}
}

func TestWriteAreaAnalysisDOT(t *testing.T) {
	analysis := policy.AreaAnalysis{
		Areas: []policy.AreaStatistics{{ID: "api", NodeCount: 4}, {ID: "storage", NodeCount: 2}},
		Dependencies: []policy.AreaDependency{
			{SourceArea: "api", TargetArea: "storage", EdgeCount: 5, RelationCounts: map[string]int{"calls": 5}},
			{SourceArea: "api", TargetArea: policy.ExternalArea, EdgeCount: 1, RelationCounts: map[string]int{"imports": 1}},
		},
		Cycles: [][]string{{"api", "storage"}},
	}
	var output bytes.Buffer
	if err := WriteAreaAnalysisDOT(&output, analysis); err != nil {
		t.Fatal(err)
	}
	text := output.String()
	for _, expected := range []string{"digraph pitlord_architecture", "color=red", "external", "calls=5"} {
		if !strings.Contains(text, expected) {
			t.Fatalf("expected %q in DOT:\n%s", expected, text)
		}
	}
}
