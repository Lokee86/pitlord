package arcana

import (
	"encoding/json"
	"testing"
)

func TestRequestEncodesZeroNodeID(t *testing.T) {
	encoded, err := json.Marshal(request{
		ID:     "neighbors-0",
		Op:     "neighbors",
		NodeID: 0,
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(encoded, &payload); err != nil {
		t.Fatalf("unmarshal request: %v", err)
	}

	value, exists := payload["node_id"]
	if !exists {
		t.Fatal("zero node_id was omitted from Arcana request")
	}
	if value != float64(0) {
		t.Fatalf("node_id = %v, want 0", value)
	}
}
