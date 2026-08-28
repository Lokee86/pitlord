package calibration

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func LoadReference(filename string) (Reference, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Reference{}, fmt.Errorf("open calibration reference: %w", err)
	}
	defer file.Close()
	return DecodeReference(file)
}

func DecodeReference(reader io.Reader) (Reference, error) {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	var reference Reference
	if err := decoder.Decode(&reference); err != nil {
		return Reference{}, fmt.Errorf("decode calibration reference: %w", err)
	}
	if err := validateReference(reference); err != nil {
		return Reference{}, err
	}
	return reference, nil
}

func validateReference(reference Reference) error {
	if reference.Schema != Schema {
		return fmt.Errorf("unsupported calibration schema %q", reference.Schema)
	}
	if strings.TrimSpace(reference.Corpus) == "" {
		return fmt.Errorf("calibration corpus is required")
	}
	if strings.TrimSpace(reference.SourceRevision) == "" {
		return fmt.Errorf("calibration source_revision is required")
	}
	if strings.TrimSpace(reference.Detector) == "" {
		return fmt.Errorf("calibration detector is required")
	}
	if len(reference.Expectations) == 0 {
		return fmt.Errorf("calibration expectations are required")
	}
	seen := map[string]struct{}{}
	for index, expectation := range reference.Expectations {
		if err := validateExpectation(expectation); err != nil {
			return fmt.Errorf("expectation %d: %w", index, err)
		}
		if _, ok := seen[expectation.ID]; ok {
			return fmt.Errorf("duplicate calibration expectation id %q", expectation.ID)
		}
		seen[expectation.ID] = struct{}{}
	}
	return nil
}

func validateExpectation(expectation Expectation) error {
	if strings.TrimSpace(expectation.ID) == "" {
		return fmt.Errorf("id is required")
	}
	if (expectation.Path == "") == (expectation.PathPrefix == "") {
		return fmt.Errorf("exactly one of path or path_prefix is required")
	}
	if expectation.Path != "" {
		if err := validateRelativePath(expectation.Path); err != nil {
			return fmt.Errorf("path: %w", err)
		}
	}
	if expectation.PathPrefix != "" {
		if err := validateRelativePath(expectation.PathPrefix); err != nil {
			return fmt.Errorf("path_prefix: %w", err)
		}
	}
	if strings.TrimSpace(expectation.Class) == "" {
		return fmt.Errorf("class is required")
	}
	switch expectation.Finding {
	case FindingRequired, FindingAbsent, FindingAllowed:
	default:
		return fmt.Errorf("unsupported finding expectation %q", expectation.Finding)
	}
	if expectation.Finding == FindingAbsent && (expectation.MinSeverity != "" || expectation.MaxSeverity != "") {
		return fmt.Errorf("absent finding expectation cannot constrain severity")
	}
	if !validSeverity(string(expectation.MinSeverity)) || !validSeverity(string(expectation.MaxSeverity)) {
		return fmt.Errorf("invalid severity constraint")
	}
	return nil
}

func validateRelativePath(value string) error {
	value = strings.ReplaceAll(strings.TrimSpace(value), "\\", "/")
	cleaned := path.Clean(value)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") || strings.HasPrefix(cleaned, "/") {
		return fmt.Errorf("must be a repository-relative path")
	}
	return nil
}

func validSeverity(value string) bool {
	switch value {
	case "", "info", "warning", "high", "critical":
		return true
	default:
		return false
	}
}
