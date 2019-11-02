package test

import (
	"fmt"
	"path/filepath"
	"strings"
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

	visited := make(map[string]bool)

	checkLen := lengthChecker(map[string]int{
		config.fork.ID():             0,
		config.firstNameExtract.ID(): 1,
		config.lastNameExtract.ID():  1,
		config.merge.ID():            2,
	})

	tree.Iterate(func(prev []*configuration.FlowNode, curr *configuration.FlowNode) {
		id := curr.Processor().ID()
		visited[id] = true
		if err := checkLen(id, len(prev)); err != nil {
			t.Fatal(err)
		}
	})

	if len(visited) != 4 {
		t.Fatal("Expected to visit 4 nodes, but got:", len(visited))
	}
	errs := checkContainsAll(
		visited,
		config.fork.ID(),
		config.firstNameExtract.ID(),
		config.lastNameExtract.ID(),
		config.merge.ID(),
	)
	if len(errs) > 0 {
		err := fmt.Errorf("Errors visiting nodes: %v", strings.Join(errorsToStrings(errs), ", "))
		t.Fatal(err)
	}
}
