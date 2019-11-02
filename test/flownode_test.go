package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jaqmol/approx/configuration"
)

// TestFlowNodes ...
func TestFlowNodes(t *testing.T) {
	forkNode, fneNode, lneNode, mergeNode := createTestFlow()

	checkNodeNextCount(t, forkNode, 2, "forkNode")
	checkNodePreviousCount(t, fneNode, 1, "fneNode")
	checkNodeNextCount(t, fneNode, 1, "fneNode")
	checkNodePreviousCount(t, lneNode, 1, "lneNode")
	checkNodeNextCount(t, lneNode, 1, "lneNode")
	checkNodePreviousCount(t, mergeNode, 2, "mergeNode")

	visited := make(map[string]bool)

	checkLen := lengthChecker(map[string]int{
		forkNode.Processor().ID():  0,
		fneNode.Processor().ID():   1,
		lneNode.Processor().ID():   1,
		mergeNode.Processor().ID(): 2,
	})

	forkNode.Iterate(func(prev []*configuration.FlowNode, curr *configuration.FlowNode) {
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
		forkNode.Processor().ID(),
		fneNode.Processor().ID(),
		lneNode.Processor().ID(),
		mergeNode.Processor().ID(),
	)
	if len(errs) > 0 {
		err := fmt.Errorf("Errors visiting nodes: %v", strings.Join(errorsToStrings(errs), ", "))
		t.Fatal(err)
	}
}

func lengthChecker(expected map[string]int) func(string, int) error {
	return func(id string, givenLen int) error {
		expectedLen := expected[id]
		if expectedLen != givenLen {
			return fmt.Errorf("Expected %v node to have %v predecessors, but got %v", id, expectedLen, givenLen)
		}
		return nil
	}
}

func checkContainsAll(checkIn map[string]bool, checkForAll ...string) []error {
	acc := make([]error, 0)
	for _, checkFor := range checkForAll {
		if _, ok := checkIn[checkFor]; !ok {
			err := fmt.Errorf("Expected to visit %v node, but didn't", checkFor)
			acc = append(acc, err)
		}
	}
	return acc
}

func errorsToStrings(errs []error) []string {
	acc := make([]string, len(errs))
	for i, e := range errs {
		acc[i] = e.Error()
	}
	return acc
}

func createTestFlow() (
	forkNode *configuration.FlowNode,
	fneNode *configuration.FlowNode,
	lneNode *configuration.FlowNode,
	mergeNode *configuration.FlowNode,
) {
	config := makeSimpleSequenceConfig()

	forkNode = configuration.NewFlowNode(&config.fork)
	fneNode = configuration.NewFlowNode(&config.firstNameExtract)
	lneNode = configuration.NewFlowNode(&config.lastNameExtract)
	mergeNode = configuration.NewFlowNode(&config.merge)

	forkNode.AppendNext(fneNode, lneNode)
	fneNode.AppendPrevious(forkNode)
	fneNode.AppendNext(mergeNode)
	lneNode.AppendPrevious(forkNode)
	lneNode.AppendNext(mergeNode)
	mergeNode.AppendPrevious(fneNode, lneNode)

	return
}

func nodesEqual(a *configuration.FlowNode, b *configuration.FlowNode) bool {
	return a.Processor().ID() == b.Processor().ID()
}

func checkNodePreviousCount(t *testing.T, node *configuration.FlowNode, count int, name string) {
	length := len(node.Previous())
	if length != count {
		t.Fatalf("Expected %v to have %v predecessors, but found: %v", name, count, length)
	}
}

func checkNodeNextCount(t *testing.T, node *configuration.FlowNode, count int, name string) {
	length := len(node.Next())
	if length != count {
		t.Fatalf("Expected %v to have %v successors, but found: %v", name, count, length)
	}
}
