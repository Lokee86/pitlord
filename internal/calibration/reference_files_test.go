package calibration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRepositoryCalibrationReferencesDecode(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("..", "..", "docs", "development", "calibration", "references", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 64 {
		t.Fatalf("reference files = %d, want 64", len(files))
	}
	for _, filename := range files {
		filename := filename
		t.Run(filepath.Base(filename), func(t *testing.T) {
			file, err := os.Open(filename)
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()
			if _, err := DecodeReference(file); err != nil {
				t.Fatal(err)
			}
		})
	}
}
