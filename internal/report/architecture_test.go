package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestWriteArchitectureText(t *testing.T) {
	summary := arcana.ArchitectureSummary{
		Relations:         []string{"calls", "imports"},
		NodeCount:         10,
		InternalEdgeCount: 7,
		BoundaryEdgeCount: 3,
		CommunityCount:    1,
		Returned:          1,
		Communities: []arcana.ArchitectureCommunity{
			{
				NodeCount:           5,
				EdgeCount:           4,
				Paths:               []arcana.CommunityPath{{Path: "api", NodeCount: 5}},
				RelationCounts:      map[string]int{"calls": 4},
				RepresentativeNodes: []arcana.Node{{Name: "CreateUser", Path: "api/user.go"}},
			},
		},
	}
	var output bytes.Buffer
	if err := WriteArchitectureText(&output, summary); err != nil {
		t.Fatal(err)
	}
	text := output.String()
	for _, expected := range []string{
		"Architecture: 10 nodes, 7 internal edges, 3 boundary edges",
		"Paths: api (5)",
		"CreateUser [api/user.go]",
		"Internal: calls=4",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("expected output to contain %q:\n%s", expected, text)
		}
	}
}
