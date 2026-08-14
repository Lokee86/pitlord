package arcana

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	if os.Getenv("PITLORD_ARCANA_HELPER") == "1" {
		runArcanaProtocolHelper()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func runArcanaProtocolHelper() {
	switch os.Getenv("PITLORD_ARCANA_HELPER_MODE") {
	case "hang":
		time.Sleep(30 * time.Second)
	case "bad-protocol":
		emitHelperResponses("arcana.query.v999", -1)
	case "incomplete":
		emitHelperResponses(protocolID, 1)
	}
}

func emitHelperResponses(protocol string, maximum int) {
	scanner := bufio.NewScanner(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)
	emitted := 0
	for scanner.Scan() {
		if maximum >= 0 && emitted >= maximum {
			return
		}
		var current request
		if err := json.Unmarshal(scanner.Bytes(), &current); err != nil {
			return
		}
		_ = encoder.Encode(response{
			Protocol: protocol,
			ID:       current.ID,
			OK:       true,
			Result:   json.RawMessage(`{}`),
		})
		emitted++
	}
}

func TestRequestEncodesZeroNodeID(t *testing.T) {
	encoded, err := json.Marshal(request{ID: "neighbors-0", Op: "neighbors", NodeID: 0})
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err := json.Unmarshal(encoded, &payload); err != nil {
		t.Fatal(err)
	}
	value, exists := payload["node_id"]
	if !exists {
		t.Fatal("zero node_id was omitted from the Arcana request")
	}
	if value != float64(0) {
		t.Fatalf("node_id = %#v, want 0", value)
	}
}

func TestRunProtocolHonorsContextCancellation(t *testing.T) {
	t.Setenv("PITLORD_ARCANA_HELPER", "1")
	t.Setenv("PITLORD_ARCANA_HELPER_MODE", "hang")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := runProtocol(ctx, os.Args[0], "snapshot", []request{{ID: "request-1", Op: "list_nodes"}})
	if err == nil {
		t.Fatal("expected cancelled Arcana process to fail")
	}
	if ctx.Err() != context.DeadlineExceeded {
		t.Fatalf("expected deadline exceeded, got %v", ctx.Err())
	}
}

func TestRunProtocolRejectsIncompatibleProtocol(t *testing.T) {
	t.Setenv("PITLORD_ARCANA_HELPER", "1")
	t.Setenv("PITLORD_ARCANA_HELPER_MODE", "bad-protocol")

	_, err := runProtocol(context.Background(), os.Args[0], "snapshot", []request{{ID: "request-1", Op: "list_nodes"}})
	if err == nil {
		t.Fatal("expected incompatible protocol to fail")
	}
}

func TestRunProtocolRejectsIncompleteResponseStream(t *testing.T) {
	t.Setenv("PITLORD_ARCANA_HELPER", "1")
	t.Setenv("PITLORD_ARCANA_HELPER_MODE", "incomplete")

	_, err := runProtocol(context.Background(), os.Args[0], "snapshot", []request{
		{ID: "request-1", Op: "list_nodes"},
		{ID: "request-2", Op: "neighbors"},
	})
	if err == nil {
		t.Fatal("expected incomplete response stream to fail")
	}
}
