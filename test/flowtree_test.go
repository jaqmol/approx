package test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jaqmol/approx/config"
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

	conf := MakeSimpleSequenceConfig()
	procs := map[string]config.Processor{
		conf.Fork.Ident:             &conf.Fork,
		conf.FirstNameExtract.Ident: &conf.FirstNameExtract,
		conf.LastNameExtract.Ident:  &conf.LastNameExtract,
		conf.Merge.Ident:            &conf.Merge,
	}
	tree, err := config.NewFlowTree(flows, procs)
	if err != nil {
		t.Fatal(err)
	}

	visited := make(map[string]int)

	checkLen := lengthChecker(map[string][]int{
		conf.Fork.ID():             []int{0, 2},
		conf.FirstNameExtract.ID(): []int{1, 1},
		conf.LastNameExtract.ID():  []int{1, 1},
		conf.Merge.ID():            []int{2, 0},
	})

	err = tree.Iterate(func(prev []*config.FlowNode, curr *config.FlowNode, next []*config.FlowNode) error {
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
			conf.Fork.ID():             1,
			conf.FirstNameExtract.ID(): 1,
			conf.LastNameExtract.ID():  1,
			conf.Merge.ID():            1,
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
