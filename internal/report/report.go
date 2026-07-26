package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
	"github.com/Lokee86/pitlord/internal/policy"
)

func WriteJSON(writer io.Writer, value policy.Report) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func WriteText(writer io.Writer, value policy.Report) error {
	if len(value.Diagnostics) == 0 {
		if _, err := fmt.Fprintf(
			writer,
			"Pitlord: no architecture violations (%d rules, %d source nodes, %d relationships).\n",
			value.Summary.RulesChecked,
			value.Summary.SourceNodesScanned,
			value.Summary.RelationshipsScanned,
		); err != nil {
			return err
		}
	} else {
		for _, diagnostic := range value.Diagnostics {
			if err := WriteDiagnosticText(writer, diagnostic); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(
			writer,
			"Pitlord: %d errors, %d warnings (%d rules, %d source nodes, %d relationships).\n",
			value.Summary.Errors,
			value.Summary.Warnings,
			value.Summary.RulesChecked,
			value.Summary.SourceNodesScanned,
			value.Summary.RelationshipsScanned,
		); err != nil {
			return err
		}
	}
	if value.Summary.SuppressedEvidence > 0 {
		if _, err := fmt.Fprintf(
			writer,
			"Pitlord: suppressed %d baseline finding(s).\n",
			value.Summary.SuppressedEvidence,
		); err != nil {
			return err
		}
	}
	if value.Expectation != nil && !value.Expectation.Matched {
		if len(value.Expectation.MissingIDs) > 0 {
			if _, err := fmt.Fprintf(writer, "Missing expected diagnostics: %s\n", strings.Join(value.Expectation.MissingIDs, ", ")); err != nil {
				return err
			}
		}
		if len(value.Expectation.Unexpected) > 0 {
			if _, err := fmt.Fprintf(writer, "Unexpected diagnostics: %s\n", strings.Join(value.Expectation.Unexpected, ", ")); err != nil {
				return err
			}
		}
	}
	return nil
}

func WriteDiagnosticText(writer io.Writer, diagnostic policy.Diagnostic) error {
	if _, err := fmt.Fprintf(
		writer,
		"%s %s: %s\n",
		strings.ToUpper(diagnostic.Severity),
		diagnostic.RuleID,
		diagnostic.Message,
	); err != nil {
		return err
	}
	for _, evidence := range diagnostic.Evidence {
		if _, err := fmt.Fprintf(writer, "  %s\n", evidenceText(evidence)); err != nil {
			return err
		}
	}
	return nil
}

func evidenceText(evidence policy.Evidence) string {
	prefix := fmt.Sprintf("%s %s", location(evidence.Source), displayNode(evidence.Source))
	switch evidence.Issue {
	case "forbidden_dependency":
		if evidence.Target == nil {
			return prefix + " has malformed dependency evidence"
		}
		return fmt.Sprintf(
			"%s --%s--> %s",
			prefix,
			evidence.Relation,
			displayNode(*evidence.Target),
		)
	case "unowned":
		return prefix + " has no declared owner"
	case "multiple_owners":
		return fmt.Sprintf("%s matches multiple owners: %s", prefix, strings.Join(evidence.Areas, ", "))
	case "area_cycle":
		if evidence.Target == nil {
			return prefix + " has malformed area-cycle evidence"
		}
		return fmt.Sprintf(
			"%s [%s] --%s--> %s [%s] closes cycle among %s",
			prefix,
			evidence.SourceArea,
			evidence.Relation,
			displayNode(*evidence.Target),
			evidence.TargetArea,
			strings.Join(evidence.Areas, ", "),
		)
	case "area_dependency":
		if evidence.Target == nil {
			return prefix + " has malformed area-dependency evidence"
		}
		return fmt.Sprintf(
			"%s [%s] --%s--> %s [%s]",
			prefix,
			evidence.SourceArea,
			evidence.Relation,
			displayNode(*evidence.Target),
			evidence.TargetArea,
		)
	default:
		return prefix + " violates policy"
	}
}

func location(node arcana.Node) string {
	if node.Span == nil || node.Span.StartLine <= 0 {
		return node.Path
	}
	return fmt.Sprintf("%s:%d:%d", node.Span.Path, node.Span.StartLine, node.Span.StartColumn)
}

func displayNode(node arcana.Node) string {
	if node.Name == "" {
		return node.Path
	}
	return node.Name
}
