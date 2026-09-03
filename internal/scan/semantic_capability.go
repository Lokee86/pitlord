package scan

import (
	"fmt"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

type SemanticCapability string

const (
	CapabilityControlFlow   SemanticCapability = "control-flow"
	CapabilityErrorHandling SemanticCapability = "error-handling"
	CapabilityCalls         SemanticCapability = "calls"
	CapabilitySourceSpans   SemanticCapability = "source-spans"
)

const semanticCapabilityPrefix = "semantic-capabilities:"

func validateAnalyzerCapabilities(metadata AnalyzerMetadata, graph arcana.Graph) error {
	if len(metadata.Capabilities) == 0 {
		return nil
	}
	for _, node := range graph.Sources {
		if node.Kind != "protocol" || !strings.HasPrefix(node.Name, semanticCapabilityPrefix) {
			continue
		}
		if capabilityNodeSatisfies(node.Name, metadata.Capabilities) {
			return nil
		}
	}
	return fmt.Errorf("scan analyzer %q requires semantic capabilities: %s", metadata.ID, joinCapabilities(metadata.Capabilities))
}

func capabilityNodeSatisfies(name string, required []SemanticCapability) bool {
	parts := strings.SplitN(name, ":", 3)
	if len(parts) != 3 {
		return false
	}
	available := make(map[SemanticCapability]struct{})
	for _, value := range strings.Split(parts[2], ",") {
		available[SemanticCapability(strings.TrimSpace(value))] = struct{}{}
	}
	for _, capability := range required {
		if _, ok := available[capability]; !ok {
			return false
		}
	}
	return true
}

func joinCapabilities(capabilities []SemanticCapability) string {
	values := make([]string, 0, len(capabilities))
	for _, capability := range capabilities {
		values = append(values, string(capability))
	}
	return strings.Join(values, ", ")
}
