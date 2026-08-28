package scan

import "testing"

func TestClassifySourceRoleSeparatesCommonNonProductionTrees(t *testing.T) {
	tests := []struct {
		path string
		want sourceRole
	}{
		{"src/main/java/example/App.java", sourceRoleProduction},
		{"src/test/java/example/AppTest.java", sourceRoleNonProduction},
		{"tests/Dapper.Tests/AsyncTests.cs", sourceRoleNonProduction},
		{"bench/Polly.Core.Benchmarks/BridgeBenchmark.cs", sourceRoleNonProduction},
		{"benchmarks/perf.go", sourceRoleNonProduction},
		{"src/Spectre.Console.Tests/Fixture.cs", sourceRoleNonProduction},
		{"impl/maven-testing/src/main/java/Fixture.java", sourceRoleNonProduction},
		{"test-shrinker/src/main/java/Fixture.java", sourceRoleNonProduction},
		{"src/Snippets/Docs/Chaos.Fault.cs", sourceRoleNonProduction},
		{"jmh-samples/src/main/java/Sample.java", sourceRoleNonProduction},
		{"services/game-server/internal/game/collision_spatial_index_test.go", sourceRoleNonProduction},
		{"services/game-server/internal/tooling/controller.go", sourceRoleNonProduction},
		{"src/widget.spec.ts", sourceRoleNonProduction},
		{"test_parser.py", sourceRoleNonProduction},
	}
	for _, test := range tests {
		if got := classifySourceRole(test.path); got != test.want {
			t.Errorf("classifySourceRole(%q) = %q, want %q", test.path, got, test.want)
		}
	}
}
