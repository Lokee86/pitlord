package snapshotdiff

import (
	"sort"

	"github.com/Lokee86/pitlord/internal/baseline"
	"github.com/Lokee86/pitlord/internal/policy"
)

func Compare(
	policySource string,
	beforeSnapshot string,
	afterSnapshot string,
	before []policy.Diagnostic,
	after []policy.Diagnostic,
) Result {
	beforeFindings := flatten(before)
	afterFindings := flatten(after)
	beforeByFingerprint := index(beforeFindings)
	afterByFingerprint := index(afterFindings)

	introduced := make([]Finding, 0)
	resolved := make([]Finding, 0)
	persistent := make([]Finding, 0)
	for fingerprint, finding := range afterByFingerprint {
		if _, exists := beforeByFingerprint[fingerprint]; exists {
			persistent = append(persistent, finding)
		} else {
			introduced = append(introduced, finding)
		}
	}
	for fingerprint, finding := range beforeByFingerprint {
		if _, exists := afterByFingerprint[fingerprint]; !exists {
			resolved = append(resolved, finding)
		}
	}
	sortFindings(introduced)
	sortFindings(resolved)
	sortFindings(persistent)
	return Result{
		Schema:         "pitlord.snapshot-diff.v1",
		PolicySource:   policySource,
		BeforeSnapshot: beforeSnapshot,
		AfterSnapshot:  afterSnapshot,
		Introduced:     introduced,
		Resolved:       resolved,
		Persistent:     persistent,
		Summary: Summary{
			BeforeFindings: len(beforeFindings),
			AfterFindings:  len(afterFindings),
			Introduced:     len(introduced),
			Resolved:       len(resolved),
			Persistent:     len(persistent),
		},
	}
}

func Diagnostics(findings []Finding) []policy.Diagnostic {
	byRule := make(map[string]*policy.Diagnostic)
	for _, finding := range findings {
		diagnostic, exists := byRule[finding.RuleID]
		if !exists {
			diagnostic = &policy.Diagnostic{
				RuleID:   finding.RuleID,
				Message:  finding.Message,
				Severity: finding.Severity,
			}
			byRule[finding.RuleID] = diagnostic
		}
		diagnostic.Evidence = append(diagnostic.Evidence, finding.Evidence)
	}
	ids := make([]string, 0, len(byRule))
	for id := range byRule {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	diagnostics := make([]policy.Diagnostic, 0, len(ids))
	for _, id := range ids {
		diagnostics = append(diagnostics, *byRule[id])
	}
	return diagnostics
}

func flatten(diagnostics []policy.Diagnostic) []Finding {
	findings := make([]Finding, 0)
	for _, diagnostic := range diagnostics {
		for _, evidence := range diagnostic.Evidence {
			findings = append(findings, Finding{
				Fingerprint: baseline.Fingerprint(diagnostic.RuleID, evidence),
				RuleID:      diagnostic.RuleID,
				Message:     diagnostic.Message,
				Severity:    diagnostic.Severity,
				Evidence:    evidence,
			})
		}
	}
	sortFindings(findings)
	return findings
}

func index(findings []Finding) map[string]Finding {
	result := make(map[string]Finding, len(findings))
	for _, finding := range findings {
		result[finding.Fingerprint] = finding
	}
	return result
}

func sortFindings(findings []Finding) {
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].RuleID != findings[j].RuleID {
			return findings[i].RuleID < findings[j].RuleID
		}
		return findings[i].Fingerprint < findings[j].Fingerprint
	})
}
