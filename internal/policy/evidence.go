package policy

import "sort"

func diagnosticFromEvidence(rule Rule, evidence []Evidence) *Diagnostic {
	if len(evidence) == 0 {
		return nil
	}
	sortEvidence(evidence)
	return &Diagnostic{
		RuleID:   rule.ID,
		Message:  rule.Message,
		Severity: rule.Severity,
		Evidence: evidence,
	}
}

func sortEvidence(evidence []Evidence) {
	sort.Slice(evidence, func(i, j int) bool {
		left, right := evidence[i], evidence[j]
		if left.Source.Path != right.Source.Path {
			return left.Source.Path < right.Source.Path
		}
		leftLine, rightLine := sourceLine(left), sourceLine(right)
		if leftLine != rightLine {
			return leftLine < rightLine
		}
		if left.Source.Name != right.Source.Name {
			return left.Source.Name < right.Source.Name
		}
		if left.Issue != right.Issue {
			return left.Issue < right.Issue
		}
		if left.SourceArea != right.SourceArea {
			return left.SourceArea < right.SourceArea
		}
		if left.TargetArea != right.TargetArea {
			return left.TargetArea < right.TargetArea
		}
		if left.Relation != right.Relation {
			return left.Relation < right.Relation
		}
		leftTargetPath, rightTargetPath := targetPath(left), targetPath(right)
		if leftTargetPath != rightTargetPath {
			return leftTargetPath < rightTargetPath
		}
		return targetName(left) < targetName(right)
	})
}

func sourceLine(evidence Evidence) int {
	if evidence.Source.Span == nil {
		return 0
	}
	return evidence.Source.Span.StartLine
}

func targetPath(evidence Evidence) string {
	if evidence.Target == nil {
		return ""
	}
	return evidence.Target.Path
}

func targetName(evidence Evidence) string {
	if evidence.Target == nil {
		return ""
	}
	return evidence.Target.Name
}
