package processor

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jaqmol/approx/logger"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/project"
)

// TODO: Testing

// Formation ...
type Formation struct {
	Configuration *configuration.Formation
	Processors    map[string]Processor
	Logger        *logger.Logger
}

// NewFormation ...
func NewFormation(stdin *StdInOut, stdout *StdInOut, errWriter io.Writer) (*Formation, error) {
	projForm, err := loadProjectFormation()
	if err != nil {
		return nil, err
	}
	confForm, err := configuration.NewFormation(projForm)
	if err != nil {
		return nil, err
	}

	f := Formation{
		Configuration: confForm,
		Processors:    make(map[string]Processor),
		Logger:        logger.NewLogger(errWriter),
	}

	f.Processors[stdin.Conf().ID()] = stdin
	f.Processors[stdout.Conf().ID()] = stdout

	err = f.createProcessors()
	if err != nil {
		return nil, err
	}

	err = f.connectProcessors()
	if err != nil {
		return nil, err
	}

	f.connectLogger()
	return &f, nil
}

// // Stdin ...
// func (f *Formation) Stdin() *StdInOut {
// 	proc := f.Processors[configuration.Stdin.ID()]
// 	si, ok := proc.(*StdInOut)
// 	if ok {
// 		return si
// 	}
// 	return nil
// }

// // Stdout ...
// func (f *Formation) Stdout() *StdInOut {
// 	proc := f.Processors[configuration.Stdout.ID()]
// 	so, ok := proc.(*StdInOut)
// 	if ok {
// 		return so
// 	}
// 	return nil
// }

// Root ...
func (f *Formation) Root() Processor {
	id := f.Configuration.FlowTree.Root.Processor().ID()
	return f.Processors[id]
}

// Input ...
func (f *Formation) Input() Processor {
	inputNode := f.Configuration.FlowTree.Input
	if inputNode == nil {
		return nil
	}
	id := inputNode.Processor().ID()
	return f.Processors[id]
}

// Output ...
func (f *Formation) Output() Processor {
	outputNode := f.Configuration.FlowTree.Output
	if outputNode == nil {
		return nil
	}
	id := outputNode.Processor().ID()
	return f.Processors[id]
}

// Start ...
func (f *Formation) Start() {
	for _, p := range f.Processors {
		p.Start()
	}
	go f.Logger.Start()
	// TODO: Last step is connect a Collector
	//       and pull data through the flow.
}

// WaitForCommands ...
func (f *Formation) WaitForCommands() {
	for _, p := range f.Processors {
		if p.Conf().Type() == configuration.CommandType {
			cmd := p.(*Command)
			cmd.Wait()
		}
	}
}

func (f *Formation) createProcessors() error {
	return f.Configuration.FlowTree.Iterate(func(
		prev []*configuration.FlowNode,
		curr *configuration.FlowNode,
		_ []*configuration.FlowNode,
	) error {
		switch curr.Processor().Type() {
		case configuration.StdinType:
			f.getCreateProcessor(curr.Processor())
		case configuration.MergeType:
			f.getCreateProcessor(curr.Processor())
		default:
			if len(prev) != 1 {
				return fmt.Errorf("Expected precisely 1 input for type of processor \"%v\"", curr.Processor().ID())
			}
			f.getCreateProcessor(curr.Processor())
		}
		return nil
	})
}

func (f *Formation) getCreateProcessor(currConfProc configuration.Processor) (Processor, error) {
	id := currConfProc.ID()
	pp, ok := f.Processors[id]
	var err error
	if !ok {
		switch currConfProc.Type() {
		case configuration.StdinType:
			pp = NewStdin()
		case configuration.CommandType:
			pp, err = NewCommand(currConfProc.(*configuration.Command))
		case configuration.ForkType:
			pp, err = NewFork(currConfProc.(*configuration.Fork))
		case configuration.MergeType:
			pp, err = NewMerge(currConfProc.(*configuration.Merge))
		case configuration.StdoutType:
			pp = NewStdout()
		}
		f.Processors[id] = pp
	}
	return pp, err
}

func (f *Formation) connectProcessors() error {
	return f.Configuration.FlowTree.Iterate(func(
		prev []*configuration.FlowNode,
		curr *configuration.FlowNode,
		_ []*configuration.FlowNode,
	) error {
		prevIDs := getNodeIDs(prev)
		readers, err := f.getOutputReaders(prevIDs)
		if err != nil {
			return err
		}
		currID := curr.Processor().ID()
		currProc, ok := f.Processors[currID]
		if !ok {
			return fmt.Errorf("Processor to connect to not found: %v", currID)
		}

		if len(readers) > 0 {
			return currProc.Connect(readers...)
		}
		return nil
	})
}

func (f *Formation) connectLogger() {
	for _, proc := range f.Processors {
		f.Logger.Add(proc.Err())
	}
}

func (f *Formation) getOutputReaders(procIDs []string) ([]io.Reader, error) {
	acc := make([]io.Reader, len(procIDs))
	for i, id := range procIDs {
		procProc, ok := f.Processors[id]
		if !ok {
			return nil, fmt.Errorf("Could not find processor for ID %v", id)
		}
		output := procProc.Out()
		if output == nil {
			return nil, fmt.Errorf("Processor %v returned nil output", id)
		}
		acc[i] = output
	}
	return acc, nil
}

func getNodeIDs(nodes []*configuration.FlowNode) []string {
	acc := make([]string, len(nodes))
	for i, n := range nodes {
		acc[i] = n.Processor().ID()
	}
	return acc
}

func loadProjectFormation() (*project.Formation, error) {
	var projPath string
	var err error
	if len(os.Args) == 2 {
		projPath = os.Args[1]
	} else {
		projPath, err = os.Getwd()
	}

	if err != nil {
		return nil, err
	}
	projPath, err = filepath.Abs(projPath)
	if err != nil {
		return nil, err
	}
	projForm, err := project.LoadFormation(projPath)
	if err != nil {
		return nil, err
	}
	return projForm, nil
}

// func loadConfigFormation(projForm *project.Formation) (*configuration.Formation, error) {
// 	confForm, err := configuration.NewFormation(projForm)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return confForm, nil
// }

// func toProcList(procForID map[string]Processor) []Processor {
// 	acc := make([]Processor, len(procForID))
// 	i := 0
// 	for _, p := range procForID {
// 		acc[i] = p
// 		i++
// 	}
// 	return acc
// }
