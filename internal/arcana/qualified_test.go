package arcana

import (
	"context"
	"encoding/json"
	"testing"
)

func TestResolveQualifiedUsesPathAndName(t *testing.T) {
	run := func(_ context.Context, _, _ string, requests []request) (map[string]response, error) {
		current := requests[0]
		if current.Op != "resolve_symbol" || current.Path != "api/api.go" || current.Name != "CreateUser" {
			t.Fatalf("unexpected request: %+v", current)
		}
		payload, _ := json.Marshal(nodeListResult{
			Count:    1,
			Returned: 1,
			Nodes:    []Node{{NodeID: 1, Path: "api/api.go", Name: "CreateUser"}},
		})
		return map[string]response{current.ID: {OK: true, Result: payload}}, nil
	}
	nodes, err := (Client{Run: run}).ResolveQualified(context.Background(), "snapshot", "api/api.go::CreateUser")
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 || nodes[0].NodeID != 1 {
		t.Fatalf("unexpected nodes: %+v", nodes)
	}
}

func TestHasQualifiedRelationship(t *testing.T) {
	calls := 0
	run := func(_ context.Context, _, _ string, requests []request) (map[string]response, error) {
		calls++
		current := requests[0]
		if calls == 1 {
			payload, _ := json.Marshal(nodeListResult{
				Count:    1,
				Returned: 1,
				Nodes:    []Node{{NodeID: 1, Path: "api/api.go", Name: "CreateUser"}},
			})
			return map[string]response{current.ID: {OK: true, Result: payload}}, nil
		}
		if current.Op != "neighbors" || current.Relation != "calls" || current.Limit != maxResultLimit {
			t.Fatalf("unexpected neighbor request: %+v", current)
		}
		payload, _ := json.Marshal(neighborResult{
			Node:     Node{NodeID: 1},
			Count:    1,
			Returned: 1,
			Relationships: []Relationship{
				{Relation: "calls", Node: Node{NodeID: 2, Path: "storage/users.go", Name: "InsertUser"}},
			},
		})
		return map[string]response{current.ID: {OK: true, Result: payload}}, nil
	}

	found, err := (Client{Run: run}).HasQualifiedRelationship(
		context.Background(),
		"snapshot",
		"api/api.go::CreateUser",
		"calls",
		"storage/users.go::InsertUser",
	)
	if err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Fatal("expected relationship to be found")
	}
}

func TestResolveQualifiedRejectsUnqualifiedName(t *testing.T) {
	if _, err := (Client{}).ResolveQualified(context.Background(), "snapshot", "CreateUser"); err == nil {
		t.Fatal("expected invalid qualified name to fail")
	}
}
