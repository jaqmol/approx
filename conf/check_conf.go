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
		name:        name,
		ins:         []string{in},
		outs:        []string{out},
		assign:      assign,
		required:    required,
		environment: []string{},
	}
	return &fc, nil
}

// CheckConf ...
type CheckConf struct {
	name        string
	ins         []string
	outs        []string
	assign      map[string]string
	required    map[string]RequiredType
	environment []string
}

// Type ...
func (cc *CheckConf) Type() Type {
	return TypeCheck
}

// Name ...
func (cc *CheckConf) Name() string {
	return cc.name
}

// Inputs ...
func (cc *CheckConf) Inputs() []string {
	return cc.ins
}

// Outputs ...
func (cc *CheckConf) Outputs() []string {
	return cc.outs
}

// Assign ...
func (cc *CheckConf) Assign() map[string]string {
	return cc.assign
}

// Required ...
func (cc *CheckConf) Required() map[string]RequiredType {
	return cc.required
}

// Environment ...
func (cc *CheckConf) Environment() []string {
	return cc.environment
}
