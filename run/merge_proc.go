package run

// import (
// 	"github.com/jaqmol/approx/conf"
// 	"github.com/jaqmol/approx/proc"
// )

// // NewMergeProc ...
// func NewMergeProc(conf *conf.MergeConf) (*MergeProc, error) {
// 	return &MergeProc{
// 		conf: conf,
// 		ins:  make(map[string]proc.Conn),
// 		outs: make(map[string]proc.Conn),
// 	}, nil
// }

// // MergeProc ...
// type MergeProc struct {
// 	conf *conf.MergeConf
// 	ins  map[string]proc.Conn
// 	outs map[string]proc.Conn
// }

// // Type ...
// func (mp *MergeProc) Type() conf.Type {
// 	return mp.conf.Type()
// }

// // Conf ...
// func (mp *MergeProc) Conf() conf.Conf {
// 	return mp.conf
// }

// // In ...
// func (mp *MergeProc) In(name string) proc.Conn {
// 	return mp.ins[name]
// }

// // Out ...
// func (mp *MergeProc) Out(name string) proc.Conn {
// 	return mp.outs[name]
// }

// // AddIn ...
// func (mp *MergeProc) AddIn(name string, c proc.Conn) {
// 	mp.ins[name] = c
// }

// // AddOut ...
// func (mp *MergeProc) AddOut(name string, c proc.Conn) {
// 	mp.outs[name] = c
// }

// // Start ...
// func (mp *MergeProc) Start(errChan chan<- error) {

// }
