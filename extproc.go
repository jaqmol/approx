package main

import "fmt"

// NewExtProc ...
func NewExtProc(name string, dec *SpecDecoder) (*ExtProc, error) {
	cmd := dec.String("command")
	if cmd == nil {
		return nil, fmt.Errorf("Please provide a command for process \"%v\"", name)
	}
	return &ExtProc{
		name:        &name,
		environment: dec.StringStringMap("environment"),
		arguments:   dec.StringSlice("arguments"),
		command:     cmd,
		in:          dec.String("in"),
		out:         dec.String("out"),
		props:       dec.StringSlice("props"),
	}, nil
}

// ExtProc ...
type ExtProc struct {
	name        *string
	environment map[string]string
	arguments   []string
	command     *string
	in          *string
	out         *string
	props       []string
}

// RequiredProps ...
func (p ExtProc) RequiredProps() []string {
	return p.props
}

// Outputs ...
func (p ExtProc) Outputs() []string {
	if p.out == nil {
		return []string{"stdout"}
	}
	return []string{*p.out}
}

// Inputs ...
func (p ExtProc) Inputs() []string {
	if p.in == nil {
		return []string{"stdin"}
	}
	return []string{*p.in}
}

// Name ...
func (p ExtProc) Name() *string {
	return p.name
}
