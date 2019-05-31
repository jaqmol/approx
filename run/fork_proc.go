package run

import (
	"github.com/jaqmol/approx/conf"
)

// NewForkProc ...
func NewForkProc(conf *conf.ForkConf) (*ForkProc, error) {
	return &ForkProc{
		conf: conf,
	}, nil
}

// ForkProc ...
type ForkProc struct {
	conf *conf.ForkConf
	ins  []*Conn
	outs []*Conn
}

// Type ...
func (fp *ForkProc) Type() conf.Type {
	return fp.conf.Type()
}

// // Name ...
// func (fp *ForkProc) Name() string {
// 	return fp.conf.Name()
// }

// Conf ...
func (fp *ForkProc) Conf() conf.Conf {
	return fp.conf
}

// Ins ...
func (fp *ForkProc) Ins() []*Conn {
	return fp.ins
}

// Outs ...
func (fp *ForkProc) Outs() []*Conn {
	return fp.outs
}

// AddIn ...
func (fp *ForkProc) AddIn(c *Conn) {
	fp.ins = append(fp.ins, c)
}

// AddOut ...
func (fp *ForkProc) AddOut(c *Conn) {
	fp.outs = append(fp.outs, c)
}
