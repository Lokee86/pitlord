package report

import (
	"encoding/json"
	"io"
	"sort"

	"github.com/Lokee86/pitlord/internal/arcana"
	"github.com/Lokee86/pitlord/internal/baseline"
	"github.com/Lokee86/pitlord/internal/policy"
)

const sarifSchema = "https://json.schemastore.org/sarif-2.1.0.json"

type sarifLog struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool       sarifTool      `json:"tool"`
	Results    []sarifResult  `json:"results"`
	Properties map[string]any `json:"properties,omitempty"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string      `json:"name"`
	Version        string      `json:"version"`
	InformationURI string      `json:"informationUri,omitempty"`
	Rules          []sarifRule `json:"rules,omitempty"`
}

type sarifRule struct {
	ID               string       `json:"id"`
	ShortDescription sarifMessage `json:"shortDescription"`
	DefaultConfig    sarifConfig  `json:"defaultConfiguration"`
}

type sarifConfig struct {
	Level string `json:"level"`
}

type sarifResult struct {
	RuleID              string            `json:"ruleId"`
	Level               string            `json:"level"`
	Message             sarifMessage      `json:"message"`
	Locations           []sarifLocation   `json:"locations,omitempty"`
	RelatedLocations    []sarifLocation   `json:"relatedLocations,omitempty"`
	PartialFingerprints map[string]string `json:"partialFingerprints"`
	Properties          map[string]any    `json:"properties,omitempty"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

type sarifLocation struct {
	ID               int                   `json:"id,omitempty"`
	Message          *sarifMessage         `json:"message,omitempty"`
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           *sarifRegion          `json:"region,omitempty"`
}

type sarifArtifactLocation struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine   int `json:"startLine,omitempty"`
	StartColumn int `json:"startColumn,omitempty"`
	EndLine     int `json:"endLine,omitempty"`
	EndColumn   int `json:"endColumn,omitempty"`
}

func WriteSARIF(writer io.Writer, value policy.Report, toolVersion string) error {
	rules := sarifRules(value.Diagnostics)
	results := make([]sarifResult, 0)
	for _, diagnostic := range value.Diagnostics {
		for _, evidence := range diagnostic.Evidence {
			result := sarifResult{
				RuleID:  diagnostic.RuleID,
				Level:   sarifLevel(diagnostic.Severity),
				Message: sarifMessage{Text: evidenceText(evidence)},
				Locations: []sarifLocation{
					locationForNode(evidence.Source, 0, "policy source"),
				},
				PartialFingerprints: map[string]string{
					"pitlordEvidenceFingerprint/v1": baseline.Fingerprint(diagnostic.RuleID, evidence),
				},
				Properties: map[string]any{"issue": evidence.Issue},
			}
			if evidence.Relation != "" {
				result.Properties["relation"] = evidence.Relation
			}
			if len(evidence.Areas) > 0 {
				result.Properties["areas"] = evidence.Areas
			}
			if evidence.SourceArea != "" {
				result.Properties["sourceArea"] = evidence.SourceArea
			}
			if evidence.TargetArea != "" {
				result.Properties["targetArea"] = evidence.TargetArea
			}
			if evidence.Target != nil {
				result.RelatedLocations = []sarifLocation{
					locationForNode(*evidence.Target, 1, "dependency target"),
				}
			}
			results = append(results, result)
		}
	}
	log := sarifLog{
		Version: "2.1.0",
		Schema:  sarifSchema,
		Runs: []sarifRun{
			{
				Tool: sarifTool{Driver: sarifDriver{
					Name:    "Pitlord",
					Version: toolVersion,
					Rules:   rules,
				}},
				Results: results,
				Properties: map[string]any{
					"snapshot":           value.Snapshot,
					"policySource":       value.PolicySource,
					"suppressedEvidence": value.Summary.SuppressedEvidence,
				},
			},
		},
	}
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(log)
}

func sarifRules(diagnostics []policy.Diagnostic) []sarifRule {
	byID := make(map[string]sarifRule)
	for _, diagnostic := range diagnostics {
		byID[diagnostic.RuleID] = sarifRule{
			ID:               diagnostic.RuleID,
			ShortDescription: sarifMessage{Text: diagnostic.Message},
			DefaultConfig:    sarifConfig{Level: sarifLevel(diagnostic.Severity)},
		}
	}
	ids := make([]string, 0, len(byID))
	for id := range byID {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	rules := make([]sarifRule, 0, len(ids))
	for _, id := range ids {
		rules = append(rules, byID[id])
	}
	return rules
}

func locationForNode(node arcana.Node, id int, message string) sarifLocation {
	location := sarifLocation{
		ID:      id,
		Message: &sarifMessage{Text: message},
		PhysicalLocation: sarifPhysicalLocation{
			ArtifactLocation: sarifArtifactLocation{URI: node.Path},
		},
	}
	if node.Span != nil && node.Span.StartLine > 0 {
		location.PhysicalLocation.ArtifactLocation.URI = node.Span.Path
		location.PhysicalLocation.Region = &sarifRegion{
			StartLine:   node.Span.StartLine,
			StartColumn: node.Span.StartColumn,
			EndLine:     node.Span.EndLine,
			EndColumn:   node.Span.EndColumn,
		}
	}
	return location
}

func sarifLevel(severity string) string {
	if severity == "warning" {
		return "warning"
	}
	return "error"
}
