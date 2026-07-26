package policy

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

type repositoryEntry struct {
	path string
	dir  bool
}

func validateRepoGlob(value string) error {
	if value == "" || filepath.IsAbs(filepath.FromSlash(value)) || strings.HasPrefix(value, "../") || value == ".." {
		return fmt.Errorf("invalid repository glob %q", value)
	}
	if _, err := path.Match(value, ""); err != nil {
		return fmt.Errorf("invalid repository glob %q: %w", value, err)
	}
	return nil
}

func evaluateContentRule(rule Rule, repositoryRoot string, entries []repositoryEntry) *Diagnostic {
	matcher := contentMatcher(rule)
	evidence := make([]Evidence, 0)
	for _, entry := range entries {
		if entry.dir || !matchesAnyGlobOrParent(entry.path, rule.IncludePaths) || matchesAnyGlobOrParent(entry.path, rule.ExcludePaths) {
			continue
		}
		data, err := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(entry.path)))
		if err != nil {
			continue
		}
		for lineNumber, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSuffix(line, "\r")
			if !matcher(line) {
				continue
			}
			evidence = append(evidence, Evidence{
				Issue:  "forbidden_content",
				Source: arcanaNodeForPath(entry.path, lineNumber+1),
			})
		}
	}
	return diagnosticFromEvidence(rule, evidence)
}

func evaluatePathRule(rule Rule, entries []repositoryEntry) *Diagnostic {
	evidence := make([]Evidence, 0)
	for _, entry := range entries {
		if !matchesGlob(entry.path, rule.Path) {
			continue
		}
		if rule.Type == RuleRequirePath {
			return nil
		}
		evidence = append(evidence, Evidence{
			Issue:  "forbidden_path",
			Source: arcanaNodeForPath(entry.path, 0),
		})
	}
	if rule.Type == RuleRequirePath {
		evidence = append(evidence, Evidence{
			Issue:  "missing_path",
			Source: arcanaNodeForPath(rule.Path, 0),
		})
	}
	return diagnosticFromEvidence(rule, evidence)
}

func repositoryEntriesForDocument(root string, document Document) []repositoryEntry {
	if root == "" {
		root = "."
	}
	searchRoots := repositorySearchRoots(document)
	seen := make(map[string]repositoryEntry)
	for _, searchRoot := range searchRoots {
		absolute := root
		if searchRoot != "." {
			absolute = filepath.Join(root, filepath.FromSlash(searchRoot))
		}
		info, err := os.Stat(absolute)
		if err != nil {
			continue
		}
		if !info.IsDir() {
			relative, err := filepath.Rel(root, absolute)
			if err == nil {
				relative = filepath.ToSlash(relative)
				seen[relative] = repositoryEntry{path: relative}
			}
			continue
		}
		_ = filepath.WalkDir(absolute, func(current string, directoryEntry os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			relative, err := filepath.Rel(root, current)
			if err != nil || relative == "." {
				return nil
			}
			relative = filepath.ToSlash(relative)
			if directoryEntry.IsDir() && isIgnoredRepositoryDirectory(directoryEntry.Name()) {
				return filepath.SkipDir
			}
			seen[relative] = repositoryEntry{path: relative, dir: directoryEntry.IsDir()}
			return nil
		})
	}
	entries := make([]repositoryEntry, 0, len(seen))
	for _, entry := range seen {
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].path < entries[j].path })
	return entries
}

func repositorySearchRoots(document Document) []string {
	roots := make(map[string]struct{})
	for _, rule := range document.Rules {
		switch rule.Type {
		case RuleForbidContent:
			for _, pattern := range rule.IncludePaths {
				roots[repositorySearchRoot(pattern)] = struct{}{}
			}
		case RuleRequirePath, RuleForbidPath:
			roots[repositorySearchRoot(rule.Path)] = struct{}{}
		}
	}
	result := make([]string, 0, len(roots))
	for candidate := range roots {
		covered := false
		for parent := range roots {
			if candidate == parent {
				continue
			}
			if parent == "." {
				covered = true
				break
			}
			if strings.HasPrefix(candidate+"/", strings.TrimSuffix(parent, "/")+"/") {
				covered = true
				break
			}
		}
		if !covered {
			result = append(result, candidate)
		}
	}
	if len(result) == 0 {
		return []string{"."}
	}
	sort.Strings(result)
	return result
}

func repositorySearchRoot(pattern string) string {
	wildcard := strings.IndexAny(pattern, "*?[")
	if wildcard < 0 {
		return pattern
	}
	prefix := pattern[:wildcard]
	if strings.HasSuffix(prefix, "/") {
		prefix = strings.TrimSuffix(prefix, "/")
		if prefix == "" {
			return "."
		}
		return prefix
	}
	root := path.Dir(prefix)
	if root == "" {
		return "."
	}
	return root
}

func isIgnoredRepositoryDirectory(name string) bool {
	switch name {
	case ".git", ".worktrees", ".workingtrees", ".godot", "node_modules", "__pycache__",
		".arcana", ".lexicon", ".grimoire", ".pitlord", ".ddocs", ".homunculus",
		".cantrip", ".incubus", ".ritual", ".warlock", ".obsidian":
		return true
	default:
		return false
	}
}

func contentMatcher(rule Rule) func(string) bool {
	pattern := rule.Pattern
	patternType := rule.PatternType
	if rule.Literal != "" {
		pattern, patternType = rule.Literal, "literal"
	} else if rule.Regex != "" {
		pattern, patternType = rule.Regex, "regex"
	}
	if patternType == "regex" {
		compiled := regexp.MustCompile(pattern)
		return compiled.MatchString
	}
	return func(line string) bool { return strings.Contains(line, pattern) }
}

func matchesAnyGlobOrParent(candidate string, globs []string) bool {
	if matchesAnyGlob(candidate, globs) {
		return true
	}
	for parent := path.Dir(candidate); parent != "."; parent = path.Dir(parent) {
		if matchesAnyGlob(parent, globs) {
			return true
		}
	}
	return false
}

func matchesAnyGlob(candidate string, globs []string) bool {
	for _, glob := range globs {
		if matchesGlob(candidate, glob) {
			return true
		}
	}
	return false
}

func matchesGlob(candidate, glob string) bool {
	pattern := "^"
	for index := 0; index < len(glob); index++ {
		switch glob[index] {
		case '*':
			if index+1 < len(glob) && glob[index+1] == '*' {
				index++
				if index+1 < len(glob) && glob[index+1] == '/' {
					pattern += "(?:.*/)?"
					index++
				} else {
					pattern += ".*"
				}
			} else {
				pattern += "[^/]*"
			}
		case '?':
			pattern += "[^/]"
		case '[':
			end := strings.IndexByte(glob[index+1:], ']')
			if end < 0 {
				return false
			}
			end += index + 1
			pattern += glob[index : end+1]
			index = end
		default:
			pattern += regexp.QuoteMeta(string(glob[index]))
		}
	}
	compiled, err := regexp.Compile(pattern + "$")
	return err == nil && compiled.MatchString(candidate)
}

func arcanaNodeForPath(relativePath string, line int) arcana.Node {
	node := arcana.Node{Kind: "file", Path: relativePath, Name: filepath.Base(relativePath)}
	if line > 0 {
		node.Span = &arcana.Span{Path: relativePath, StartLine: line, EndLine: line}
	}
	return node
}
