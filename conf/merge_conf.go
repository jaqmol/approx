package conf

import "fmt"

// NewMergeConf ...
func NewMergeConf(name string, dec *specDec) (*MergeConf, error) {
	pickStr, ok := dec.string("pick")
	var pick MergePick
	if ok {
		switch pickStr {
		case "as_comes":
			pick = MergePickAsComes
		case "round_robin":
			pick = MergePickRoundRobin
		}
	} else {
		pick = MergePickAsComes
	}
	ins, ok := dec.stringSlice("ins")
	if !ok {
		return nil, fmt.Errorf("Please provide inputs for fork \"%v\"", name)
	}
	out, ok := dec.string("out")
	if !ok {
		out = "stdout"
	}
	assign, ok := dec.stringStringMap("assign")
	required := make(map[string]RequiredType)
	if ok {
		addAssignmentsToRequired(assign, required)
	} else {
		assign = map[string]string{}
	}
	mc := MergeConf{
		name:     name,
		pick:     pick,
		ins:      ins,
		outs:     []string{out},
		assign:   assign,
		required: map[string]RequiredType{},
	}
	return &mc, nil
}

// MergeConf ...
type MergeConf struct {
	name     string
	pick     MergePick
	ins      []string
	outs     []string
	assign   map[string]string
	required map[string]RequiredType
}

// MergePick ...
type MergePick int

// MergePicks
const (
	MergePickAsComes MergePick = iota
	MergePickRoundRobin
)

// Type ...
func (mc *MergeConf) Type() Type {
	return TypeMerge
}

// Name ...
func (mc *MergeConf) Name() string {
	return mc.name
}

// Pick ...
func (mc *MergeConf) Pick() MergePick {
	return mc.pick
}

// Outputs ...
func (mc *MergeConf) Outputs() []string {
	return mc.outs
}

// Inputs ...
func (mc *MergeConf) Inputs() []string {
	return mc.ins
}

// Assign ...
func (mc *MergeConf) Assign() map[string]string {
	return mc.assign
}

// Required ...
func (mc *MergeConf) Required() map[string]RequiredType {
	return mc.required
}
