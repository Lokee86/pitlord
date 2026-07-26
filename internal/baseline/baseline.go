package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Lokee86/pitlord/internal/policy"
)

func Build(diagnostics []policy.Diagnostic) Document {
	byFingerprint := make(map[string]Entry)
	for _, diagnostic := range diagnostics {
		for _, evidence := range diagnostic.Evidence {
			fingerprint := Fingerprint(diagnostic.RuleID, evidence)
			entry := Entry{
				Fingerprint: fingerprint,
				RuleID:      diagnostic.RuleID,
				Issue:       evidence.Issue,
				Source:      evidence.Source.Path,
				Relation:    evidence.Relation,
				Areas:       append([]string(nil), evidence.Areas...),
				SourceArea:  evidence.SourceArea,
				TargetArea:  evidence.TargetArea,
			}
			if evidence.Target != nil {
				entry.Target = evidence.Target.Path
			}
			sort.Strings(entry.Areas)
			byFingerprint[fingerprint] = entry
		}
	}
	entries := make([]Entry, 0, len(byFingerprint))
	for _, entry := range byFingerprint {
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Fingerprint < entries[j].Fingerprint
	})
	return Document{Version: Version, Entries: entries}
}

func Load(path string) (Document, error) {
	file, err := os.Open(path)
	if err != nil {
		return Document{}, fmt.Errorf("open baseline: %w", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	var document Document
	if err := decoder.Decode(&document); err != nil {
		return Document{}, fmt.Errorf("decode baseline: %w", err)
	}
	if err := document.Validate(); err != nil {
		return Document{}, err
	}
	return document, nil
}

func Write(path string, document Document) error {
	if document.Version == 0 {
		document.Version = Version
	}
	if err := document.Validate(); err != nil {
		return err
	}
	data, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return fmt.Errorf("encode baseline: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write baseline: %w", err)
	}
	return nil
}

func (document Document) Validate() error {
	if document.Version != Version {
		return fmt.Errorf("unsupported baseline version %d", document.Version)
	}
	seen := make(map[string]struct{}, len(document.Entries))
	for index, entry := range document.Entries {
		if !strings.HasPrefix(entry.Fingerprint, "sha256:") {
			return fmt.Errorf("baseline entry %d has invalid fingerprint %q", index+1, entry.Fingerprint)
		}
		if entry.RuleID == "" || entry.Issue == "" {
			return fmt.Errorf("baseline entry %d is missing rule_id or issue", index+1)
		}
		if _, exists := seen[entry.Fingerprint]; exists {
			return fmt.Errorf("duplicate baseline fingerprint %q", entry.Fingerprint)
		}
		seen[entry.Fingerprint] = struct{}{}
	}
	return nil
}

func Apply(document Document, diagnostics []policy.Diagnostic) ([]policy.Diagnostic, int) {
	known := make(map[string]struct{}, len(document.Entries))
	for _, entry := range document.Entries {
		known[entry.Fingerprint] = struct{}{}
	}
	filtered := make([]policy.Diagnostic, 0, len(diagnostics))
	suppressed := 0
	for _, diagnostic := range diagnostics {
		remaining := make([]policy.Evidence, 0, len(diagnostic.Evidence))
		for _, evidence := range diagnostic.Evidence {
			if _, exists := known[Fingerprint(diagnostic.RuleID, evidence)]; exists {
				suppressed++
				continue
			}
			remaining = append(remaining, evidence)
		}
		if len(remaining) == 0 {
			continue
		}
		diagnostic.Evidence = remaining
		filtered = append(filtered, diagnostic)
	}
	return filtered, suppressed
}
