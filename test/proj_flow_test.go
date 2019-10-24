package test

import (
	"path/filepath"
	"testing"

	"github.com/jaqmol/approx/project"
)

// TestProjectFlow ...
func TestProjectFlow(t *testing.T) {
	projDir, err := filepath.Abs("simpl-test-proj") // /flow.yaml
	if err != nil {
		t.Fatal(err)
	}
	flow, err := project.LoadFlow(projDir)
	if err != nil {
		t.Fatal(err)
	}

	expected := [][]string{
		[]string{"fork", "extract-first-name", "merge"},
		[]string{"fork", "extract-last-name", "merge"},
	}

	for i, f := range flow {
		for j, name := range f {
			expName := expected[i][j]
			if name != expName {
				t.Fatalf("Expected \"%v\", but got \"%v\"", expName, name)
			}
		}
	}
}
