package run

import (
	"io"
)

// Conn ...
type Conn struct {
	FromProc Proc
	ToProc   Proc
	From     io.Reader
	To       io.Writer
}

// NewConn ...
func NewConn(fromProc Proc, toProc Proc) *Conn {
	reader, writer := io.Pipe()
	return &Conn{
		FromProc: fromProc,
		ToProc:   toProc,
		From:     reader,
		To:       writer,
	}
}

// func makePipes(fo *conf.Formation) {
// 	inputs, outputs := collectInputsAndOutputs(fo)
// 	log.Printf("inputs: %v\n", inputs)
// 	log.Printf("outputs: %v\n", outputs)
// }

// func collectInputsAndOutputs(fo *conf.Formation) (
// 	inputs map[string]bool, outputs map[string]bool,
// ) {
// 	inputs, outputs = make(map[string]bool), make(map[string]bool)
// 	appendInputsAndOutputs(fo.PrivateConfs, inputs, outputs)
// 	appendInputsAndOutputs(fo.PublicConfs, inputs, outputs)
// 	return
// }

// func appendInputsAndOutputs(confs []conf.Conf, inputs map[string]bool, outputs map[string]bool) {
// 	for _, co := range confs {
// 		for _, i := range co.Inputs() {
// 			inputs[i] = true
// 		}
// 		for _, i := range co.Outputs() {
// 			outputs[i] = true
// 		}
// 	}
// }
