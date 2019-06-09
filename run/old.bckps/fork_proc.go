package run

// import (
// 	"github.com/jaqmol/approx/conf"
// )

// // NewForkProc ...
// func NewForkProc(conf *conf.ForkConf) (*ForkProc, error) {
// 	return &ForkProc{
// 		conf: conf,
// 		ins:  make(map[string]*flow.Conn),
// 		outs: make(map[string]*flow.Conn),
// 	}, nil
// }

// // ForkProc ...
// type ForkProc struct {
// 	conf *conf.ForkConf
// 	ins  map[string]*flow.Conn
// 	outs map[string]*flow.Conn
// }

// // Type ...
// func (fp *ForkProc) Type() conf.Type {
// 	return fp.conf.Type()
// }

// // Conf ...
// func (fp *ForkProc) Conf() conf.Conf {
// 	return fp.conf
// }

// // In ...
// func (fp *ForkProc) In(name string) *flow.Conn {
// 	return fp.ins[name]
// }

// // Out ...
// func (fp *ForkProc) Out(name string) *flow.Conn {
// 	return fp.outs[name]
// }

// // AddIn ...
// func (fp *ForkProc) AddIn(name string, c *flow.Conn) {
// 	fp.ins[name] = c
// }

// // AddOut ...
// func (fp *ForkProc) AddOut(name string, c *flow.Conn) {
// 	fp.outs[name] = c
// }

// // Start ...
// func (fp *ForkProc) Start(errChan chan<- error) {
// }
