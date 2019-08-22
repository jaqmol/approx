package processor

import (
	"io"

	"github.com/jaqmol/approx/definition"
)

// Processor ...
type Processor interface {
	SetStdin(io.Reader)
	SetStdout(io.Writer)
	SetStderr(io.Writer)
	Definition() *definition.Definition
	Start()
}
