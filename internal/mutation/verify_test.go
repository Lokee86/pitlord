package mutation

import "testing"

func TestRelationshipExpectationSemantics(t *testing.T) {
	cases := []struct {
		expectation string
		baseline    bool
		mutated     bool
		matched     bool
	}{
		{"added", false, true, true},
		{"added", true, true, false},
		{"added", false, false, false},
		{"removed", true, false, true},
		{"removed", true, true, false},
		{"removed", false, false, false},
	}
	for _, current := range cases {
		matched := false
		switch current.expectation {
		case "added":
			matched = !current.baseline && current.mutated
		case "removed":
			matched = current.baseline && !current.mutated
		}
		if matched != current.matched {
			t.Fatalf("unexpected match for %+v", current)
		}
	}
}
