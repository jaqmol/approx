package configuration

import (
	"fmt"
	"strings"

	"github.com/jaqmol/approx/project"
)

// TODO: Implement test

// FlowNode ...
type FlowNode struct {
	previous  []FlowNode
	next      []FlowNode
	processor Processor
}

// NewFlowTree ...
func NewFlowTree(flows []project.Flow, procs map[string]Processor) (*FlowNode, error) {
	nodeForName := make(map[string]FlowNode)
	for _, line := range flows {
		for j, toName := range line {
			if j > 0 {
				fromName := line[j-1]

				fromNode := getCreateNode(fromName, procs[fromName], nodeForName)
				toNode := getCreateNode(toName, procs[toName], nodeForName)

				fromNode.next = append(fromNode.next, *toNode)
				toNode.previous = append(toNode.previous, *fromNode)
			} else {
				getCreateNode(toName, procs[toName], nodeForName)
			}
		}
	}

	for _, node := range nodeForName {
		node.previous = makeUniqueSet(node.previous)
		node.next = makeUniqueSet(node.next)
	}

	return findNodeWithoutPredecessor(nodeForName)
}

// // Previous ...
// func (fn *FlowNode) Previous() []FlowNode {
// 	return fn.previous
// }

// // Next ...
// func (fn *FlowNode) Next() []FlowNode {
// 	return fn.next
// }

// Processor ...
func (fn *FlowNode) Processor() Processor {
	return fn.processor
}

// // ID ...
// func (fn *FlowNode) ID() string {
// 	return fn.processor.ID()
// }

// Iterate ...
func (fn *FlowNode) Iterate(callback func(prev []FlowNode, curr *FlowNode)) {
	callback(fn.previous, fn)
	for _, next := range fn.next {
		next.Iterate(callback)
	}
}

func findNodeWithoutPredecessor(nodeForName map[string]FlowNode) (*FlowNode, error) {
	inputNodes := make([]FlowNode, 0)
	for _, node := range nodeForName {
		if len(node.previous) == 0 {
			inputNodes = append(inputNodes, node)
		}
	}
	if len(inputNodes) == 1 {
		node := inputNodes[0]
		return &node, nil
	}
	ids := make([]string, len(inputNodes))
	for i, n := range inputNodes {
		ids[i] = n.processor.ID()
	}
	err := fmt.Errorf("Exactly one input node expected, but found %v: %v", len(inputNodes), strings.Join(ids, ", "))
	return nil, err
}

func getCreateNode(name string, proc Processor, acc map[string]FlowNode) *FlowNode {
	node, ok := acc[name]
	if !ok {
		node = FlowNode{
			previous:  make([]FlowNode, 0),
			next:      make([]FlowNode, 0),
			processor: proc,
		}
		acc[name] = node
	}
	return &node
}

func makeUniqueSet(input []FlowNode) []FlowNode {
	output := make([]FlowNode, 0)
	isContainedForID := make(map[string]bool)
	for _, node := range input {
		id := node.processor.ID()
		_, isOk := isContainedForID[id]
		if !isOk {
			isContainedForID[id] = true
			output = append(output, node)
		}
	}
	return output
}
