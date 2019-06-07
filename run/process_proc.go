package run

// import (
// 	"github.com/jaqmol/approx/conf"
// 	"github.com/jaqmol/approx/proc"
// )

// // NewProcessProc ...
// func NewProcessProc(conf *conf.ProcessConf) (*ProcessProc, error) {
// 	return &ProcessProc{
// 		conf: conf,
// 		ins:  make(map[string]proc.Conn),
// 		outs: make(map[string]proc.Conn),
// 	}, nil
// }

// // ProcessProc ...
// type ProcessProc struct {
// 	conf *conf.ProcessConf
// 	ins  map[string]proc.Conn
// 	outs map[string]proc.Conn
// }

// // Type ...
// func (pp *ProcessProc) Type() conf.Type {
// 	return pp.conf.Type()
// }

// // Conf ...
// func (pp *ProcessProc) Conf() conf.Conf {
// 	return pp.conf
// }

// // In ...
// func (pp *ProcessProc) In(name string) proc.Conn {
// 	return pp.ins[name]
// }

// // Out ...
// func (pp *ProcessProc) Out(name string) proc.Conn {
// 	return pp.outs[name]
// }

// // AddIn ...
// func (pp *ProcessProc) AddIn(name string, c proc.Conn) {
// 	pp.ins[name] = c
// }

// // AddOut ...
// func (pp *ProcessProc) AddOut(name string, c proc.Conn) {
// 	pp.outs[name] = c
// }

// // Start ...
// func (pp *ProcessProc) Start(errChan chan<- error) {

// }
