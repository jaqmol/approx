package processor

import (
	"github.com/jaqmol/approx/channel"
	"github.com/jaqmol/approx/definition"
)

// Processor ...
type Processor interface {
	SetStdin(channel.Reader)
	SetStdout(channel.Writer)
	SetStderr(channel.Writer)
	Definition() *definition.Definition
	Start()
}
