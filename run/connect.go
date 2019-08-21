package run

import "github.com/jaqmol/approx/processor"

// Connect ...
func Connect(processors []processor.Processor, flows map[string][]string, pipes map[string]Pipe) {
	processorForName := makeProcessorForNameMap(processors)

	for _, fromProc := range processors {
		fromName := fromProc.Definition().Name
		fromProc := processorForName[fromName]
		toNames := flows[fromName]

		for _, toName := range toNames {
			toProc := processorForName[toName]
			pipe := pipes[PipeKey(fromName, toName)]

			fromProc.SetStdout(pipe.Writer)
			toProc.SetStdin(pipe.Reader)
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
