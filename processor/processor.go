package processor

import (
	"io"

	"github.com/jaqmol/approx/configuration"
)

// Processor ...
type Processor interface {
	Conf() configuration.Processor
	Start()
	Outs() []io.Reader
	Err() io.Reader
}

type procPipe struct {
	reader io.Reader
	writer io.Writer
}

func newProcPipe() procPipe {
	r, w := io.Pipe()
	return procPipe{
		reader: r,
		writer: w,
	}
}
