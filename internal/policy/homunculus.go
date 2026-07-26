package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"
)

type homunculusManifest struct {
	Version      int                    `json:"version"`
	Architecture homunculusArchitecture `json:"architecture"`
	Mutations    []homunculusMutation   `json:"mutations"`
}

type homunculusArchitecture struct {
	Areas     []string `json:"areas"`
	Forbidden []string `json:"forbidden"`
}

type homunculusMutation struct {
	ExpectedDiagnostics []string `json:"expected_diagnostics"`
}

func FromHomunculus(path string) (Document, []string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Document{}, nil, fmt.Errorf("read Homunculus manifest: %w", err)
	}
	var manifest homunculusManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return Document{}, nil, fmt.Errorf("decode Homunculus manifest: %w", err)
	}
	if manifest.Version == 0 {
		return Document{}, nil, fmt.Errorf("Homunculus manifest has no version")
	}
	areas := make(map[string]struct{}, len(manifest.Architecture.Areas))
	document := Document{Version: Version}
	for _, area := range manifest.Architecture.Areas {
		area = normalizePath(area)
		if area == "" {
			continue
		}
		if _, exists := areas[area]; exists {
			continue
		}
		areas[area] = struct{}{}
		document.Areas = append(document.Areas, Area{ID: area, Paths: []string{area}})
	}
	for _, forbidden := range manifest.Architecture.Forbidden {
		from, to, err := parseForbidden(forbidden)
		if err != nil {
			return Document{}, nil, err
		}
		if len(areas) > 0 {
			if _, exists := areas[from]; !exists {
				return Document{}, nil, fmt.Errorf("forbidden dependency %q references unknown area %q", forbidden, from)
			}
			if _, exists := areas[to]; !exists {
				return Document{}, nil, fmt.Errorf("forbidden dependency %q references unknown area %q", forbidden, to)
			}
		}
		document.Rules = append(document.Rules, Rule{
			ID:        slug(from) + "-must-not-access-" + slug(to),
			Type:      RuleForbidDependency,
			Message:   fmt.Sprintf("%s must not access %s", from, to),
			Severity:  "error",
			FromAreas: []string{from},
			ToAreas:   []string{to},
			Relations: []string{"calls", "imports"},
		})
	}
	if err := document.NormalizeAndValidate(); err != nil {
		return Document{}, nil, fmt.Errorf("build policy from Homunculus manifest: %w", err)
	}

	expectedSet := make(map[string]struct{})
	for _, mutation := range manifest.Mutations {
		for _, id := range mutation.ExpectedDiagnostics {
			id = strings.TrimSpace(id)
			if id != "" {
				expectedSet[id] = struct{}{}
			}
		}
	}
	expected := make([]string, 0, len(expectedSet))
	for id := range expectedSet {
		expected = append(expected, id)
	}
	sort.Strings(expected)
	return document, expected, nil
}

func parseForbidden(value string) (string, string, error) {
	parts := strings.Split(value, "->")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid Homunculus forbidden dependency %q; expected 'from -> to'", value)
	}
	from := normalizePath(parts[0])
	to := normalizePath(parts[1])
	if from == "" || to == "" {
		return "", "", fmt.Errorf("invalid Homunculus forbidden dependency %q", value)
	}
	return from, to, nil
}

func slug(value string) string {
	var builder strings.Builder
	lastDash := false
	for _, current := range strings.ToLower(value) {
		if unicode.IsLetter(current) || unicode.IsDigit(current) {
			builder.WriteRune(current)
			lastDash = false
			continue
		}
		if builder.Len() > 0 && !lastDash {
			builder.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(builder.String(), "-")
}
