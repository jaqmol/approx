package run

import (
	"github.com/jaqmol/approx/conf"
)

// Proc ...
type Proc interface {
	Type() conf.Type
	Conf() conf.Conf
	Ins() []*Conn
	Outs() []*Conn
	AddIn(*Conn)
	AddOut(*Conn)
}
