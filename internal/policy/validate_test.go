package policy

import "testing"

func TestValidationRejectsUnknownAreaReference(t *testing.T) {
	document := Document{
		Version: Version,
		Areas:   []Area{{ID: "api", Paths: []string{"api"}}},
		Rules: []Rule{
			{ID: "bad", FromAreas: []string{"api"}, ToAreas: []string{"missing"}},
		},
	}
	if err := document.NormalizeAndValidate(); err == nil {
		t.Fatal("expected unknown area reference to fail")
	}
}

func TestValidationRejectsOwnershipRuleThatEnforcesNothing(t *testing.T) {
	document := Document{
		Version: Version,
		Areas:   []Area{{ID: "api", Paths: []string{"api"}}},
		Rules: []Rule{
			{
				ID:            "bad",
				Type:          RuleRequireOwnership,
				ScopePaths:    []string{"."},
				AllowUnowned:  true,
				AllowOverlaps: true,
			},
		},
	}
	if err := document.NormalizeAndValidate(); err == nil {
		t.Fatal("expected no-op ownership rule to fail")
	}
}

func TestValidationDefaultsAreaCycleSelection(t *testing.T) {
	document := Document{
		Version: Version,
		Areas: []Area{
			{ID: "z", Paths: []string{"z"}},
			{ID: "a", Paths: []string{"a"}},
		},
		Rules: []Rule{{ID: "cycles", Type: RuleForbidAreaCycles}},
	}
	if err := document.NormalizeAndValidate(); err != nil {
		t.Fatal(err)
	}
	if len(document.Rules[0].Areas) != 2 || document.Rules[0].Areas[0] != "a" || document.Rules[0].Areas[1] != "z" {
		t.Fatalf("unexpected selected areas: %v", document.Rules[0].Areas)
	}
	if len(document.Rules[0].Relations) == 0 {
		t.Fatal("expected default dependency relations")
	}
}

func TestValidationRejectsUnknownAreaCycleSelection(t *testing.T) {
	document := Document{
		Version: Version,
		Areas: []Area{
			{ID: "a", Paths: []string{"a"}},
			{ID: "b", Paths: []string{"b"}},
		},
		Rules: []Rule{{ID: "cycles", Type: RuleForbidAreaCycles, Areas: []string{"a", "missing"}}},
	}
	if err := document.NormalizeAndValidate(); err == nil {
		t.Fatal("expected unknown cycle area to fail")
	}
}

func TestValidationRejectsSingleAreaCycleSelection(t *testing.T) {
	document := Document{
		Version: Version,
		Areas: []Area{
			{ID: "a", Paths: []string{"a"}},
			{ID: "b", Paths: []string{"b"}},
		},
		Rules: []Rule{{ID: "cycles", Type: RuleForbidAreaCycles, Areas: []string{"a"}}},
	}
	if err := document.NormalizeAndValidate(); err == nil {
		t.Fatal("expected one-area cycle rule to fail")
	}
}

func TestValidationDefaultsRuleTypesAndOwnershipKinds(t *testing.T) {
	document := Document{
		Version: Version,
		Areas:   []Area{{ID: "api", Paths: []string{"api"}}},
		Rules: []Rule{
			{ID: "dependency", FromPaths: []string{"api"}, ToPaths: []string{"storage"}},
			{ID: "ownership", Type: RuleRequireOwnership, ScopePaths: []string{"api"}},
		},
	}
	if err := document.NormalizeAndValidate(); err != nil {
		t.Fatal(err)
	}
	if document.Rules[0].Type != RuleForbidDependency {
		t.Fatalf("unexpected dependency rule type %q", document.Rules[0].Type)
	}
	if len(document.Rules[1].SourceKinds) != 1 || document.Rules[1].SourceKinds[0] != "file" {
		t.Fatalf("unexpected ownership kinds: %v", document.Rules[1].SourceKinds)
	}
}
