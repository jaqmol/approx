package configuration

import (
	"fmt"

	"github.com/jaqmol/approx/project"
)

// Formation ...
type Formation struct {
	Processors map[string]Processor
	FlowTree   *FlowTree
}

// NewFormation ...
func NewFormation(projForm *project.Formation) (*Formation, error) {
	procs := make(map[string]Processor, len(projForm.Definitions))
	for name, def := range projForm.Definitions {
		switch def.Type() {
		case project.StdinType:
			procs[name] = &Stdin
		case project.CommandType:
			prCmd := def.(*project.Command)
			procs[name] = &Command{
				Ident: prCmd.Ident(),
				Cmd:   prCmd.Cmd(),
				Env:   joinKeyValues(prCmd.Env()),
			}
		case project.ForkType:
			procs[name] = &Fork{
				Ident: def.Ident(),
			}
		case project.MergeType:
			procs[name] = &Merge{
				Ident: def.Ident(),
			}
		case project.StdoutType:
			procs[name] = &Stdout
		}
	}
	tree, err := NewFlowTree(projForm.Flows, procs)
	if err != nil {
		return nil, err
	}
	tree.Iterate(func(prev []*FlowNode, curr *FlowNode, next []*FlowNode) error {
		if curr.Processor().Type() == ForkType {
			fork := curr.Processor().(*Fork)
			fork.Count = len(next)
		} else if curr.Processor().Type() == MergeType {
			merge := curr.Processor().(*Merge)
			merge.Count = len(prev)
		}
		return nil
	})
	return &Formation{procs, tree}, nil
}

func joinKeyValues(mapping map[string]string) []string {
	acc := make([]string, len(mapping))
	idx := 0
	for key, value := range mapping {
		joining := fmt.Sprintf("%v=%v", key, value)
		acc[idx] = joining
		idx++
	}
	return acc
}
