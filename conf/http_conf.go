package conf

// NewHTTPConf ...
func NewHTTPConf(name string, dec *specDec) (*HTTPConf, error) {
	required := make(map[string]RequiredType)
	endpoint, ok := dec.string("endpoint")
	if !ok {
		required["endpoint"] = RequiredTypeProperty
		endpoint = ""
	}
	out, ok := dec.string("proxy_out")
	if !ok {
		out = "stdout"
	}
	in, ok := dec.string("proxy_in")
	if !ok {
		in = "stdin"
	}
	assign, ok := dec.stringStringMap("assign")
	if ok {
		addAssignmentsToRequired(assign, required)
	} else {
		assign = map[string]string{}
	}
	hc := HTTPConf{
		name:     name,
		endpoint: endpoint,
		outs:     []string{out},
		ins:      []string{in},
		assign:   assign,
		required: required,
	}
	return &hc, nil
}

// HTTPConf ...
type HTTPConf struct {
	name     string
	endpoint string
	outs     []string
	ins      []string
	assign   map[string]string
	required map[string]RequiredType
}

// Type ...
func (hc *HTTPConf) Type() Type {
	return TypeHTTP
}

// Name ...
func (hc *HTTPConf) Name() string {
	return hc.name
}

// Endpoint ...
func (hc *HTTPConf) Endpoint() string {
	return hc.endpoint
}

// Outputs ...
func (hc *HTTPConf) Outputs() []string {
	return hc.outs
}

// Inputs ...
func (hc *HTTPConf) Inputs() []string {
	return hc.ins
}

// Assign ...
func (hc *HTTPConf) Assign() map[string]string {
	return hc.assign
}

// Required ...
func (hc *HTTPConf) Required() map[string]RequiredType {
	return hc.required
}
