package run

import (
	"fmt"
	"io"

	"github.com/jaqmol/approx/definition"
)

// MakePipes ...
func MakePipes(definitions []definition.Definition, flows map[string][]string) map[string]Pipe {
	acc := make(map[string]Pipe)

	for _, fromDef := range definitions {
		fromName := fromDef.Name
		toNames := flows[fromName]

		for _, toName := range toNames {
			key := PipeKey(fromName, toName)
			reader, writer := io.Pipe()
			acc[key] = Pipe{Reader: reader, Writer: writer}
		}
	}

	return acc
}

// MakeStderrs ...
func MakeStderrs(definitions []definition.Definition) map[string]Pipe {
	acc := make(map[string]Pipe)

	for _, def := range definitions {
		reader, writer := io.Pipe()
		acc[def.Name] = Pipe{Reader: reader, Writer: writer}
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
