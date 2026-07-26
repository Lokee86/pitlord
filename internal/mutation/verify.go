package mutation

import (
	"context"
	"fmt"

	"github.com/Lokee86/pitlord/internal/arcana"
)

type VerifyRequest struct {
	ManifestPath     string
	BaselineSnapshot string
	MutatedSnapshot  string
	ArcanaCommand    string
}

func Verify(ctx context.Context, request VerifyRequest) (Result, error) {
	manifest, err := Load(request.ManifestPath)
	if err != nil {
		return Result{}, err
	}
	if request.BaselineSnapshot == "" || request.MutatedSnapshot == "" {
		return Result{}, fmt.Errorf("baseline and mutated Arcana snapshots are required")
	}
	client := arcana.Client{Command: request.ArcanaCommand}
	checks := make([]Check, 0, len(manifest.ExpectedAddedRelationships)+len(manifest.ExpectedRemovedRelationships))
	for _, relationship := range manifest.ExpectedAddedRelationships {
		check, err := verifyRelationship(ctx, client, request, "added", relationship)
		if err != nil {
			return Result{}, err
		}
		checks = append(checks, check)
	}
	for _, relationship := range manifest.ExpectedRemovedRelationships {
		check, err := verifyRelationship(ctx, client, request, "removed", relationship)
		if err != nil {
			return Result{}, err
		}
		checks = append(checks, check)
	}
	matched := true
	for _, check := range checks {
		if !check.Matched {
			matched = false
			break
		}
	}
	return Result{
		Schema:           "pitlord.mutation-verification.v1",
		ManifestPath:     request.ManifestPath,
		MutationID:       manifest.MutationID,
		Kind:             manifest.Kind,
		Language:         manifest.Language,
		BaseCommit:       manifest.BaseCommit,
		BaselineSnapshot: request.BaselineSnapshot,
		MutatedSnapshot:  request.MutatedSnapshot,
		Checks:           checks,
		Matched:          matched,
	}, nil
}

func verifyRelationship(
	ctx context.Context,
	client arcana.Client,
	request VerifyRequest,
	expectation string,
	relationship Relationship,
) (Check, error) {
	baselinePresent, err := client.HasQualifiedRelationship(
		ctx,
		request.BaselineSnapshot,
		relationship.Source,
		relationship.Relation,
		relationship.Target,
	)
	if err != nil {
		return Check{}, fmt.Errorf("check baseline relationship %s: %w", relationshipKey(relationship), err)
	}
	mutatedPresent, err := client.HasQualifiedRelationship(
		ctx,
		request.MutatedSnapshot,
		relationship.Source,
		relationship.Relation,
		relationship.Target,
	)
	if err != nil {
		return Check{}, fmt.Errorf("check mutated relationship %s: %w", relationshipKey(relationship), err)
	}
	matched := false
	switch expectation {
	case "added":
		matched = !baselinePresent && mutatedPresent
	case "removed":
		matched = baselinePresent && !mutatedPresent
	}
	return Check{
		Expectation:     expectation,
		Relationship:    relationship,
		BaselinePresent: baselinePresent,
		MutatedPresent:  mutatedPresent,
		Matched:         matched,
	}, nil
}
