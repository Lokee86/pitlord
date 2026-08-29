package scan

import "testing"

func TestBoundaryBypassTreatsNestedDirectoriesAsOneOwnershipTree(t *testing.T) {
	for _, test := range []struct {
		source string
		target string
		want   bool
	}{
		{source: "pkg/internal/bind/adapter.go", target: "pkg/internal/helper.go", want: true},
		{source: "game/root.go", target: "game/entities/pickup.go", want: true},
		{source: "ui/view.go", target: "storage/store.go", want: false},
		{source: "main.go", target: "storage/store.go", want: false},
	} {
		if got := nestedFileOwnership(test.source, test.target); got != test.want {
			t.Errorf("nestedFileOwnership(%q, %q) = %v, want %v", test.source, test.target, got, test.want)
		}
	}
}
