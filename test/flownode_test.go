package test

import (
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

	forkNode.Iterate(func(prev []*configuration.FlowNode, curr *configuration.FlowNode) {
		visited[curr.Processor().ID()] = true
		if nodesEqual(curr, forkNode) {
			if len(prev) != 0 {
				t.Fatal("Expected fork node to have 0 predecessors")
			}
		} else if nodesEqual(curr, fneNode) {
			if len(prev) != 1 {
				t.Fatal("Expected fne node to have 1 predecessors")
			}
		} else if nodesEqual(curr, lneNode) {
			if len(prev) != 1 {
				t.Fatal("Expected lne node to have 1 predecessors")
			}
		} else if nodesEqual(curr, mergeNode) {
			if len(prev) != 2 {
				t.Fatal("Expected merge node to have 2 predecessors")
			}
		}
	})

	if len(visited) != 4 {
		t.Fatal("Expected to visit 4 nodes, but got:", len(visited))
	}
	if _, ok := visited[forkNode.Processor().ID()]; !ok {
		t.Fatal("Expected to visit fork nodes, but didn't")
	}
	if _, ok := visited[fneNode.Processor().ID()]; !ok {
		t.Fatal("Expected to visit fne nodes, but didn't")
	}
	if _, ok := visited[lneNode.Processor().ID()]; !ok {
		t.Fatal("Expected to visit lne nodes, but didn't")
	}
	if _, ok := visited[mergeNode.Processor().ID()]; !ok {
		t.Fatal("Expected to visit merge nodes, but didn't")
	}
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
