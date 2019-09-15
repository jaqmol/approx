package run

import (
	"fmt"

	"github.com/jaqmol/approx/channel"
	"github.com/jaqmol/approx/processor"
	"github.com/jaqmol/approx/utils"
)

// Connect ...
func Connect(
	processors []processor.Processor,
	flows map[string][]string,
	tappedPipeNames map[string]string,
	pipes map[string]channel.Pipe,
	stdErrs map[string]channel.Pipe,
) {
	processorForName := makeProcessorForNameMap(processors)

	for _, fromProc := range processors {
		fromName := fromProc.Definition().Name
		errPipe := stdErrs[fromName]
		fromProc.SetStderr(errPipe)
		fromProc := processorForName[fromName]
		toNames := flows[fromName]

		for _, toName := range toNames {
			toProc := processorForName[toName]
			key := utils.PipeKey(fromName, toName)
			aPipe := pipes[key]

			fromProc.SetStdout(aPipe)
			toProc.SetStdin(aPipe)

			if aPipe.IsTapped() {
				tapName := tappedPipeNames[key]
				fmt.Printf("Configuring tapped pipe %v\n", tapName)
				tapErrPipe := stdErrs[tapName]
				aPipe.SetStderr(tapErrPipe)
			}
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
