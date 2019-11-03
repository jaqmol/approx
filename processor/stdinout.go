package processor

import (
	"fmt"
	"io"
	"os"

	"github.com/jaqmol/approx/configuration"
)

// stdinProc ...
type stdinProc struct {
	err *procPipe
}

// Start ...
func (s *stdinProc) Start() {}

// Conf ...
func (s *stdinProc) Conf() configuration.Processor {
	return &configuration.Stdin
}

// Outs ...
func (s *stdinProc) Outs() []io.Reader {
	return []io.Reader{os.Stdin}
}

// Out ...
func (s *stdinProc) Out() io.Reader {
	return os.Stdin
}

// Err ...
func (s *stdinProc) Err() io.Reader {
	return s.err.reader()
}

// Connect ...
func (s *stdinProc) Connect(inputs ...io.Reader) error {
	return nil
}

// stdoutProc ...
type stdoutProc struct {
	err *procPipe
	in  io.Reader
}

// Start ...
func (s *stdoutProc) Start() {
	go s.start()
}
func (s *stdoutProc) start() {
	io.Copy(os.Stdout, s.in)
}

// Conf ...
func (s *stdoutProc) Conf() configuration.Processor {
	return &configuration.Stdout
}

// Outs ...
func (s *stdoutProc) Outs() []io.Reader {
	return nil
}

// Out ...
func (s *stdoutProc) Out() io.Reader {
	return nil
}

// Err ...
func (s *stdoutProc) Err() io.Reader {
	return s.err.reader()
}

// Connect ...
func (s *stdoutProc) Connect(inputs ...io.Reader) error {
	if s.in != nil {
		return fmt.Errorf("Stdout can only be connected once")
	}
	s.in = inputs[0]
	return nil
}

// Stdin ...
var Stdin stdinProc

// Stdout ...
var Stdout stdoutProc

func init() {
	Stdin = stdinProc{newProcPipe()}
	Stdout = stdoutProc{
		err: newProcPipe(),
		in:  nil,
	}
}
