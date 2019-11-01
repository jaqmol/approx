package configuration

import (
	"fmt"

	"github.com/jaqmol/approx/project"
)

// TODO: Write Test

// Formation ...
type Formation struct {
	Processors map[string]Processor
	FlowTree   *FlowTree
}

// NewFormation ...
func NewFormation(projForm *project.Formation) (*Formation, error) {
	acc := make(map[string]Processor, len(projForm.Definitions))
	for name, def := range projForm.Definitions {
		switch def.Type() {
		case project.CommandType:
			prCmd := def.(*project.Command)
			acc[name] = &Command{
				Ident: prCmd.Ident(),
				Cmd:   prCmd.Cmd(),
				Env:   joinKeyValues(prCmd.Env()),
			}
		case project.ForkType:
			acc[name] = &Fork{
				Ident: def.Ident(),
			}
		case project.MergeType:
			acc[name] = &Merge{
				Ident: def.Ident(),
			}
		}
	}
	flow, err := NewFlowTree(projForm.Flows, acc)
	if err != nil {
		return nil, err
	}
	return &Formation{acc, flow}, nil
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
