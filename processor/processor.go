package processor

import (
	"github.com/jaqmol/approx/definition"
	"github.com/jaqmol/approx/pipe"
)

// Processor ...
type Processor interface {
	SetStdin(*pipe.Reader)
	SetStdout(*pipe.Writer)
	SetStderr(*pipe.Writer)
	Definition() *definition.Definition
	Start()
}
