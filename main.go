package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jaqmol/approx/processor"

	"github.com/jaqmol/approx/configuration"

	"github.com/jaqmol/approx/project"
)

func main() {
	// TODO: Goes to package processor.Formation:

	projForm := loadProjectFormation()
	confForm := loadConfigFormation(projForm)
	procProcForID := make(map[string]processor.Processor)

	confForm.FlowTree.Iterate(func(prev []configuration.FlowNode, curr *configuration.FlowNode) {
		switch curr.Processor().Type() {
		case configuration.MergeType:
			ids := collectNodeIDs(prev)
			getCreateProcProc(curr.Processor(), procProcForID, ids...)
		default:
			if len(prev) != 1 {
				log.Fatalf("Expected precisely 1 input for type of processor \"%v\"\n", curr.Processor().ID())
			}
			id := prev[0].Processor().ID()
			getCreateProcProc(curr.Processor(), procProcForID, id)
		}
	})

	// TODO: Add Start() to processor.Formation
}

// TODO: Goes to package processor.Formation:
// TODO: Testing

func getCreateProcProc(currConfProc configuration.Processor, procProcForID map[string]processor.Processor, prevIDs ...string) processor.Processor {
	id := currConfProc.ID()
	pp, ok := procProcForID[id]
	if !ok {
		switch currConfProc.Type() {
		case configuration.CommandType:
			pp = processor.NewCommand(
				currConfProc.(*configuration.Command),
				procProcForID[prevIDs[0]].Out(),
			)
		case configuration.ForkType:
			pp = processor.NewFork(
				currConfProc.(*configuration.Fork),
				procProcForID[prevIDs[0]].Out(),
			)
		case configuration.MergeType:
			pp = processor.NewMerge(
				currConfProc.(*configuration.Merge),
				collectOutputReaders(prevIDs, procProcForID),
			)
		}
		procProcForID[id] = pp
	}
	return pp
}

func collectNodeIDs(nodes []configuration.FlowNode) []string {
	acc := make([]string, len(nodes))
	for i, n := range nodes {
		acc[i] = n.Processor().ID()
	}
	return acc
}

func collectOutputReaders(procIDs []string, procProcForID map[string]processor.Processor) []io.Reader {
	acc := make([]io.Reader, len(procIDs))
	for i, id := range procIDs {
		procProc := procProcForID[id]
		acc[i] = procProc.Out()
	}
	return acc
}

func loadProjectFormation() *project.Formation {
	var projPath string
	var err error
	if len(os.Args) == 2 {
		projPath = os.Args[1]
	} else {
		projPath, err = os.Getwd()
	}
	if err != nil {
		log.Fatalln(err)
	}
	projPath, err = filepath.Abs(projPath)
	if err != nil {
		log.Fatalln(err)
	}
	projForm, err := project.LoadFormation(projPath)
	if err != nil {
		log.Fatalln(err)
	}
	return projForm
}

func loadConfigFormation(projForm *project.Formation) *configuration.Formation {
	confForm, err := configuration.NewFormation(projForm)
	if err != nil {
		log.Fatalln(err)
	}
	return confForm
}
