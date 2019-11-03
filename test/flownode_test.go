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

	visited := make(map[string]int)

	checkLen := lengthChecker(map[string][]int{
		forkNode.Processor().ID():  []int{0, 2},
		fneNode.Processor().ID():   []int{1, 1},
		lneNode.Processor().ID():   []int{1, 1},
		mergeNode.Processor().ID(): []int{2, 0},
	})

	err := forkNode.Iterate(func(prev []*configuration.FlowNode, curr *configuration.FlowNode, next []*configuration.FlowNode) error {
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
			forkNode.Processor().ID():  1,
			fneNode.Processor().ID():   1,
			lneNode.Processor().ID():   1,
			mergeNode.Processor().ID(): 2,
		},
	)
	if len(errs) > 0 {
		err := fmt.Errorf("Errors visiting nodes: %v", strings.Join(errorsToStrings(errs), ", "))
		t.Fatal(err)
	}
}

func lengthChecker(expected map[string][]int) func(string, int, int) error {
	return func(id string, givenPrevLen, givenNextLen int) error {
		expectedLens := expected[id]
		expectedPrevLen := expectedLens[0]
		expectedNextLen := expectedLens[1]
		if expectedPrevLen != givenPrevLen {
			return fmt.Errorf("Expected %v node to have %v predecessors, but got %v", id, expectedPrevLen, givenPrevLen)
		}
		if expectedNextLen != givenNextLen {
			return fmt.Errorf("Expected %v node to have %v successors, but got %v", id, expectedNextLen, givenNextLen)
		}
		return nil
	}
}

func checkContainsAllTimes(checkIn map[string]int, checkForTimes map[string]int) []error {
	acc := make([]error, 0)
	for checkFor, expectedTimes := range checkForTimes {
		if givenTimes := checkIn[checkFor]; givenTimes != expectedTimes {
			err := fmt.Errorf("Expected to visit %v node %v times, but got %v", checkFor, expectedTimes, givenTimes)
			acc = append(acc, err)
		}
	}
	return acc
}

// func checkContainsTimes(checkIn map[string]int, checkFor string, expectedTimes int) error {
// 	givenTimes := checkIn[checkFor]
// 	if givenTimes != expectedTimes {
// 		return fmt.Errorf("Expected to visit %v node %v times, but got %v", checkFor, expectedTimes, givenTimes)
// 	}
// 	return nil
// }

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
