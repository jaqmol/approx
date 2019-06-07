package conf

// NewHTTPServerConf ...
func NewHTTPServerConf(name string, dec *specDec) (*HTTPServerConf, error) {
	required := make(map[string]RequiredType)
	endpoint, ok := dec.string("endpoint")
	if !ok {
		required["endpoint"] = RequiredTypeProperty
		endpoint = "/"
	}
	port, ok := dec.integer("port")
	if !ok {
		port = 3000
	}
	out, ok := dec.string("out")
	if !ok {
		out = "stdout"
	}
	in, ok := dec.string("in")
	if !ok {
		in = "stdin"
	}
	assign, ok := dec.stringStringMap("assign")
	if ok {
		addAssignmentsToRequired(assign, required)
	} else {
		assign = map[string]string{}
	}
	hc := HTTPServerConf{
		name:     name,
		endpoint: endpoint,
		port:     port,
		outs:     []string{out},
		ins:      []string{in},
		assign:   assign,
		required: required,
	}
	return &hc, nil
}

// HTTPServerConf ...
type HTTPServerConf struct {
	name     string
	endpoint string
	port     int
	outs     []string
	ins      []string
	assign   map[string]string
	required map[string]RequiredType
}

// Type ...
func (hc *HTTPServerConf) Type() Type {
	return TypeHTTPServer
}

// Name ...
func (hc *HTTPServerConf) Name() string {
	return hc.name
}

// Endpoint ...
func (hc *HTTPServerConf) Endpoint() string {
	return hc.endpoint
}

// Port ...
func (hc *HTTPServerConf) Port() int {
	return hc.port
}

// Outputs ...
func (hc *HTTPServerConf) Outputs() []string {
	return hc.outs
}

// Inputs ...
func (hc *HTTPServerConf) Inputs() []string {
	return hc.ins
}

// Assign ...
func (hc *HTTPServerConf) Assign() map[string]string {
	return hc.assign
}

// Required ...
func (hc *HTTPServerConf) Required() map[string]RequiredType {
	return hc.required
}
