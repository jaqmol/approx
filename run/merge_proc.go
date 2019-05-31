package run

import "github.com/jaqmol/approx/conf"

// NewMergeProc ...
func NewMergeProc(conf *conf.MergeConf) (*MergeProc, error) {
	return &MergeProc{
		conf: conf,
	}, nil
}

// MergeProc ...
type MergeProc struct {
	conf *conf.MergeConf
	ins  []*Conn
	outs []*Conn
}

// Type ...
func (mp *MergeProc) Type() conf.Type {
	return mp.conf.Type()
}

// // Name ...
// func (mp *MergeProc) Name() string {
// 	return mp.conf.Name()
// }

// Conf ...
func (mp *MergeProc) Conf() conf.Conf {
	return mp.conf
}

// Ins ...
func (mp *MergeProc) Ins() []*Conn {
	return mp.ins
}

// Outs ...
func (mp *MergeProc) Outs() []*Conn {
	return mp.outs
}

// AddIn ...
func (mp *MergeProc) AddIn(c *Conn) {
	mp.ins = append(mp.ins, c)
}

// AddOut ...
func (mp *MergeProc) AddOut(c *Conn) {
	mp.outs = append(mp.outs, c)
}
