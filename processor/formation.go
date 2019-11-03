package processor

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/project"
)

// TODO: Testing

// Formation ...
type Formation struct {
	Configuration *configuration.Formation
	Processors    map[string]Processor
}

// NewFormation ...
func NewFormation() (*Formation, error) {
	projForm, err := loadProjectFormation()
	if err != nil {
		return nil, err
	}
	confForm, err := loadConfigFormation(projForm)
	if err != nil {
		return nil, err
	}
	procForID := make(map[string]Processor)

	err = createProcessors(confForm, procForID)
	if err != nil {
		return nil, err
	}

	err = connectProcessors(confForm, procForID)
	if err != nil {
		return nil, err
	}

	return &Formation{
		Configuration: confForm,
		Processors:    procForID,
	}, nil
}

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

func createProcessors(form *configuration.Formation, procForID map[string]Processor) error {
	return form.FlowTree.Iterate(func(prev []*configuration.FlowNode, curr *configuration.FlowNode, _ []*configuration.FlowNode) error {
		switch curr.Processor().Type() {
		case configuration.StdinType:
			getCreateProcessor(curr.Processor(), procForID)
		case configuration.MergeType:
			getCreateProcessor(curr.Processor(), procForID)
		default:
			if len(prev) != 1 {
				return fmt.Errorf("Expected precisely 1 input for type of processor \"%v\"", curr.Processor().ID())
			}
			getCreateProcessor(curr.Processor(), procForID)
		}
		return nil
	})
}

func getCreateProcessor(
	currConfProc configuration.Processor,
	procForID map[string]Processor,
) (Processor, error) {
	id := currConfProc.ID()
	pp, ok := procForID[id]
	var err error
	if !ok {
		switch currConfProc.Type() {
		case configuration.StdinType:
			pp = &Stdin
		case configuration.CommandType:
			pp, err = NewCommand(currConfProc.(*configuration.Command))
		case configuration.ForkType:
			pp, err = NewFork(currConfProc.(*configuration.Fork))
		case configuration.MergeType:
			pp, err = NewMerge(currConfProc.(*configuration.Merge))
		case configuration.StdoutType:
			pp = &Stdout
		}
		procForID[id] = pp
	}
	return pp, err
}

func connectProcessors(form *configuration.Formation, procForID map[string]Processor) error {
	return form.FlowTree.Iterate(func(prev []*configuration.FlowNode, curr *configuration.FlowNode, _ []*configuration.FlowNode) error {
		prevIDs := collectNodeIDs(prev)
		readers := collectOutputReaders(prevIDs, procForID)
		currID := curr.Processor().ID()
		currProc, ok := procForID[currID]
		if !ok {
			return fmt.Errorf("Processor to connect to not found: %v", currID)
		}
		return currProc.Connect(readers...)
	})
}

func collectNodeIDs(nodes []*configuration.FlowNode) []string {
	acc := make([]string, len(nodes))
	for i, n := range nodes {
		acc[i] = n.Processor().ID()
	}
	return acc
}

func collectOutputReaders(procIDs []string, procForID map[string]Processor) []io.Reader {
	acc := make([]io.Reader, len(procIDs))
	for i, id := range procIDs {
		procProc := procForID[id]
		acc[i] = procProc.Out()
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

func loadConfigFormation(projForm *project.Formation) (*configuration.Formation, error) {
	confForm, err := configuration.NewFormation(projForm)
	if err != nil {
		return nil, err
	}
	return confForm, nil
}

// func toProcList(procForID map[string]Processor) []Processor {
// 	acc := make([]Processor, len(procForID))
// 	i := 0
// 	for _, p := range procForID {
// 		acc[i] = p
// 		i++
// 	}
// 	return acc
// }
