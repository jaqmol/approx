package run

// import (
// 	"io"
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
// 	"syscall"

// 	"github.com/jaqmol/approx/proc"
// )

// // NamedPipeConn ...
// type NamedPipeConn struct {
// 	name     string
// 	Path     string
// 	fromProc proc.Proc
// 	toProc   proc.Proc
// }

// // NewNamedPipeConn ...
// func NewNamedPipeConn(fromProc proc.Proc, toProc proc.Proc) (proc.Conn, error) {
// 	tmpDir, _ := ioutil.TempDir("", "approx")
// 	name := proc.ConnName(fromProc, toProc)
// 	path := filepath.Join(tmpDir, name)
// 	err := syscall.Mkfifo(name, 0600)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &NamedPipeConn{
// 		name:     name,
// 		Path:     path,
// 		fromProc: fromProc,
// 		toProc:   toProc,
// 	}, nil
// }

// // FromProc ...
// func (npc *NamedPipeConn) FromProc() proc.Proc {
// 	return npc.fromProc
// }

// // Out ...
// func (npc *NamedPipeConn) Out() (io.Writer, error) {
// 	return os.OpenFile(npc.Path, os.O_WRONLY, 0600)
// }

// // ToProc ...
// func (npc *NamedPipeConn) ToProc() proc.Proc {
// 	return npc.toProc
// }

// // In ...
// func (npc *NamedPipeConn) In() (io.Reader, error) {
// 	return os.OpenFile(npc.Path, os.O_RDONLY, 0600)
// }

// // Name ...
// func (npc *NamedPipeConn) Name() string {
// 	return npc.name
// }

// // Type ...
// func (npc *NamedPipeConn) Type() proc.ConnType {
// 	return proc.ConnTypeNamedPipe
// }
