package conf

import "fmt"

// NewProcessConf ...
func NewProcessConf(name string, dec *specDec) (*ProcessConf, error) {
	environment, ok := dec.stringStringMap("environment")
	if !ok {
		environment = map[string]string{}
	}
	arguments, ok := dec.stringSlice("arguments")
	if !ok {
		arguments = []string{}
	}
	command, ok := dec.string("command")
	if !ok {
		return nil, fmt.Errorf("Please provide a command for process \"%v\"", name)
	}
	in, ok := dec.string("in")
	if !ok {
		in = "stdin"
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
	return &ProcessConf{
		name:        name,
		environment: environment,
		arguments:   arguments,
		command:     command,
		ins:         []string{in},
		outs:        []string{out},
		assign:      assign,
		required:    required,
	}, nil
}

// ProcessConf ...
type ProcessConf struct {
	name        string
	environment map[string]string
	arguments   []string
	command     string
	ins         []string
	outs        []string
	assign      map[string]string
	required    map[string]RequiredType
}

// Type ...
func (pc *ProcessConf) Type() Type {
	return TypeProcess
}

// Name ...
func (pc *ProcessConf) Name() string {
	return pc.name
}

// Environment ...
func (pc *ProcessConf) Environment() map[string]string {
	return pc.environment
}

// Arguments ...
func (pc *ProcessConf) Arguments() []string {
	return pc.arguments
}

// Command ...
func (pc *ProcessConf) Command() string {
	return pc.command
}

// Inputs ...
func (pc *ProcessConf) Inputs() []string {
	return pc.ins
}

// Outputs ...
func (pc *ProcessConf) Outputs() []string {
	return pc.outs
}

// Assign ...
func (pc *ProcessConf) Assign() map[string]string {
	return pc.assign
}

// Required ...
func (pc *ProcessConf) Required() map[string]RequiredType {
	return pc.required
}
