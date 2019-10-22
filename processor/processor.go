package processor

import (
	"io"

	"github.com/jaqmol/approx/configuration"
)

/*
	A	processor is initialized with
	- it's specific type of configuration
	- the output(s) of it's predecessor(s) as it's input(s)
*/

// Processor ...
type Processor interface {
	Start()
	Conf() configuration.Processor
	Outs() []io.Reader
	Out() io.Reader
	Err() io.Reader
}
