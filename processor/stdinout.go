package processor

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/event"
)

// stdinReader ...
var stdinReader io.Reader

// stdoutWriter ...
var stdoutWriter io.Writer

// stderrWriter ...
var stderrWriter io.Writer

// Stdin ...
var Stdin *stdinProc

// Stdout ...
var Stdout *stdoutProc

func init() {
	stdinReader = os.Stdin
	stdoutWriter = os.Stdout
	stderrWriter = os.Stderr

	Stdin = &stdinProc{newProcPipe()}
	Stdout = &stdoutProc{
		err: newProcPipe(),
		in:  nil,
	}
}

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
	return []io.Reader{stdinReader}
}

// Out ...
func (s *stdinProc) Out() io.Reader {
	return stdinReader
}

// Err ...
func (s *stdinProc) Err() io.Reader {
	return s.err.reader()
}

// Connect ...
func (s *stdinProc) Connect(inputs ...io.Reader) error {
	if len(inputs) > 0 {
		return fmt.Errorf("Stdin cannot be connected")
	}
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
	scanner := event.NewScanner(s.in)
	for scanner.Scan() {
		msg := evntEndedCopy(scanner.Bytes())
		n, err := stdoutWriter.Write(msg)
		if err != nil {
			log.Fatalln(err.Error())
		}
		if n != len(msg) {
			log.Fatalln("Stdout couldn't write complete event")
		}
	}
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
	// log.Println("processor.Stdout did connect to an input") TODO: REMOVE
	return nil
}

// DebugChangeStdin for testing
func DebugChangeStdin(altStdin io.Reader) error {
	if !inApproxDevEnv() {
		return fmt.Errorf("Stdin can only be changed in development environment")
	}
	stdinReader = altStdin
	return nil
}

// DebugResetStdin ...
func DebugResetStdin() {
	stdinReader = os.Stdin
}

// DebugChangeStdout for testing
func DebugChangeStdout(altStdout io.Writer) error {
	if !inApproxDevEnv() {
		return fmt.Errorf("Stdout can only be changed in development environment")
	}
	stdoutWriter = altStdout
	return nil
}

// DebugResetStdout ...
func DebugResetStdout() {
	stdoutWriter = os.Stdout
}

// DebugChangeStderr for testing
func DebugChangeStderr(altStderr io.Writer) error {
	if !inApproxDevEnv() {
		return fmt.Errorf("Stderr can only be changed in development environment")
	}
	stderrWriter = altStderr
	return nil
}

// DebugResetStderr ...
func DebugResetStderr() {
	stderrWriter = os.Stderr
}

func inApproxDevEnv() bool {
	approxEnv := strings.ToLower(os.Getenv("APPROX_ENV"))
	return approxEnv == "development" || approxEnv == "dev"
}
