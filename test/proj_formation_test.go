package test

import (
	"path/filepath"
	"testing"

	"github.com/jaqmol/approx/project"
)

// TestComplexProjectFormation ...
func TestComplexProjectFormation(t *testing.T) {
	// t.SkipNow()
	projDir, err := filepath.Abs("complx-test-proj") // /formation.yaml
	if err != nil {
		t.Fatal(err)
	}
	form, err := project.LoadFormation(projDir)
	if err != nil {
		t.Fatal(err)
	}

	checkProjectDefinitions(t, form.Definitions)
	checkProjectFlows(t, form.Flows)
}

// TestSimpleProjectFormation ...
func TestSimpleProjectFormation(t *testing.T) {
	// t.SkipNow()
	projDir, err := filepath.Abs("simpl-test-proj") // /definition.yaml /flow.yaml
	if err != nil {
		t.Fatal(err)
	}
	form, err := project.LoadFormation(projDir)
	if err != nil {
		t.Fatal(err)
	}

	checkProjectDefinitions(t, form.Definitions)
	checkProjectFlows(t, form.Flows)
}
