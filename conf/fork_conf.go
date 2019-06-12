package conf

import (
	"fmt"
)

// NewForkConf ...
func NewForkConf(name string, dec *specDec) (*ForkConf, error) {
	environment := make([]string, 0)
	distributeStr, ok := dec.string("distribute")
	var distribute ForkDistribute
	if ok {
		switch distributeStr {
		case "copy":
			distribute = ForkDistributeCopy
			environment = append(environment, "DISTRIBUTE=copy")
		case "round_robin":
			distribute = ForkDistributeRoundRobin
			environment = append(environment, "DISTRIBUTE=round_robin")
		}
	} else {
		distribute = ForkDistributeCopy
	}
	in, ok := dec.string("in")
	if !ok {
		in = "stdin"
	}
	outs, ok := dec.stringSlice("outs")
	if !ok {
		return nil, fmt.Errorf("Please provide outputs for fork \"%v\"", name)
	}
	assign, ok := dec.stringStringMap("assign")
	required := make(map[string]RequiredType)
	if ok {
		addAssignmentsToRequired(assign, required)
	} else {
		assign = map[string]string{}
	}

	fc := ForkConf{
		name:        name,
		distribute:  distribute,
		ins:         []string{in},
		outs:        outs,
		assign:      assign,
		required:    required,
		environment: environment,
	}
	return &fc, nil
}

// ForkConf ...
type ForkConf struct {
	name        string
	distribute  ForkDistribute
	ins         []string
	outs        []string
	assign      map[string]string
	required    map[string]RequiredType
	environment []string
}

// ForkDistribute ...
type ForkDistribute int

// ForkDistributes
const (
	ForkDistributeCopy ForkDistribute = iota
	ForkDistributeRoundRobin
)

// Type ...
func (fc *ForkConf) Type() Type {
	return TypeFork
}

// Name ...
func (fc *ForkConf) Name() string {
	return fc.name
}

// Distribute ...
func (fc *ForkConf) Distribute() ForkDistribute {
	return fc.distribute
}

// Inputs ...
func (fc *ForkConf) Inputs() []string {
	return fc.ins
}

// Outputs ...
func (fc *ForkConf) Outputs() []string {
	return fc.outs
}

// Assign ...
func (fc *ForkConf) Assign() map[string]string {
	return fc.assign
}

// Required ...
func (fc *ForkConf) Required() map[string]RequiredType {
	return fc.required
}

// Environment ...
func (fc *ForkConf) Environment() []string {
	return fc.environment
}
