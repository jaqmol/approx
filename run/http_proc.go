package run

import "github.com/jaqmol/approx/conf"

// NewHTTPProc ...
func NewHTTPProc(conf *conf.HTTPConf) (*HTTPProc, error) {
	return &HTTPProc{
		conf: conf,
	}, nil
}

// HTTPProc ...
type HTTPProc struct {
	conf *conf.HTTPConf
	ins  []*Conn
	outs []*Conn
}

// Type ...
func (hp *HTTPProc) Type() conf.Type {
	return hp.conf.Type()
}

// // Name ...
// func (hp *HTTPProc) Name() string {
// 	return hp.conf.Name()
// }

// Conf ...
func (hp *HTTPProc) Conf() conf.Conf {
	return hp.conf
}

// Ins ...
func (hp *HTTPProc) Ins() []*Conn {
	return hp.ins
}

// Outs ...
func (hp *HTTPProc) Outs() []*Conn {
	return hp.outs
}

// AddIn ...
func (hp *HTTPProc) AddIn(c *Conn) {
	hp.ins = append(hp.ins, c)
}

// AddOut ...
func (hp *HTTPProc) AddOut(c *Conn) {
	hp.outs = append(hp.outs, c)
}
