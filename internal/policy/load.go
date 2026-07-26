package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Load(path string) (Document, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return Document{}, fmt.Errorf("resolve policy path: %w", err)
	}
	loader := policyLoader{
		loaded:   make(map[string]struct{}),
		visiting: make(map[string]int),
	}
	document, err := loader.load(filepath.Clean(absolute), nil)
	if err != nil {
		return Document{}, err
	}
	document.Includes = nil
	if err := document.NormalizeAndValidate(); err != nil {
		return Document{}, err
	}
	return document, nil
}

type policyLoader struct {
	loaded   map[string]struct{}
	visiting map[string]int
	stack    []string
}

func (loader *policyLoader) load(path string, chain []string) (Document, error) {
	path = filepath.Clean(path)
	if _, exists := loader.loaded[path]; exists {
		return Document{}, nil
	}
	if index, exists := loader.visiting[path]; exists {
		cycle := append(append([]string(nil), loader.stack[index:]...), path)
		return Document{}, fmt.Errorf("policy include cycle: %s", strings.Join(cycle, " -> "))
	}

	raw, err := decodeDocument(path)
	if err != nil {
		if len(chain) > 0 {
			return Document{}, fmt.Errorf("include %s: %w", strings.Join(chain, " -> "), err)
		}
		return Document{}, err
	}
	if raw.Version == 0 {
		raw.Version = Version
	}
	if raw.Version != Version {
		return Document{}, fmt.Errorf("policy %q has unsupported version %d", path, raw.Version)
	}

	loader.visiting[path] = len(loader.stack)
	loader.stack = append(loader.stack, path)
	merged := Document{Version: Version}
	baseDirectory := filepath.Dir(path)
	for _, include := range normalizeIncludes(raw.Includes) {
		includePath := include
		if !filepath.IsAbs(includePath) {
			includePath = filepath.Join(baseDirectory, includePath)
		}
		included, err := loader.load(
			filepath.Clean(includePath),
			append(append([]string(nil), chain...), path),
		)
		if err != nil {
			return Document{}, err
		}
		merged.Areas = append(merged.Areas, included.Areas...)
		merged.Rules = append(merged.Rules, included.Rules...)
	}
	merged.Areas = append(merged.Areas, raw.Areas...)
	merged.Rules = append(merged.Rules, raw.Rules...)

	loader.stack = loader.stack[:len(loader.stack)-1]
	delete(loader.visiting, path)
	loader.loaded[path] = struct{}{}
	return merged, nil
}

func decodeDocument(path string) (Document, error) {
	file, err := os.Open(path)
	if err != nil {
		return Document{}, fmt.Errorf("open policy %q: %w", path, err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	var document Document
	if err := decoder.Decode(&document); err != nil {
		return Document{}, fmt.Errorf("decode policy %q: %w", path, err)
	}
	return document, nil
}

func Write(path string, document Document) error {
	if err := document.NormalizeAndValidate(); err != nil {
		return err
	}
	data, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return fmt.Errorf("encode policy: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write policy: %w", err)
	}
	return nil
}

func normalizeIncludes(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
