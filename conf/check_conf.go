package conf

// NewCheckConf ...
func NewCheckConf(name string, dec *specDec) (*CheckConf, error) {
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
	fc := CheckConf{
		name:     name,
		ins:      []string{in},
		outs:     []string{out},
		assign:   assign,
		required: required,
	}
	return &fc, nil
}

// CheckConf ...
type CheckConf struct {
	name     string
	ins      []string
	outs     []string
	assign   map[string]string
	required map[string]RequiredType
}

// Type ...
func (fc *CheckConf) Type() Type {
	return TypeCheck
}

// Name ...
func (fc *CheckConf) Name() string {
	return fc.name
}

// Inputs ...
func (fc *CheckConf) Inputs() []string {
	return fc.ins
}

// Outputs ...
func (fc *CheckConf) Outputs() []string {
	return fc.outs
}

// Assign ...
func (fc *CheckConf) Assign() map[string]string {
	return fc.assign
}

// Required ...
func (fc *CheckConf) Required() map[string]RequiredType {
	return fc.required
}
