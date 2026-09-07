package scan

import "testing"

func TestBuiltInAnalyzerResolvesExactAnalyzerAndGraphRequirement(t *testing.T) {
	semantic, err := BuiltInAnalyzer(AnalyzerSwallowedError)
	if err != nil {
		t.Fatal(err)
	}
	if !semantic.Metadata().RequiresGraph {
		t.Fatal("swallowed-error should require graph state")
	}

	clippy, err := BuiltInAnalyzer(AnalyzerClippy)
	if err != nil {
		t.Fatal(err)
	}
	if clippy.Metadata().RequiresGraph {
		t.Fatal("clippy should remain graph-free")
	}

	if _, err := BuiltInAnalyzer("missing-analyzer"); err == nil {
		t.Fatal("unknown analyzer should fail closed")
	}
}
