package main

// NewForkProc ...
func NewForkProc(name string, dec *SpecDecoder) *ForkProc {
	return &ForkProc{
		name:       &name,
		distribute: dec.Int("distribute"),
		in:         dec.String("in"),
		outs:       dec.StringSlice("outs"),
		props:      dec.StringSlice("props"),
	}
}

// ForkProc ...
type ForkProc struct {
	name       *string
	distribute *int
	in         *string
	outs       []string
	props      []string
}

// Constants ...
const (
	ForkDistributeCopy int = iota
	ForkDistributeRoundRobin
)

// RequiredProps ...
func (p ForkProc) RequiredProps() []string {
	acc := make([]string, 0)
	if p.distribute == nil {
		acc = append(acc, "distribute")
	}
	if p.outs == nil || len(p.outs) == 0 {
		acc = append(acc, "outs")
	}
	if p.props != nil {
		acc = append(acc, p.props...)
	}
	return acc
}

// Outputs ...
func (p ForkProc) Outputs() []string {
	return p.outs
}

// Inputs ...
func (p ForkProc) Inputs() []string {
	if p.in == nil {
		return []string{"stdin"}
	}
	return []string{*p.in}
}

// Name ...
func (p ForkProc) Name() *string {
	return p.name
}
