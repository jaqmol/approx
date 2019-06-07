package run

// import (
// 	"fmt"

// 	"github.com/jaqmol/approx/conf"
// 	"github.com/jaqmol/approx/proc"
// )

// func newHubBuilder() *hubBuilder {
// 	return &hubBuilder{
// 		procForName: make(map[string]proc.Proc),
// 		connForName: make(map[string]proc.Conn),
// 	}
// }

// type hubBuilder struct {
// 	procForName map[string]proc.Proc
// 	connForName map[string]proc.Conn
// }

// func (hb *hubBuilder) publicProcs(fo *conf.Formation) (procs []proc.Proc) {
// 	procs = make([]proc.Proc, 0)

// 	for _, co := range fo.PublicConfs {
// 		procs = append(procs, hb.procForName[co.Name()])
// 	}
// 	return
// }

// func (hb *hubBuilder) connectAllProcs() (err error) {
// 	for procName, sourceProc := range hb.procForName {
// 		for _, destinationName := range sourceProc.Conf().Outputs() {
// 			if destinationName == "stdout" {
// 				continue
// 			}
// 			destinationProc, ok := hb.procForName[destinationName]
// 			if !ok {
// 				err = fmt.Errorf("Cannot connect \"%v\" to in of \"%v\", \"%v\" not found", procName, destinationName, destinationName)
// 				return
// 			}
// 			hb.connectSourceAndDestinationProcs(sourceProc, destinationProc)
// 		}
// 		destinationProc := sourceProc
// 		for _, sourceName := range destinationProc.Conf().Inputs() {
// 			if sourceName == "stdin" {
// 				continue
// 			}
// 			sourceProc, ok := hb.procForName[sourceName]
// 			if !ok {
// 				err = fmt.Errorf("Cannot connect \"%v\" to out of \"%v\", \"%v\" not found", procName, sourceName, sourceName)
// 				return
// 			}
// 			hb.connectSourceAndDestinationProcs(sourceProc, destinationProc)
// 		}
// 	}
// 	return
// }

// func (hb *hubBuilder) connectSourceAndDestinationProcs(sourceProc proc.Proc, destinationProc proc.Proc) error {
// 	connName := proc.ConnName(sourceProc, destinationProc)
// 	conn, ok := hb.connForName[connName]
// 	if !ok {
// 		conn, err := newConn(sourceProc, destinationProc)
// 		if err != nil {
// 			return err
// 		}
// 		// conn = proc.NewConn(sourceProc, destinationProc)
// 		hb.connForName[connName] = conn
// 	}
// 	sourceProc.AddOut(conn.Name(), conn)
// 	destinationProc.AddIn(conn.Name(), conn)
// 	return nil
// }

// func (hb *hubBuilder) initAllProcs(fo *conf.Formation) (err error) {
// 	err = hb.initSomeProcs(fo.PublicConfs)
// 	if err != nil {
// 		return
// 	}
// 	err = hb.initSomeProcs(fo.PrivateConfs)
// 	return
// }

// func (hb *hubBuilder) initSomeProcs(confs []conf.Conf) (err error) {
// 	var proc proc.Proc
// 	for _, co := range confs {
// 		name := co.Name()
// 		_, exists := hb.procForName[name]
// 		if exists {
// 			err = fmt.Errorf("Processor specification \"%v\" is not unique", name)
// 			return
// 		}
// 		proc, err = initOneProc(co)
// 		if err != nil {
// 			return
// 		}
// 		hb.procForName[name] = proc
// 	}
// 	return
// }

// func initOneProc(co conf.Conf) (proc proc.Proc, err error) {
// 	switch co.Type() {
// 	case conf.TypeProcess:
// 		proc, err = NewProcessProc(co.(*conf.ProcessConf))
// 	case conf.TypeHTTPServer:
// 		proc, err = NewHTTPServerProc(co.(*conf.HTTPServerConf))
// 	case conf.TypeFork:
// 		proc, err = NewForkProc(co.(*conf.ForkConf))
// 	case conf.TypeMerge:
// 		proc, err = NewMergeProc(co.(*conf.MergeConf))
// 	}
// 	return
// }
