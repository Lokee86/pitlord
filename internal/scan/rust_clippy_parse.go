package scan

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

type cargoMessage struct {
	Reason  string             `json:"reason"`
	Message compilerDiagnostic `json:"message"`
}

type compilerDiagnostic struct {
	Code     *diagnosticCode      `json:"code"`
	Level    string               `json:"level"`
	Message  string               `json:"message"`
	Spans    []diagnosticSpan     `json:"spans"`
	Children []compilerDiagnostic `json:"children"`
}

type diagnosticCode struct {
	Code string `json:"code"`
}

type diagnosticSpan struct {
	FileName                string  `json:"file_name"`
	LineStart               int     `json:"line_start"`
	LineEnd                 int     `json:"line_end"`
	ColumnStart             int     `json:"column_start"`
	ColumnEnd               int     `json:"column_end"`
	IsPrimary               bool    `json:"is_primary"`
	SuggestedReplacement    *string `json:"suggested_replacement"`
	SuggestionApplicability *string `json:"suggestion_applicability"`
}

func parseClippyOutput(output []byte, repositoryRoot, manifestDirectory, pathPrefix string) ([]Finding, error) {
	byID := make(map[string]Finding)
	scanner := bufio.NewScanner(bytes.NewReader(output))
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		var message cargoMessage
		if err := json.Unmarshal(scanner.Bytes(), &message); err != nil {
			return nil, fmt.Errorf("decode cargo JSON output: %w", err)
		}
		finding, ok := clippyFinding(message, repositoryRoot, manifestDirectory)
		if !ok || !matchesPathPrefix(finding.Scope.Path, pathPrefix) {
			continue
		}
		byID[finding.ID] = finding
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read cargo JSON output: %w", err)
	}
	findings := make([]Finding, 0, len(byID))
	for _, finding := range byID {
		findings = append(findings, finding)
	}
	sort.Slice(findings, func(left, right int) bool { return findings[left].ID < findings[right].ID })
	return findings, nil
}

func clippyFinding(message cargoMessage, repositoryRoot, manifestDirectory string) (Finding, bool) {
	if message.Reason != "compiler-message" || message.Message.Code == nil ||
		!strings.HasPrefix(message.Message.Code.Code, "clippy::") {
		return Finding{}, false
	}
	primary, ok := primaryDiagnosticSpan(message.Message.Spans)
	if !ok {
		return Finding{}, false
	}
	path := clippyRepositoryPath(repositoryRoot, manifestDirectory, primary.FileName)
	location := sourceSpan(primary, path)
	ruleID := message.Message.Code.Code
	scopeKey := fmt.Sprintf("%s\x00%d\x00%d\x00%s", path, primary.LineStart, primary.ColumnStart, message.Message.Message)
	finding := Finding{
		ID:          findingID(AnalyzerClippy, ruleID+"\x00"+scopeKey),
		RuleID:      ruleID,
		Detector:    AnalyzerClippy,
		Disposition: DispositionAdvisory,
		Severity:    clippySeverity(message.Message.Level),
		Scope:       Scope{Kind: "file", Path: path},
		Location:    &location,
		Summary:     message.Message.Message,
		Evidence:    []Evidence{{Kind: "clippy", Message: ruleID}},
	}
	finding.SuggestedFix = clippySuggestedFix(message.Message, repositoryRoot, manifestDirectory)
	return finding, true
}

func primaryDiagnosticSpan(spans []diagnosticSpan) (diagnosticSpan, bool) {
	for _, span := range spans {
		if span.IsPrimary {
			return span, true
		}
	}
	if len(spans) > 0 {
		return spans[0], true
	}
	return diagnosticSpan{}, false
}

func clippySuggestedFix(diagnostic compilerDiagnostic, repositoryRoot, manifestDirectory string) *SuggestedFix {
	for _, child := range diagnostic.Children {
		edits, applicability := suggestedEdits(child.Spans, repositoryRoot, manifestDirectory)
		if len(edits) > 0 {
			return &SuggestedFix{Message: child.Message, Applicability: applicability, Edits: edits}
		}
	}
	edits, applicability := suggestedEdits(diagnostic.Spans, repositoryRoot, manifestDirectory)
	if len(edits) == 0 {
		return nil
	}
	return &SuggestedFix{Message: "Apply Clippy suggestion", Applicability: applicability, Edits: edits}
}

func suggestedEdits(spans []diagnosticSpan, repositoryRoot, manifestDirectory string) ([]TextEdit, Applicability) {
	var edits []TextEdit
	applicability := ApplicabilityMachine
	for _, span := range spans {
		if span.SuggestedReplacement == nil {
			continue
		}
		path := clippyRepositoryPath(repositoryRoot, manifestDirectory, span.FileName)
		edits = append(edits, TextEdit{Span: sourceSpan(span, path), Replacement: *span.SuggestedReplacement})
		applicability = lessSafeApplicability(applicability, clippyApplicability(span.SuggestionApplicability))
	}
	return edits, applicability
}

func sourceSpan(span diagnosticSpan, path string) SourceSpan {
	return SourceSpan{Path: path, StartLine: span.LineStart, StartColumn: span.ColumnStart, EndLine: span.LineEnd, EndColumn: span.ColumnEnd}
}

func clippyRepositoryPath(repositoryRoot, manifestDirectory, fileName string) string {
	path := fileName
	if !filepath.IsAbs(path) {
		path = filepath.Join(manifestDirectory, path)
	}
	if relative, err := filepath.Rel(repositoryRoot, path); err == nil && pathWithinRoot(repositoryRoot, path) {
		return normalizedUnitPath(filepath.ToSlash(relative))
	}
	return normalizedUnitPath(filepath.ToSlash(fileName))
}

func matchesPathPrefix(path, prefix string) bool {
	prefix = normalizedUnitPath(prefix)
	return prefix == "" || prefix == "." || path == prefix || strings.HasPrefix(path, prefix+"/")
}
