package processor

import (
	"io"

	"github.com/jaqmol/approx/configuration"
)

// StdInOut ...
type StdInOut struct {
	conf  configuration.Processor
	err   *procPipe
	inout io.Reader
}

// NewStdin ...
func NewStdin() *StdInOut {
	return &StdInOut{
		conf:  configuration.Stdin,
		err:   newProcPipe(),
		inout: nil,
	}
}

// NewStdout ...
func NewStdout() *StdInOut {
	return &StdInOut{
		conf:  configuration.Stdout,
		err:   newProcPipe(),
		inout: nil,
	}
}

// Start ...
func (s *StdInOut) Start() {}

// Conf ...
func (s *StdInOut) Conf() configuration.Processor {
	return s.conf
}

// Outs ...
func (s *StdInOut) Outs() []io.Reader {
	return []io.Reader{s.inout}
}

// Out ...
func (s *StdInOut) Out() io.Reader {
	return s.inout
}

// Err ...
func (s *StdInOut) Err() io.Reader {
	return s.err.reader()
}

// Connect ...
func (s *StdInOut) Connect(inputs ...io.Reader) error {
	err := errorIfInvalidConnect(configuration.Stdin.ID(), inputs, s.inout != nil)
	if err != nil {
		return err
	}
	s.inout = inputs[0]
	return nil
}