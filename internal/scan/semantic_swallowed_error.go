package scan

import (
	"context"
	"strconv"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

const AnalyzerSwallowedError = "swallowed-error"

var swallowedErrorCapabilities = []SemanticCapability{
	CapabilityControlFlow,
	CapabilityErrorHandling,
	CapabilityCalls,
	CapabilitySourceSpans,
}

type swallowedErrorAnalyzer struct{}

func (swallowedErrorAnalyzer) Metadata() AnalyzerMetadata {
	return AnalyzerMetadata{
		ID:            AnalyzerSwallowedError,
		Category:      "error-handling",
		RequiresGraph: true,
		Capabilities:  swallowedErrorCapabilities,
	}
}

func (swallowedErrorAnalyzer) Analyze(_ context.Context, input AnalyzerContext) ([]Finding, error) {
	supportedPaths := semanticCapabilityPaths(input.Graph, swallowedErrorCapabilities)
	findings := make([]Finding, 0)
	for _, node := range input.Graph.Sources {
		if node.Kind != "protocol" || !strings.HasPrefix(node.Name, "error-handler:") {
			continue
		}
		if _, supported := supportedPaths[normalizedUnitPath(node.Path)]; !supported {
			continue
		}
		if handlerHasAction(input.Graph, node.NodeID) {
			continue
		}
		findings = append(findings, swallowedErrorFinding(node))
	}
	return findings, nil
}

func semanticCapabilityPaths(graph arcana.Graph, required []SemanticCapability) map[string]struct{} {
	paths := make(map[string]struct{})
	for _, node := range graph.Sources {
		if node.Kind == "protocol" && strings.HasPrefix(node.Name, semanticCapabilityPrefix) && capabilityNodeSatisfies(node.Name, required) {
			paths[normalizedUnitPath(node.Path)] = struct{}{}
		}
	}
	return paths
}

func handlerHasAction(graph arcana.Graph, nodeID uint32) bool {
	for _, relationship := range graph.Outgoing[nodeID] {
		if relationship.Relation != "contains" || relationship.Node.Kind != "protocol" {
			continue
		}
		switch relationship.Node.Name {
		case "error-action:propagate", "error-action:record", "error-action:recover":
			return true
		}
	}
	return false
}

func swallowedErrorFinding(node arcana.Node) Finding {
	language := strings.TrimPrefix(node.Name, "error-handler:")
	location := sourceSpanFromArcana(node)
	scopeKey := node.Identity
	if scopeKey == "" {
		scopeKey = node.Path
		if node.Span != nil {
			scopeKey += ":" + spanIdentity(node.Span)
		}
	}
	return Finding{
		ID:                findingID(AnalyzerSwallowedError, scopeKey),
		RuleID:            AnalyzerSwallowedError,
		Detector:          AnalyzerSwallowedError,
		Language:          language,
		Disposition:       DispositionAdvisory,
		Severity:          SeverityWarning,
		Scope:             Scope{Kind: "file", Path: normalizedUnitPath(node.Path)},
		Location:          location,
		Summary:           "Error handler does not propagate, record, or explicitly recover from the error",
		Rationale:         "The language adapter identified an error-handling branch without a semantic error action.",
		Evidence:          []Evidence{{Kind: "semantic-error-handler", Message: language}},
		RequiredOutcome:   "The handler must propagate the error, record it, or perform an explicit recovery action.",
		RecommendedAction: "Handle the error explicitly instead of silently discarding it.",
	}
}

func sourceSpanFromArcana(node arcana.Node) *SourceSpan {
	if node.Span == nil {
		return nil
	}
	return &SourceSpan{
		Path:        normalizedUnitPath(node.Span.Path),
		StartLine:   node.Span.StartLine,
		StartColumn: node.Span.StartColumn,
		EndLine:     node.Span.EndLine,
		EndColumn:   node.Span.EndColumn,
	}
}

func spanIdentity(span *arcana.Span) string {
	return strings.Join([]string{
		strconv.Itoa(span.StartLine), strconv.Itoa(span.StartColumn), strconv.Itoa(span.EndLine), strconv.Itoa(span.EndColumn),
	}, ":")
}
