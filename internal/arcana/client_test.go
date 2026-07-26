package arcana

import (
	"context"
	"encoding/json"
	"testing"
)

func TestClientLoadsPathScopedGraph(t *testing.T) {
	calls := 0
	run := func(_ context.Context, _, _ string, requests []request) (map[string]response, error) {
		calls++
		results := make(map[string]response, len(requests))
		if calls == 1 {
			payload, _ := json.Marshal(nodeListResult{
				Count:    1,
				Returned: 1,
				Nodes: []Node{
					{NodeID: 7, Path: "api/handler.go", Name: "CreateUser"},
				},
			})
			results[requests[0].ID] = response{OK: true, Result: payload}
			return results, nil
		}
		payload, _ := json.Marshal(neighborResult{
			Node:  Node{NodeID: 7, Path: "api/handler.go", Name: "CreateUser"},
			Count: 1,
			Relationships: []Relationship{
				{Relation: "calls", Node: Node{NodeID: 9, Path: "storage/users.go", Name: "InsertUser"}},
			},
		})
		results[requests[0].ID] = response{OK: true, Result: payload}
		return results, nil
	}

	graph, err := (Client{Run: run}).LoadGraph(context.Background(), "snapshot", []string{"api"})
	if err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatalf("expected two Arcana batches, got %d", calls)
	}
	if len(graph.Sources) != 1 || graph.Relationships != 1 {
		t.Fatalf("unexpected graph: %+v", graph)
	}
	if graph.Outgoing[7][0].Node.Path != "storage/users.go" {
		t.Fatalf("unexpected relationship: %+v", graph.Outgoing[7][0])
	}
}

func TestClientCanLoadNodesWithoutOutgoingRelationships(t *testing.T) {
	calls := 0
	run := func(_ context.Context, _, _ string, requests []request) (map[string]response, error) {
		calls++
		if len(requests) != 1 || requests[0].Op != "list_nodes" {
			t.Fatalf("unexpected requests: %+v", requests)
		}
		payload, _ := json.Marshal(nodeListResult{
			Count:    2,
			Returned: 2,
			Nodes: []Node{
				{NodeID: 1, Path: "src/api.go", Kind: "file"},
				{NodeID: 2, Path: "src/store.go", Kind: "file"},
			},
		})
		return map[string]response{requests[0].ID: {OK: true, Result: payload}}, nil
	}

	graph, err := (Client{Run: run}).LoadGraphWithOptions(
		context.Background(),
		"snapshot",
		LoadOptions{SourcePrefixes: []string{"src"}},
	)
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("expected only the node-list request, got %d calls", calls)
	}
	if len(graph.Sources) != 2 || graph.Relationships != 0 {
		t.Fatalf("unexpected graph: %+v", graph)
	}
}

func TestClientLoadsOutgoingRelationshipsOnlyForSelectedPrefixes(t *testing.T) {
	calls := 0
	run := func(_ context.Context, _, _ string, requests []request) (map[string]response, error) {
		calls++
		if calls == 1 {
			payload, _ := json.Marshal(nodeListResult{
				Count:    2,
				Returned: 2,
				Nodes: []Node{
					{NodeID: 1, Path: "src/api/a.go", Kind: "function"},
					{NodeID: 2, Path: "src/storage/b.go", Kind: "function"},
				},
			})
			return map[string]response{requests[0].ID: {OK: true, Result: payload}}, nil
		}
		if len(requests) != 1 || requests[0].NodeID != 1 {
			t.Fatalf("expected only api neighbors, got %+v", requests)
		}
		payload, _ := json.Marshal(neighborResult{Node: Node{NodeID: 1}, Count: 0})
		return map[string]response{requests[0].ID: {OK: true, Result: payload}}, nil
	}

	_, err := (Client{Run: run}).LoadGraphWithOptions(
		context.Background(),
		"snapshot",
		LoadOptions{
			SourcePrefixes:   []string{"src"},
			OutgoingPrefixes: []string{"src/api"},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatalf("expected list and neighbor calls, got %d", calls)
	}
}

func TestClientPaginatesNodeLists(t *testing.T) {
	calls := 0
	run := func(_ context.Context, _, _ string, requests []request) (map[string]response, error) {
		calls++
		current := requests[0]
		if current.Op == "neighbors" {
			results := make(map[string]response, len(requests))
			for _, neighborRequest := range requests {
				payload, _ := json.Marshal(neighborResult{
					Node:  Node{NodeID: neighborRequest.NodeID},
					Count: 0,
				})
				results[neighborRequest.ID] = response{OK: true, Result: payload}
			}
			return results, nil
		}
		if current.Offset == 0 {
			next := 1
			payload, _ := json.Marshal(nodeListResult{
				Count:      2,
				Offset:     0,
				Returned:   1,
				Truncated:  true,
				NextOffset: &next,
				Nodes:      []Node{{NodeID: 1, Path: "api/a.go", Name: "A"}},
			})
			return map[string]response{current.ID: {OK: true, Result: payload}}, nil
		}
		payload, _ := json.Marshal(nodeListResult{
			Count:    2,
			Offset:   1,
			Returned: 1,
			Nodes:    []Node{{NodeID: 2, Path: "api/b.go", Name: "B"}},
		})
		return map[string]response{current.ID: {OK: true, Result: payload}}, nil
	}

	graph, err := (Client{Run: run}).LoadGraph(context.Background(), "snapshot", []string{"api"})
	if err != nil {
		t.Fatal(err)
	}
	if calls != 3 {
		t.Fatalf("expected two list pages and one neighbor batch, got %d calls", calls)
	}
	if len(graph.Sources) != 2 {
		t.Fatalf("expected two sources, got %d", len(graph.Sources))
	}
}

func TestClientFailsClosedOnTruncatedNodeList(t *testing.T) {
	run := func(_ context.Context, _, _ string, requests []request) (map[string]response, error) {
		payload, _ := json.Marshal(nodeListResult{Count: 10_001, Returned: 10_000, Truncated: true})
		return map[string]response{
			requests[0].ID: {OK: true, Result: payload},
		}, nil
	}
	if _, err := (Client{Run: run}).LoadGraph(context.Background(), "snapshot", []string{"api"}); err == nil {
		t.Fatal("expected truncation to fail")
	}
}
