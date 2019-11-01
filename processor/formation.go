package processor

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/project"
)

// TODO: Testing

// Formation ...
type Formation struct {
	Processors []Processor
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

	confForm.FlowTree.Iterate(func(prev []*configuration.FlowNode, curr *configuration.FlowNode) {
		switch curr.Processor().Type() {
		case configuration.MergeType:
			ids := collectNodeIDs(prev)
			getCreateProcessor(curr.Processor(), procForID, ids...)
		default:
			if len(prev) != 1 {
				log.Fatalf("Expected precisely 1 input for type of processor \"%v\"\n", curr.Processor().ID())
			}
			id := prev[0].Processor().ID()
			getCreateProcessor(curr.Processor(), procForID, id)
		}
	})

	return &Formation{
		Processors: toProcList(procForID),
	}, nil
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

func getCreateProcessor(
	currConfProc configuration.Processor,
	procForID map[string]Processor,
	prevIDs ...string,
) Processor {
	id := currConfProc.ID()
	pp, ok := procForID[id]
	if !ok {
		switch currConfProc.Type() {
		case configuration.CommandType:
			pp = NewCommand(
				currConfProc.(*configuration.Command),
				procForID[prevIDs[0]].Out(),
			)
		case configuration.ForkType:
			pp = NewFork(
				currConfProc.(*configuration.Fork),
				procForID[prevIDs[0]].Out(),
			)
		case configuration.MergeType:
			pp = NewMerge(
				currConfProc.(*configuration.Merge),
				collectOutputReaders(prevIDs, procForID),
			)
		}
		procForID[id] = pp
	}
	return pp
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

func toProcList(procForID map[string]Processor) []Processor {
	acc := make([]Processor, len(procForID))
	idx := 0
	for _, p := range procForID {
		acc[0] = p
		idx++
	}
	return acc
}
