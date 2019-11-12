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
	projDir, err := filepath.Abs("alpha-test-proj") // /flow.yaml
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

	visited := make(map[string]int)

	checkLen := lengthChecker(map[string][]int{
		config.fork.ID():             []int{0, 2},
		config.firstNameExtract.ID(): []int{1, 1},
		config.lastNameExtract.ID():  []int{1, 1},
		config.merge.ID():            []int{2, 0},
	})

	err = tree.Iterate(func(prev []*configuration.FlowNode, curr *configuration.FlowNode, next []*configuration.FlowNode) error {
		id := curr.Processor().ID()
		visited[id]++
		return checkLen(id, len(prev), len(next))
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(visited) != 4 {
		t.Fatal("Expected to visit 4 nodes, but got:", len(visited))
	}
	errs := checkContainsAllTimes(
		visited,
		map[string]int{
			config.fork.ID():             1,
			config.firstNameExtract.ID(): 1,
			config.lastNameExtract.ID():  1,
			config.merge.ID():            1,
		},
	)
	if len(errs) > 0 {
		err := fmt.Errorf("Errors visiting nodes: %v", strings.Join(errorsToStrings(errs), ", "))
		t.Fatal(err)
	}

	if tree.Input == nil {
		t.Fatal("Expected node tree to have an input node")
	}
	if tree.Output == nil {
		t.Fatal("Expected node tree to have an output node")
	}
}
