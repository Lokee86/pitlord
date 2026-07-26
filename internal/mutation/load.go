package mutation

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

func Load(path string) (Manifest, error) {
	file, err := os.Open(path)
	if err != nil {
		return Manifest{}, fmt.Errorf("open Homunculus mutation manifest: %w", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	var manifest Manifest
	if err := decoder.Decode(&manifest); err != nil {
		return Manifest{}, fmt.Errorf("decode Homunculus mutation manifest: %w", err)
	}
	if err := manifest.Validate(); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

func (manifest *Manifest) Validate() error {
	if manifest.Version != ContractVersion {
		return fmt.Errorf("unsupported Homunculus mutation contract version %d", manifest.Version)
	}
	manifest.MutationID = strings.TrimSpace(manifest.MutationID)
	manifest.Kind = strings.TrimSpace(manifest.Kind)
	manifest.Language = strings.TrimSpace(manifest.Language)
	manifest.BaseCommit = strings.TrimSpace(manifest.BaseCommit)
	manifest.SourceRepository = strings.TrimSpace(manifest.SourceRepository)
	manifest.Worktree = strings.TrimSpace(manifest.Worktree)
	if manifest.MutationID == "" || manifest.Kind == "" || manifest.Language == "" || manifest.BaseCommit == "" {
		return fmt.Errorf("Homunculus mutation manifest is missing mutation_id, kind, language, or base_commit")
	}
	if manifest.SourceRepository == "" || manifest.Worktree == "" {
		return fmt.Errorf("Homunculus mutation manifest is missing source_repository or worktree")
	}
	if len(manifest.ExpectedAddedRelationships) == 0 && len(manifest.ExpectedRemovedRelationships) == 0 {
		return fmt.Errorf("Homunculus mutation manifest has no relationship expectations")
	}
	seen := make(map[string]struct{})
	for index := range manifest.ExpectedAddedRelationships {
		if err := normalizeRelationship(&manifest.ExpectedAddedRelationships[index]); err != nil {
			return fmt.Errorf("invalid expected added relationship %d: %w", index+1, err)
		}
		key := "added\x00" + relationshipKey(manifest.ExpectedAddedRelationships[index])
		if _, exists := seen[key]; exists {
			return fmt.Errorf("duplicate expected added relationship %d", index+1)
		}
		seen[key] = struct{}{}
	}
	for index := range manifest.ExpectedRemovedRelationships {
		if err := normalizeRelationship(&manifest.ExpectedRemovedRelationships[index]); err != nil {
			return fmt.Errorf("invalid expected removed relationship %d: %w", index+1, err)
		}
		key := "removed\x00" + relationshipKey(manifest.ExpectedRemovedRelationships[index])
		if _, exists := seen[key]; exists {
			return fmt.Errorf("duplicate expected removed relationship %d", index+1)
		}
		seen[key] = struct{}{}
	}
	sortRelationships(manifest.ExpectedAddedRelationships)
	sortRelationships(manifest.ExpectedRemovedRelationships)
	return nil
}

func normalizeRelationship(relationship *Relationship) error {
	relationship.Source = strings.TrimSpace(relationship.Source)
	relationship.Relation = strings.ToLower(strings.TrimSpace(relationship.Relation))
	relationship.Target = strings.TrimSpace(relationship.Target)
	relationship.SourceArea = strings.TrimSpace(relationship.SourceArea)
	relationship.TargetArea = strings.TrimSpace(relationship.TargetArea)
	if relationship.Source == "" || relationship.Relation == "" || relationship.Target == "" {
		return fmt.Errorf("source, relation, and target are required")
	}
	if !strings.Contains(relationship.Source, "::") || !strings.Contains(relationship.Target, "::") {
		return fmt.Errorf("source and target must use Lexicon path::name qualified symbols")
	}
	return nil
}

func sortRelationships(relationships []Relationship) {
	sort.Slice(relationships, func(i, j int) bool {
		return relationshipKey(relationships[i]) < relationshipKey(relationships[j])
	})
}

func relationshipKey(relationship Relationship) string {
	return strings.Join([]string{
		relationship.Source,
		relationship.Relation,
		relationship.Target,
		relationship.SourceArea,
		relationship.TargetArea,
	}, "\x00")
}
