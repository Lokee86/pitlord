package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/Lokee86/pitlord/internal/checker"
)

type evaluationOptions struct {
	repository *string
	snapshot   *string
	policy     *string
	manifest   *string
	arcana     *string
	timeout    *time.Duration
}

func bindEvaluationOptions(flags *flag.FlagSet) evaluationOptions {
	return evaluationOptions{
		repository: flags.String("repo", ".", "repository root containing .arcana/CURRENT"),
		snapshot:   flags.String("snapshot", "", "explicit Arcana snapshot directory"),
		policy:     flags.String("policy", "", "Pitlord JSON policy"),
		manifest:   flags.String("homunculus-manifest", "", "Homunculus manifest used as policy and expected diagnostics"),
		arcana:     flags.String("arcana", "", "Arcana executable override"),
		timeout:    flags.Duration("timeout", 2*time.Minute, "maximum evaluation duration"),
	}
}

func (options evaluationOptions) validate() error {
	if (*options.policy == "") == (*options.manifest == "") {
		return fmt.Errorf("exactly one of --policy or --homunculus-manifest is required")
	}
	if *options.timeout <= 0 {
		return fmt.Errorf("--timeout must be positive")
	}
	return nil
}

func (options evaluationOptions) run() (checker.Result, error) {
	if err := options.validate(); err != nil {
		return checker.Result{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), *options.timeout)
	defer cancel()
	return checker.Run(ctx, checker.Request{
		Repository:         *options.repository,
		Snapshot:           *options.snapshot,
		PolicyPath:         *options.policy,
		HomunculusManifest: *options.manifest,
		ArcanaCommand:      *options.arcana,
	})
}
