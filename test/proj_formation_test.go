package test

import (
	"path/filepath"
	"testing"

	"github.com/jaqmol/approx/project"
)

// TestComplexProjectFormation ...
func TestComplexProjectFormation(t *testing.T) {
	projDir, err := filepath.Abs("complx-test-proj") // /formation.yaml
	if err != nil {
		t.Fatal(err)
	}
	form, err := project.LoadFormation(projDir)
	if err != nil {
		t.Fatal(err)
	}

	checkProjectDefinitions(t, form.Definitions, true)
	checkProjectFlows(t, form.Flows, true)
}

// TestSimpleProjectFormation ...
func TestSimpleProjectFormation(t *testing.T) {
	projDir, err := filepath.Abs("simpl-test-proj") // /definition.yaml /flow.yaml
	if err != nil {
		t.Fatal(err)
	}
	form, err := project.LoadFormation(projDir)
	if err != nil {
		t.Fatal(err)
	}

	checkProjectDefinitions(t, form.Definitions, false)
	checkProjectFlows(t, form.Flows, false)
}
