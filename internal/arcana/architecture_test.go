package arcana

import (
	"context"
	"encoding/json"
	"testing"
)

func TestInspectArchitectureDecodesSummary(t *testing.T) {
	run := func(_ context.Context, _, _ string, requests []request) (map[string]response, error) {
		if len(requests) != 1 {
			t.Fatalf("expected one request, got %d", len(requests))
		}
		current := requests[0]
		if current.Op != "architecture_summary" {
			t.Fatalf("unexpected operation %q", current.Op)
		}
		if current.PathPrefix != "" {
			t.Fatalf("expected root prefix to be omitted, got %q", current.PathPrefix)
		}
		if current.MinCommunitySize != 3 || current.Limit != 7 {
			t.Fatalf("unexpected bounds: %+v", current)
		}
		payload, _ := json.Marshal(ArchitectureSummary{
			Relations:         []string{"calls"},
			NodeCount:         12,
			InternalEdgeCount: 8,
			CommunityCount:    1,
			Returned:          1,
			Communities: []ArchitectureCommunity{
				{
					CommunityID: 1,
					NodeCount:   5,
					EdgeCount:   4,
					Paths:       []CommunityPath{{Path: "api", NodeCount: 5}},
				},
			},
		})
		return map[string]response{current.ID: {OK: true, Result: payload}}, nil
	}

	summary, err := (Client{Run: run}).InspectArchitecture(
		context.Background(),
		"snapshot",
		".",
		[]string{"calls"},
		3,
		7,
	)
	if err != nil {
		t.Fatal(err)
	}
	if summary.NodeCount != 12 || len(summary.Communities) != 1 {
		t.Fatalf("unexpected summary: %+v", summary)
	}
}

func TestInspectArchitectureRejectsInconsistentCommunityCount(t *testing.T) {
	run := func(_ context.Context, _, _ string, requests []request) (map[string]response, error) {
		payload, _ := json.Marshal(ArchitectureSummary{Returned: 2, Communities: []ArchitectureCommunity{{}}})
		return map[string]response{requests[0].ID: {OK: true, Result: payload}}, nil
	}
	if _, err := (Client{Run: run}).InspectArchitecture(context.Background(), "snapshot", "src", nil, 2, 20); err == nil {
		t.Fatal("expected inconsistent community count to fail")
	}
}
