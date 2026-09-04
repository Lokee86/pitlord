package scan

import (
	"context"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

const AnalyzerUnobservedOutcome = "unobserved-outcome"

var unobservedOutcomeCapabilities = []SemanticCapability{
	CapabilityCalls,
	CapabilitySourceSpans,
	CapabilityOutcomeObligations,
}

type unobservedOutcomeAnalyzer struct{}

func (unobservedOutcomeAnalyzer) Metadata() AnalyzerMetadata {
	return AnalyzerMetadata{
		ID:            AnalyzerUnobservedOutcome,
		Category:      "outcome-handling",
		RequiresGraph: true,
		Capabilities:  unobservedOutcomeCapabilities,
	}
}

func (unobservedOutcomeAnalyzer) Analyze(_ context.Context, input AnalyzerContext) ([]Finding, error) {
	supportedPaths := semanticCapabilityPaths(input.Graph, unobservedOutcomeCapabilities)
	findings := make([]Finding, 0)
	for _, node := range input.Graph.Sources {
		if node.Kind != "protocol" || !strings.HasPrefix(node.Name, "outcome-operation:") {
			continue
		}
		if _, supported := supportedPaths[normalizedUnitPath(node.Path)]; !supported {
			continue
		}
		if outcomeIsConsumed(input.Graph, node.NodeID) {
			continue
		}
		findings = append(findings, unobservedOutcomeFinding(node))
	}
	return findings, nil
}

func outcomeIsConsumed(graph arcana.Graph, nodeID uint32) bool {
	for _, relationship := range graph.Outgoing[nodeID] {
		if relationship.Relation == "contains" && relationship.Node.Kind == "protocol" && relationship.Node.Name == "outcome-action:consume" {
			return true
		}
	}
	return false
}

func unobservedOutcomeFinding(node arcana.Node) Finding {
	parts := strings.SplitN(strings.TrimPrefix(node.Name, "outcome-operation:"), ":", 2)
	language := parts[0]
	obligation := "outcome"
	if len(parts) == 2 {
		obligation = parts[1]
	}
	scopeKey := node.Identity
	if scopeKey == "" {
		scopeKey = node.Path
		if node.Span != nil {
			scopeKey += ":" + spanIdentity(node.Span)
		}
	}
	return Finding{
		ID:                findingID(AnalyzerUnobservedOutcome, scopeKey),
		RuleID:            AnalyzerUnobservedOutcome,
		Detector:          AnalyzerUnobservedOutcome,
		Language:          language,
		Disposition:       DispositionAdvisory,
		Severity:          SeverityWarning,
		Scope:             Scope{Kind: "file", Path: normalizedUnitPath(node.Path)},
		Location:          sourceSpanFromArcana(node),
		Summary:           "Operation result requires observation but is discarded",
		Rationale:         "The language adapter proved that this operation creates an outcome obligation and found no semantic consumption of that outcome.",
		Evidence:          []Evidence{{Kind: "semantic-outcome-obligation", Message: obligation}},
		RequiredOutcome:   "The operation outcome must be consumed, transferred, awaited, handled, or explicitly discarded.",
		RecommendedAction: "Observe the operation outcome or make the intentional discard explicit.",
	}
}
