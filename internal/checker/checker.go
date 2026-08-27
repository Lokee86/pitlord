package checker

import (
	"context"
	"fmt"

	"github.com/Lokee86/pitlord/internal/arcana"
	"github.com/Lokee86/pitlord/internal/policy"
	"github.com/Lokee86/pitlord/internal/snapshot"
)

type Request struct {
	Repository         string
	Snapshot           string
	PolicyPath         string
	HomunculusManifest string
	ArcanaCommand      string
}

type Result struct {
	Document policy.Document
	Graph    arcana.Graph
	Report   policy.Report
}

func Run(ctx context.Context, request Request) (Result, error) {
	if (request.PolicyPath == "") == (request.HomunculusManifest == "") {
		return Result{}, fmt.Errorf("exactly one policy source is required")
	}

	var document policy.Document
	var expected []string
	var source string
	var err error
	if request.HomunculusManifest != "" {
		document, expected, err = policy.FromHomunculus(request.HomunculusManifest)
		source = request.HomunculusManifest
	} else {
		document, err = policy.Load(request.PolicyPath)
		source = request.PolicyPath
	}
	if err != nil {
		return Result{}, err
	}

	repositoryRoot := request.Repository
	if repositoryRoot == "" {
		repositoryRoot = "."
	}
	snapshotPath := ""
	if policy.RequiresGraph(document) {
		snapshotPath, err = snapshot.Resolve(repositoryRoot, request.Snapshot)
		if err != nil {
			return Result{}, err
		}
	}

	result, err := Evaluate(ctx, document, source, repositoryRoot, snapshotPath, request.ArcanaCommand)
	if err != nil {
		return Result{}, err
	}
	if request.HomunculusManifest != "" {
		result.Report.Expectation = policy.CompareExpectation(expected, result.Report.Diagnostics)
	}
	return result, nil
}

func Evaluate(
	ctx context.Context,
	document policy.Document,
	policySource string,
	repositoryRoot string,
	snapshotPath string,
	arcanaCommand string,
) (Result, error) {
	graph := arcana.Graph{}
	if policy.RequiresGraph(document) {
		resolvedCommand, err := arcana.ResolveCommand(repositoryRoot, arcanaCommand)
		if err != nil {
			return Result{}, err
		}
		loaded, err := (arcana.Client{Command: resolvedCommand}).LoadGraphWithOptions(
			ctx,
			snapshotPath,
			arcana.LoadOptions{
				SourcePrefixes:   policy.SourcePrefixes(document),
				OutgoingPrefixes: policy.RelationshipSourcePrefixes(document),
			},
		)
		if err != nil {
			return Result{}, err
		}
		graph = loaded
	}

	diagnostics := policy.Evaluate(document, graph, repositoryRoot)
	report := policy.Report{
		Schema:       "pitlord.report.v1",
		Snapshot:     snapshotPath,
		PolicySource: policySource,
		Diagnostics:  diagnostics,
		Summary:      policy.BuildSummary(document, graph, diagnostics),
	}
	return Result{Document: document, Graph: graph, Report: report}, nil
}
