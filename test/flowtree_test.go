package test

import (
	"path/filepath"
	"testing"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/project"
)

// TestFlowTree ...
func TestFlowTree(t *testing.T) {
	projDir, err := filepath.Abs("simpl-test-proj") // /flow.yaml
	if err != nil {
		t.Fatal(err)
	}
	flows, err := project.LoadFlow(projDir)
	if err != nil {
		t.Fatal(err)
	}

	config := makeSimpleSequenceConfig()
	procs := map[string]configuration.Processor{
		config.fork.Ident:             &config.fork,
		config.firstNameExtract.Ident: &config.firstNameExtract,
		config.lastNameExtract.Ident:  &config.lastNameExtract,
		config.merge.Ident:            &config.merge,
	}
	tree, err := configuration.NewFlowTree(flows, procs)
	if err != nil {
		t.Fatal(err)
	}
	// TODO: CONTINUE HERE
}
