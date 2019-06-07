package run

// import (
// 	"io"

// 	"github.com/jaqmol/approx/proc"
// )

// // PipeConn ...
// type PipeConn struct {
// 	fromProc proc.Proc
// 	out      io.Writer
// 	toProc   proc.Proc
// 	in       io.Reader
// 	name     string
// }

// // NewPipeConn ...
// func NewPipeConn(fromProc proc.Proc, toProc proc.Proc) (proc.Conn, error) {
// 	reader, writer := io.Pipe()
// 	name := proc.ConnName(fromProc, toProc)
// 	return &PipeConn{
// 		fromProc: fromProc,
// 		out:      writer,
// 		toProc:   toProc,
// 		in:       reader,
// 		name:     name,
// 	}, nil
// }

// // FromProc ...
// func (pc *PipeConn) FromProc() proc.Proc {
// 	return pc.fromProc
// }

// // Out ...
// func (pc *PipeConn) Out() (io.Writer, error) {
// 	return pc.out, nil
// }

// // ToProc ...
// func (pc *PipeConn) ToProc() proc.Proc {
// 	return pc.toProc
// }

// // In ...
// func (pc *PipeConn) In() (io.Reader, error) {
// 	return pc.in, nil
// }

// // Name ...
// func (pc *PipeConn) Name() string {
// 	return pc.name
// }

// // Type ...
// func (pc *PipeConn) Type() proc.ConnType {
// 	return proc.ConnTypePipe
// }
