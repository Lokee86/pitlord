package policy

import (
	"fmt"
	"regexp"
	"strings"
)

func (document *Document) NormalizeAndValidate() error {
	if document.Version == 0 {
		document.Version = Version
	}
	if document.Version != Version {
		return fmt.Errorf("unsupported policy version %d", document.Version)
	}
	document.Includes = normalizeIncludes(document.Includes)
	if len(document.Rules) == 0 {
		return fmt.Errorf("policy must contain at least one rule")
	}

	areas, err := document.normalizeAreas()
	if err != nil {
		return err
	}
	seenRules := make(map[string]struct{}, len(document.Rules))
	for index := range document.Rules {
		rule := &document.Rules[index]
		if err := normalizeRuleCommon(rule, index, seenRules); err != nil {
			return err
		}
		switch rule.Type {
		case RuleForbidDependency:
			if err := normalizeDependencyRule(rule, areas); err != nil {
				return err
			}
		case RuleRequireOwnership:
			if err := normalizeOwnershipRule(rule, len(areas)); err != nil {
				return err
			}
		case RuleForbidAreaCycles:
			if err := normalizeAreaCycleRule(rule, areas); err != nil {
				return err
			}
		case RuleForbidContent:
			if err := normalizeContentRule(rule); err != nil {
				return err
			}
		case RuleRequirePath, RuleForbidPath:
			if err := normalizePathRule(rule); err != nil {
				return err
			}
		default:
			return fmt.Errorf("rule %q has unsupported type %q", rule.ID, rule.Type)
		}
	}
	return nil
}

func (document *Document) normalizeAreas() (map[string]Area, error) {
	areas := make(map[string]Area, len(document.Areas))
	for index := range document.Areas {
		area := &document.Areas[index]
		area.ID = strings.TrimSpace(area.ID)
		area.Description = strings.TrimSpace(area.Description)
		if area.ID == "" {
			return nil, fmt.Errorf("area %d has no id", index+1)
		}
		if _, exists := areas[area.ID]; exists {
			return nil, fmt.Errorf("duplicate area id %q", area.ID)
		}
		area.Paths = normalizePrefixes(area.Paths)
		area.ExcludePaths = normalizePrefixes(area.ExcludePaths)
		area.Kinds = normalizeKinds(area.Kinds)
		if len(area.Paths) == 0 {
			return nil, fmt.Errorf("area %q has no paths", area.ID)
		}
		areas[area.ID] = *area
	}
	return areas, nil
}

func normalizeRuleCommon(rule *Rule, index int, seen map[string]struct{}) error {
	rule.ID = strings.TrimSpace(rule.ID)
	if rule.ID == "" {
		return fmt.Errorf("rule %d has no id", index+1)
	}
	if _, exists := seen[rule.ID]; exists {
		return fmt.Errorf("duplicate rule id %q", rule.ID)
	}
	seen[rule.ID] = struct{}{}

	rule.Type = strings.ToLower(strings.TrimSpace(rule.Type))
	if rule.Type == "" {
		rule.Type = RuleForbidDependency
	}
	rule.Severity = strings.ToLower(strings.TrimSpace(rule.Severity))
	if rule.Severity == "" {
		rule.Severity = "error"
	}
	if rule.Severity != "error" && rule.Severity != "warning" {
		return fmt.Errorf("rule %q has invalid severity %q", rule.ID, rule.Severity)
	}
	rule.Message = strings.TrimSpace(rule.Message)
	rule.FromAreas = normalizeNames(rule.FromAreas)
	rule.ToAreas = normalizeNames(rule.ToAreas)
	rule.FromPaths = normalizePrefixes(rule.FromPaths)
	rule.ToPaths = normalizePrefixes(rule.ToPaths)
	rule.FromExcludePaths = normalizePrefixes(rule.FromExcludePaths)
	rule.ToExcludePaths = normalizePrefixes(rule.ToExcludePaths)
	rule.SourceKinds = normalizeKinds(rule.SourceKinds)
	rule.TargetKinds = normalizeKinds(rule.TargetKinds)
	rule.Relations = normalizeRelations(rule.Relations)
	rule.Areas = normalizeNames(rule.Areas)
	rule.ScopePaths = normalizePrefixes(rule.ScopePaths)
	rule.ScopeExcludePaths = normalizePrefixes(rule.ScopeExcludePaths)
	rule.IncludePaths = normalizeGlobs(rule.IncludePaths)
	rule.ExcludePaths = normalizeGlobs(rule.ExcludePaths)
	rule.Path = normalizeGlob(rule.Path)
	rule.PatternType = strings.ToLower(strings.TrimSpace(rule.PatternType))
	return nil
}

func normalizeContentRule(rule *Rule) error {
	if len(rule.IncludePaths) == 0 {
		return fmt.Errorf("rule %q has no include_paths", rule.ID)
	}
	for _, glob := range append(append([]string(nil), rule.IncludePaths...), rule.ExcludePaths...) {
		if err := validateRepoGlob(glob); err != nil {
			return fmt.Errorf("rule %q: %w", rule.ID, err)
		}
	}
	if rule.Path != "" || len(rule.ScopePaths) > 0 || len(rule.ScopeExcludePaths) > 0 || rule.AllowUnowned || rule.AllowOverlaps {
		return fmt.Errorf("rule %q uses fields incompatible with content checks", rule.ID)
	}
	if len(rule.FromAreas) > 0 || len(rule.ToAreas) > 0 || len(rule.FromPaths) > 0 || len(rule.ToPaths) > 0 || len(rule.Relations) > 0 || len(rule.Areas) > 0 {
		return fmt.Errorf("rule %q uses graph fields", rule.ID)
	}
	patternCount := 0
	if rule.Literal != "" {
		patternCount++
	}
	if rule.Regex != "" {
		patternCount++
	}
	if rule.Pattern != "" {
		patternCount++
	}
	if patternCount != 1 {
		return fmt.Errorf("rule %q must specify exactly one literal, regex, or pattern", rule.ID)
	}
	if rule.Literal != "" || rule.Regex != "" {
		if rule.PatternType != "" {
			return fmt.Errorf("rule %q cannot combine pattern_type with literal or regex", rule.ID)
		}
		if rule.Regex != "" {
			if _, err := regexp.Compile(rule.Regex); err != nil {
				return fmt.Errorf("rule %q has invalid regex: %w", rule.ID, err)
			}
		}
	} else {
		if rule.PatternType == "" {
			rule.PatternType = "literal"
		}
		if rule.PatternType != "literal" && rule.PatternType != "regex" {
			return fmt.Errorf("rule %q has invalid pattern_type %q", rule.ID, rule.PatternType)
		}
		if rule.PatternType == "regex" {
			if _, err := regexp.Compile(rule.Pattern); err != nil {
				return fmt.Errorf("rule %q has invalid regex: %w", rule.ID, err)
			}
		}
	}
	if rule.Message == "" {
		rule.Message = fmt.Sprintf("content matching %s is forbidden", contentPatternLabel(*rule))
	}
	return nil
}

func normalizePathRule(rule *Rule) error {
	if rule.Path == "" {
		return fmt.Errorf("rule %q has no path", rule.ID)
	}
	if err := validateRepoGlob(rule.Path); err != nil {
		return fmt.Errorf("rule %q: %w", rule.ID, err)
	}
	if len(rule.FromAreas) > 0 || len(rule.ToAreas) > 0 || len(rule.FromPaths) > 0 || len(rule.ToPaths) > 0 || len(rule.Relations) > 0 || len(rule.Areas) > 0 || len(rule.ScopePaths) > 0 || len(rule.ScopeExcludePaths) > 0 || len(rule.IncludePaths) > 0 || len(rule.ExcludePaths) > 0 || rule.Literal != "" || rule.Regex != "" || rule.Pattern != "" || rule.PatternType != "" || rule.AllowUnowned || rule.AllowOverlaps {
		return fmt.Errorf("rule %q uses fields incompatible with path checks", rule.ID)
	}
	if rule.Message == "" {
		rule.Message = fmt.Sprintf("path %s is forbidden or missing", rule.Path)
	}
	return nil
}

func contentPatternLabel(rule Rule) string {
	if rule.Literal != "" {
		return "literal " + rule.Literal
	}
	if rule.Regex != "" {
		return "regex " + rule.Regex
	}
	return rule.PatternType + " " + rule.Pattern
}

func normalizeDependencyRule(rule *Rule, areas map[string]Area) error {
	if len(rule.FromPaths) == 0 && len(rule.FromAreas) == 0 {
		return fmt.Errorf("rule %q has no source selector", rule.ID)
	}
	if len(rule.ToPaths) == 0 && len(rule.ToAreas) == 0 {
		return fmt.Errorf("rule %q has no target selector", rule.ID)
	}
	if err := validateAreaReferences(rule.ID, "from_areas", rule.FromAreas, areas); err != nil {
		return err
	}
	if err := validateAreaReferences(rule.ID, "to_areas", rule.ToAreas, areas); err != nil {
		return err
	}
	if len(rule.Areas) > 0 {
		return fmt.Errorf("rule %q uses area-cycle-only fields", rule.ID)
	}
	if len(rule.ScopePaths) > 0 || len(rule.ScopeExcludePaths) > 0 || rule.AllowUnowned || rule.AllowOverlaps {
		return fmt.Errorf("rule %q uses ownership-only fields", rule.ID)
	}
	if len(rule.Relations) == 0 {
		rule.Relations = append([]string(nil), defaultRelations...)
	}
	if rule.Message == "" {
		rule.Message = fmt.Sprintf(
			"%s must not depend on %s",
			selectorLabel(rule.FromAreas, rule.FromPaths),
			selectorLabel(rule.ToAreas, rule.ToPaths),
		)
	}
	return nil
}

func normalizeOwnershipRule(rule *Rule, areaCount int) error {
	if areaCount == 0 {
		return fmt.Errorf("rule %q requires at least one declared area", rule.ID)
	}
	if len(rule.ScopePaths) == 0 {
		return fmt.Errorf("rule %q has no scope_paths", rule.ID)
	}
	if len(rule.FromAreas) > 0 || len(rule.ToAreas) > 0 || len(rule.FromPaths) > 0 || len(rule.ToPaths) > 0 || len(rule.Relations) > 0 || len(rule.TargetKinds) > 0 || len(rule.Areas) > 0 {
		return fmt.Errorf("rule %q uses dependency or area-cycle fields", rule.ID)
	}
	if rule.AllowUnowned && rule.AllowOverlaps {
		return fmt.Errorf("rule %q allows both unowned and overlapping nodes and would enforce nothing", rule.ID)
	}
	if len(rule.SourceKinds) == 0 {
		rule.SourceKinds = []string{"file"}
	}
	if rule.Message == "" {
		rule.Message = fmt.Sprintf(
			"nodes in %s must belong to exactly one declared area",
			strings.Join(rule.ScopePaths, ", "),
		)
	}
	return nil
}

func normalizeAreaCycleRule(rule *Rule, areas map[string]Area) error {
	if len(areas) < 2 {
		return fmt.Errorf("rule %q requires at least two declared areas", rule.ID)
	}
	if len(rule.Areas) == 0 {
		for areaID := range areas {
			rule.Areas = append(rule.Areas, areaID)
		}
		rule.Areas = normalizeNames(rule.Areas)
	}
	if len(rule.Areas) < 2 {
		return fmt.Errorf("rule %q requires at least two selected areas", rule.ID)
	}
	if err := validateAreaReferences(rule.ID, "areas", rule.Areas, areas); err != nil {
		return err
	}
	if len(rule.FromAreas) > 0 || len(rule.ToAreas) > 0 || len(rule.FromPaths) > 0 || len(rule.ToPaths) > 0 || len(rule.FromExcludePaths) > 0 || len(rule.ToExcludePaths) > 0 || len(rule.ScopePaths) > 0 || len(rule.ScopeExcludePaths) > 0 || rule.AllowUnowned || rule.AllowOverlaps {
		return fmt.Errorf("rule %q uses fields incompatible with area-cycle checks", rule.ID)
	}
	if len(rule.Relations) == 0 {
		rule.Relations = append([]string(nil), defaultRelations...)
	}
	if rule.Message == "" {
		rule.Message = fmt.Sprintf(
			"selected architecture areas must not form dependency cycles: %s",
			strings.Join(rule.Areas, ", "),
		)
	}
	return nil
}

func validateAreaReferences(ruleID, field string, references []string, areas map[string]Area) error {
	for _, reference := range references {
		if _, exists := areas[reference]; !exists {
			return fmt.Errorf("rule %q %s references unknown area %q", ruleID, field, reference)
		}
	}
	return nil
}

func selectorLabel(areas, paths []string) string {
	parts := make([]string, 0, len(areas)+len(paths))
	for _, area := range areas {
		parts = append(parts, "area:"+area)
	}
	parts = append(parts, paths...)
	return strings.Join(parts, ", ")
}
