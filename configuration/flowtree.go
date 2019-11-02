package configuration

import (
	"fmt"
	"strings"

	"github.com/jaqmol/approx/project"
)

// FlowTree ...
type FlowTree struct {
	Input  *FlowNode
	Output *FlowNode
}

// NewFlowTree ...
func NewFlowTree(flows []project.Flow, procs map[string]Processor) (*FlowTree, error) {
	nodeForName := make(map[string]*FlowNode)
	for _, line := range flows {
		for j, toName := range line {
			if j > 0 {
				fromName := line[j-1]
				fromNode := getCreateNode(fromName, procs[fromName], nodeForName)
				toNode := getCreateNode(toName, procs[toName], nodeForName)
				fromNode.AppendNext(toNode)
				toNode.AppendPrevious(fromNode)
			} else {
				getCreateNode(toName, procs[toName], nodeForName)
			}
		}
	}

	for _, node := range nodeForName {
		node.previous = makeUniqueSet(node.previous)
		node.next = makeUniqueSet(node.next)
	}

	input, output, err := findNoPredecessorAndNoSuccessorNodes(nodeForName)
	if err != nil {
		return nil, err
	}
	return &FlowTree{input, output}, nil
}

// Iterate ...
func (ft *FlowTree) Iterate(callback func(prev []*FlowNode, curr *FlowNode)) {
	wasVisitedForID := make(map[string]bool)
	ft.Input.Iterate(func(prev []*FlowNode, curr *FlowNode) {
		id := curr.processor.ID()
		_, ok := wasVisitedForID[id]
		if !ok {
			wasVisitedForID[id] = true
			callback(prev, curr)
		}
	})
}

func findNoPredecessorAndNoSuccessorNodes(nodeForName map[string]*FlowNode) (
	noPredecessor *FlowNode,
	noSuccessor *FlowNode,
	err error,
) {
	inputNodes := make([]*FlowNode, 0)
	outputNodes := make([]*FlowNode, 0)
	for _, node := range nodeForName {
		if len(node.previous) == 0 {
			inputNodes = append(inputNodes, node)
		}
		if len(node.next) == 0 {
			outputNodes = append(outputNodes, node)
		}
	}
	// Happy path
	if len(inputNodes) == 1 {
		noPredecessor = inputNodes[0]
		if len(outputNodes) == 1 {
			noSuccessor = outputNodes[0]
		}
		return
	}
	// Fail path
	inputIds := collectIDs(inputNodes)
	outputIds := collectIDs(outputNodes)

	allErrMsgs := make([]string, 0)
	if len(inputIds) > 0 {
		allErrMsgs = append(allErrMsgs, moreThanOneErrorMsg(inputIds))
	}
	if len(outputIds) > 0 {
		allErrMsgs = append(allErrMsgs, moreThanOneErrorMsg(outputIds))
	}
	err = fmt.Errorf("Error(s) interpreting flow: %v", strings.Join(allErrMsgs, "; "))
	return nil, nil, err
}

func getCreateNode(name string, proc Processor, acc map[string]*FlowNode) *FlowNode {
	node, ok := acc[name]
	if !ok {
		node = NewFlowNode(proc)
		acc[name] = node
	}
	return node
}

func makeUniqueSet(input []*FlowNode) []*FlowNode {
	output := make([]*FlowNode, 0)
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

func collectIDs(nodes []*FlowNode) []string {
	ids := make([]string, len(nodes))
	for i, n := range nodes {
		ids[i] = n.processor.ID()
	}
	return ids
}

func moreThanOneErrorMsg(ids []string) string {
	return fmt.Sprintf("More than 1 input node: %v", strings.Join(ids, ", "))
}
