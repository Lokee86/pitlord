package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
	"github.com/Lokee86/pitlord/internal/policy"
)

func TestWriteTextIncludesStableDiagnosticAndEvidence(t *testing.T) {
	target := arcana.Node{Name: "InsertUser", Path: "storage/users.go"}
	value := policy.Report{
		Diagnostics: []policy.Diagnostic{
			{
				RuleID:   "api-must-not-access-storage",
				Message:  "api must not access storage",
				Severity: "error",
				Evidence: []policy.Evidence{
					{
						Issue:    "forbidden_dependency",
						Source:   arcana.Node{Name: "CreateUser", Path: "api/handler.go", Span: &arcana.Span{Path: "api/handler.go", StartLine: 12, StartColumn: 1}},
						Relation: "calls",
						Target:   &target,
					},
				},
			},
		},
		Summary: policy.Summary{RulesChecked: 1, SourceNodesScanned: 1, RelationshipsScanned: 1, Errors: 1},
	}
	var output bytes.Buffer
	if err := WriteText(&output, value); err != nil {
		t.Fatal(err)
	}
	text := output.String()
	for _, expected := range []string{
		"ERROR api-must-not-access-storage",
		"api/handler.go:12:1 CreateUser --calls--> InsertUser",
		"Pitlord: 1 errors, 0 warnings",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("expected output to contain %q:\n%s", expected, text)
		}
	}
}

func TestWriteTextIncludesAreaCycleEvidence(t *testing.T) {
	target := arcana.Node{Name: "B", Path: "src/b/b.go"}
	value := policy.Report{
		Diagnostics: []policy.Diagnostic{
			{
				RuleID:   "no-area-cycles",
				Message:  "areas must not cycle",
				Severity: "error",
				Evidence: []policy.Evidence{
					{
						Issue:      "area_cycle",
						Source:     arcana.Node{Name: "A", Path: "src/a/a.go"},
						Relation:   "calls",
						Target:     &target,
						Areas:      []string{"a", "b"},
						SourceArea: "a",
						TargetArea: "b",
					},
				},
			},
		},
		Summary: policy.Summary{RulesChecked: 1, Errors: 1},
	}
	var output bytes.Buffer
	if err := WriteText(&output, value); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "A [a] --calls--> B [b] closes cycle among a, b") {
		t.Fatalf("unexpected cycle output:\n%s", output.String())
	}
}

func TestWriteTextIncludesOwnershipEvidence(t *testing.T) {
	value := policy.Report{
		Diagnostics: []policy.Diagnostic{
			{
				RuleID:   "source-ownership",
				Message:  "source files need one owner",
				Severity: "error",
				Evidence: []policy.Evidence{
					{Issue: "unowned", Source: arcana.Node{Name: "shared.go", Path: "src/shared.go"}},
					{Issue: "multiple_owners", Source: arcana.Node{Name: "api.go", Path: "src/api/api.go"}, Areas: []string{"api", "root"}},
				},
			},
		},
		Summary: policy.Summary{RulesChecked: 1, SourceNodesScanned: 2, Errors: 1},
	}
	var output bytes.Buffer
	if err := WriteText(&output, value); err != nil {
		t.Fatal(err)
	}
	text := output.String()
	for _, expected := range []string{
		"src/shared.go shared.go has no declared owner",
		"src/api/api.go api.go matches multiple owners: api, root",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("expected output to contain %q:\n%s", expected, text)
		}
	}
}
