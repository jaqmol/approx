package run

import (
	"fmt"

	"github.com/jaqmol/approx/pipe"
	"github.com/jaqmol/approx/processor"
)

// Connect ...
func Connect(processors []processor.Processor, flows map[string][]string, pipes map[string]pipe.Pipe, stdErrs map[string]pipe.Pipe) {
	processorForName := makeProcessorForNameMap(processors)

	for _, fromProc := range processors {
		fromName := fromProc.Definition().Name
		errPipe := stdErrs[fromName]
		fromProc.SetStderr(errPipe.Writer)
		fromProc := processorForName[fromName]
		toNames := flows[fromName]

		for _, toName := range toNames {
			toProc := processorForName[toName]
			key := PipeKey(fromName, toName)
			aPipe := pipes[key]
			fmt.Printf("Connecting: %v: %v\n", key, aPipe)

			fromProc.SetStdout(aPipe.Writer)
			toProc.SetStdin(aPipe.Reader)
		}
	}
}

func makeProcessorForNameMap(processors []processor.Processor) map[string]processor.Processor {
	acc := make(map[string]processor.Processor)
	for _, proc := range processors {
		acc[proc.Definition().Name] = proc
	}
	return acc
}

// func findCommands(names []string, processors map[string]processor.Processor) []processor.Processor {
// 	acc := make([]processor.Processor, len(names))
// 	for idx, name := range names {
// 		acc[idx] = processors[name]
// 	}
// 	return acc
// }
