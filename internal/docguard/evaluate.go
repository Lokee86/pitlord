package docguard

import (
	"path/filepath"
	"sort"
	"strings"
)

func Evaluate(codeFiles []string, dataset Dataset, changes []Change, changedFrom string) Report {
	owners := buildOwners(dataset)
	codeSet := make(map[string]struct{}, len(codeFiles))
	for _, file := range codeFiles {
		codeSet[normalizePath(file)] = struct{}{}
	}

	findings := make([]Finding, 0)
	mapped := 0
	files := sortedKeys(codeSet)
	for _, file := range files {
		docs := owners[file]
		if len(docs) == 0 {
			findings = append(findings, Finding{Code: "missing_codemap", CodePath: file})
			continue
		}
		mapped++
	}

	changedDocs := changedDocumentation(changes, dataset)
	changedCode := selectChangedCode(changes, codeSet, owners)
	for _, file := range changedCode {
		docs := owners[file]
		if len(docs) == 0 {
			continue
		}
		if intersects(docs, changedDocs) {
			continue
		}
		findings = append(findings, Finding{Code: "documentation_not_changed", CodePath: file, Documents: append([]string(nil), docs...)})
	}

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].CodePath != findings[j].CodePath {
			return findings[i].CodePath < findings[j].CodePath
		}
		return findings[i].Code < findings[j].Code
	})
	return Report{
		Schema:          "pitlord.documentation-guard.v1",
		ChangedFrom:     changedFrom,
		CodeFiles:       len(files),
		MappedCodeFiles: mapped,
		ChangedCode:     len(changedCode),
		ChangedDocs:     len(changedDocs),
		Findings:        findings,
	}
}

func buildOwners(dataset Dataset) map[string][]string {
	sets := make(map[string]map[string]struct{})
	add := func(file, document string) {
		file, document = normalizePath(file), normalizePath(document)
		if file == "" || document == "" {
			return
		}
		if sets[file] == nil {
			sets[file] = make(map[string]struct{})
		}
		sets[file][document] = struct{}{}
	}
	for _, item := range dataset.Entries {
		if item.Entry.Kind != "directory" && item.Resolution.ResolvedPath != "" {
			add(item.Resolution.ResolvedPath, item.Entry.DocumentPath)
		}
		if item.Resolution.SemanticNode != nil {
			add(item.Resolution.SemanticNode.Path, item.Entry.DocumentPath)
		}
		for _, match := range item.Resolution.Matches {
			if !match.IsDir {
				add(match.Path, item.Entry.DocumentPath)
			}
		}
	}
	owners := make(map[string][]string, len(sets))
	for file, documents := range sets {
		owners[file] = sortedKeys(documents)
	}
	return owners
}

func changedDocumentation(changes []Change, dataset Dataset) map[string]struct{} {
	documents := make(map[string]struct{})
	for _, entry := range dataset.Entries {
		documents[normalizePath(entry.Entry.DocumentPath)] = struct{}{}
	}
	changed := make(map[string]struct{})
	for _, change := range changes {
		for _, candidate := range []string{change.Path, change.OldPath} {
			candidate = normalizePath(candidate)
			if _, ok := documents[candidate]; ok {
				changed[candidate] = struct{}{}
			}
		}
	}
	return changed
}

func selectChangedCode(changes []Change, codeSet map[string]struct{}, owners map[string][]string) []string {
	selected := make(map[string]struct{})
	for _, change := range changes {
		candidate := normalizePath(change.Path)
		if strings.HasPrefix(change.Status, "D") {
			candidate = normalizePath(change.OldPath)
			if candidate == "" {
				candidate = normalizePath(change.Path)
			}
		}
		if _, ok := codeSet[candidate]; ok || len(owners[candidate]) > 0 {
			selected[candidate] = struct{}{}
		}
	}
	return sortedKeys(selected)
}

func intersects(documents []string, changed map[string]struct{}) bool {
	for _, document := range documents {
		if _, ok := changed[document]; ok {
			return true
		}
	}
	return false
}

func normalizePath(value string) string {
	value = filepath.ToSlash(filepath.Clean(strings.TrimSpace(value)))
	value = strings.TrimPrefix(value, "./")
	if value == "." {
		return ""
	}
	return value
}

func sortedKeys[T any](values map[string]T) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		if key != "" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	return keys
}
