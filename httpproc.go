package main

// NewHTTPProc ...
func NewHTTPProc(name string, dec *SpecDecoder) *HTTPProc {
	return &HTTPProc{
		name:     &name,
		endpoint: dec.String("endpoint"),
		proxyOut: dec.String("proxy_out"),
		proxyIn:  dec.String("proxy_in"),
		props:    dec.StringSlice("props"),
	}
}

// HTTPProc ...
type HTTPProc struct {
	name     *string
	endpoint *string
	proxyOut *string
	proxyIn  *string
	props    []string
}

// RequiredProps ...
func (p HTTPProc) RequiredProps() []string {
	acc := make([]string, 0)
	if p.endpoint == nil {
		acc = append(acc, "endpoint")
	}
	if p.props != nil {
		acc = append(acc, p.props...)
	}
	return acc
}

// Outputs ...
func (p HTTPProc) Outputs() []string {
	return []string{*p.proxyOut}
}

// Inputs ...
func (p HTTPProc) Inputs() []string {
	return []string{*p.proxyIn}
}

// Name ...
func (p HTTPProc) Name() *string {
	return p.name
}
