package policy

import (
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestEvaluateGroupsMatchingRelationshipsByRule(t *testing.T) {
	document := Document{
		Version: Version,
		Rules: []Rule{
			{
				ID:        "api-must-not-access-storage",
				FromPaths: []string{"api"},
				ToPaths:   []string{"storage"},
				Relations: []string{"imports", "calls"},
			},
		},
	}
	if err := document.NormalizeAndValidate(); err != nil {
		t.Fatal(err)
	}

	apiFile := arcana.Node{NodeID: 1, Kind: "file", Path: "api/handler.go", Name: "api"}
	apiFunction := arcana.Node{NodeID: 2, Kind: "function", Path: "api/handler.go", Name: "CreateUser"}
	storageFile := arcana.Node{NodeID: 3, Kind: "file", Path: "storage/users.go", Name: "storage"}
	storageFunction := arcana.Node{NodeID: 4, Kind: "function", Path: "storage/users.go", Name: "InsertUser"}
	serviceFunction := arcana.Node{NodeID: 5, Kind: "function", Path: "service/users.go", Name: "CreateUser"}
	graph := arcana.Graph{
		Sources: []arcana.Node{apiFile, apiFunction},
		Outgoing: map[uint32][]arcana.Relationship{
			1: {
				{Relation: "imports", Node: storageFile},
			},
			2: {
				{Relation: "calls", Node: storageFunction},
				{Relation: "calls", Node: serviceFunction},
			},
		},
		Relationships: 3,
	}

	diagnostics := Evaluate(document, graph)
	if len(diagnostics) != 1 {
		t.Fatalf("expected one diagnostic, got %d", len(diagnostics))
	}
	if diagnostics[0].RuleID != "api-must-not-access-storage" {
		t.Fatalf("unexpected rule id %q", diagnostics[0].RuleID)
	}
	if len(diagnostics[0].Evidence) != 2 {
		t.Fatalf("expected two evidence edges, got %d", len(diagnostics[0].Evidence))
	}
	summary := BuildSummary(document, graph, diagnostics)
	if summary.Errors != 1 || summary.RelationshipsScanned != 3 {
		t.Fatalf("unexpected summary: %+v", summary)
	}
}

func TestEvaluateHonorsExclusionsAndKindFilters(t *testing.T) {
	document := Document{
		Version: Version,
		Rules: []Rule{
			{
				ID:               "runtime-must-not-call-storage",
				FromPaths:        []string{"."},
				FromExcludePaths: []string{"tests"},
				ToPaths:          []string{"storage"},
				SourceKinds:      []string{"function"},
				TargetKinds:      []string{"function"},
				Relations:        []string{"calls"},
			},
		},
	}
	if err := document.NormalizeAndValidate(); err != nil {
		t.Fatal(err)
	}

	mainFunction := arcana.Node{NodeID: 1, Kind: "function", Path: "runtime/run.go", Name: "Run"}
	testFunction := arcana.Node{NodeID: 2, Kind: "function", Path: "tests/run_test.go", Name: "TestRun"}
	mainFile := arcana.Node{NodeID: 3, Kind: "file", Path: "runtime/run.go", Name: "runtime"}
	storageFunction := arcana.Node{NodeID: 4, Kind: "function", Path: "storage/store.go", Name: "Write"}
	storageFile := arcana.Node{NodeID: 5, Kind: "file", Path: "storage/store.go", Name: "storage"}
	graph := arcana.Graph{
		Sources: []arcana.Node{mainFunction, testFunction, mainFile},
		Outgoing: map[uint32][]arcana.Relationship{
			1: {{Relation: "calls", Node: storageFunction}},
			2: {{Relation: "calls", Node: storageFunction}},
			3: {{Relation: "calls", Node: storageFile}},
		},
		Relationships: 3,
	}

	diagnostics := Evaluate(document, graph)
	if len(diagnostics) != 1 {
		t.Fatalf("expected one diagnostic, got %d", len(diagnostics))
	}
	if len(diagnostics[0].Evidence) != 1 {
		t.Fatalf("expected one filtered evidence edge, got %d", len(diagnostics[0].Evidence))
	}
	if diagnostics[0].Evidence[0].Source.NodeID != mainFunction.NodeID {
		t.Fatalf("unexpected source: %+v", diagnostics[0].Evidence[0].Source)
	}
}

func TestCompareExpectationReportsMissingAndUnexpectedIDs(t *testing.T) {
	diagnostics := []Diagnostic{{RuleID: "actual"}}
	comparison := CompareExpectation([]string{"expected"}, diagnostics)
	if comparison == nil || comparison.Matched {
		t.Fatal("expected an expectation mismatch")
	}
	if len(comparison.MissingIDs) != 1 || comparison.MissingIDs[0] != "expected" {
		t.Fatalf("unexpected missing ids: %v", comparison.MissingIDs)
	}
	if len(comparison.Unexpected) != 1 || comparison.Unexpected[0] != "actual" {
		t.Fatalf("unexpected extra ids: %v", comparison.Unexpected)
	}
}
