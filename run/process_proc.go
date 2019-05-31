package run

import (
	"github.com/jaqmol/approx/conf"
)

// NewProcessProc ...
func NewProcessProc(conf *conf.ProcessConf) (*ProcessProc, error) {
	return &ProcessProc{
		conf: conf,
	}, nil
}

// ProcessProc ...
type ProcessProc struct {
	conf *conf.ProcessConf
	ins  []*Conn
	outs []*Conn
}

// Type ...
func (pp *ProcessProc) Type() conf.Type {
	return pp.conf.Type()
}

// // Name ...
// func (pp *ProcessProc) Name() string {
// 	return pp.conf.Name()
// }

// Conf ...
func (pp *ProcessProc) Conf() conf.Conf {
	return pp.conf
}

// Ins ...
func (pp *ProcessProc) Ins() []*Conn {
	return pp.ins
}

// Outs ...
func (pp *ProcessProc) Outs() []*Conn {
	return pp.outs
}

// AddIn ...
func (pp *ProcessProc) AddIn(c *Conn) {
	pp.ins = append(pp.ins, c)
}

// AddOut ...
func (pp *ProcessProc) AddOut(c *Conn) {
	pp.outs = append(pp.outs, c)
}
