package main

// NewMergeProc ...
func NewMergeProc(name string, dec *SpecDecoder) *MergeProc {
	return &MergeProc{
		name:  &name,
		pick:  dec.Int("pick"),
		ins:   dec.StringSlice("ins"),
		out:   dec.String("out"),
		props: dec.StringSlice("props"),
	}
}

// MergeProc ...
type MergeProc struct {
	name  *string
	pick  *int
	ins   []string
	out   *string
	props []string
}

// Constants ...
const (
	MergePickAsComes int = iota
	MergePickRoundRobin
)

// RequiredProps ...
func (p MergeProc) RequiredProps() []string {
	acc := make([]string, 0)
	if p.pick == nil {
		acc = append(acc, "pick")
	}
	if p.ins == nil || len(p.ins) == 0 {
		acc = append(acc, "ins")
	}
	if p.props != nil {
		acc = append(acc, p.props...)
	}
	return acc
}

// Outputs ...
func (p MergeProc) Outputs() []string {
	if p.out == nil {
		return []string{"stdout"}
	}
	return []string{*p.out}
}

// Inputs ...
func (p MergeProc) Inputs() []string {
	return p.ins
}

// Name ...
func (p MergeProc) Name() *string {
	return p.name
}
