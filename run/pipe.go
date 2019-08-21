package run

import (
	"fmt"
	"io"

	"github.com/jaqmol/approx/processor"
)

// MakePipes ...
func MakePipes(processors []processor.Processor, flows map[string][]string) map[string]Pipe {
	acc := make(map[string]Pipe)

	for _, fromProc := range processors {
		fromName := fromProc.Definition().Name
		toNames := flows[fromName]

		for _, toName := range toNames {
			key := PipeKey(fromName, toName)
			reader, writer := io.Pipe()
			acc[key] = Pipe{Reader: reader, Writer: writer}
		}
	}

	return acc
}

// Pipe ...
type Pipe struct {
	Reader io.Reader
	Writer io.Writer
}

// PipeKey ...
func PipeKey(fromName string, toName string) string {
	return fmt.Sprintf("%v->%v", fromName, toName)
}
