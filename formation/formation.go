package formation

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jaqmol/approx/actor"
	"github.com/jaqmol/approx/config"
	"github.com/jaqmol/approx/logger"
	"github.com/jaqmol/approx/project"
)

const actorInboxSize = 10

// Formation ...
type Formation struct {
	conf     *config.Formation
	Actables map[string]actor.Actable
	Stdin    io.ReadCloser
	Stdout   io.WriteCloser
	Logger   *logger.Logger
}

// NewFormation ...
func NewFormation(
	stdin io.ReadCloser,
	stderr io.Writer,
	stdout io.WriteCloser,
) (*Formation, error) {
	projForm, err := loadProjectFormation()
	if err != nil {
		return nil, err
	}
	confForm, err := config.NewFormation(projForm)
	if err != nil {
		return nil, err
	}

	f := Formation{
		conf:     confForm,
		Actables: make(map[string]actor.Actable),
		Logger:   logger.NewLogger(stderr),
	}

	err = f.createActables()
	if err != nil {
		return nil, err
	}
	err = f.connectActables()
	if err != nil {
		return nil, err
	}

	return &f, nil
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

func (f *Formation) createActables() error {
	return f.conf.FlowTree.Iterate(func(
		prev []*config.FlowNode,
		curr *config.FlowNode,
		_ []*config.FlowNode,
	) error {
		// TODO: Rename FlowNode.Processor() -> FlowNode.Actor
		var err error
		switch curr.Processor().Type() {
		case config.StdinType:
			fallthrough
		case config.MergeType:
			_, err = f.getCreateActable(curr.Processor())
		default:
			if len(prev) != 1 {
				return fmt.Errorf("Expected precisely 1 input for type of processor \"%v\"", curr.Processor().ID())
			}
			_, err = f.getCreateActable(curr.Processor())
		}
		return err
	})
}

func (f *Formation) getCreateActable(currConfProc config.Processor) (actor.Actable, error) {
	id := currConfProc.ID()
	actbl, ok := f.Actables[id]
	var err error
	if !ok {
		switch currConfProc.Type() {
		case config.StdinType:
			actbl = newStdinActor(f.Stdin)
		case config.CommandType:
			actbl, err = f.newCommandActor(currConfProc.(*config.Command))
		case config.ForkType:
			actbl = f.newForkActor(currConfProc.(*config.Fork))
		case config.MergeType:
			actbl = f.newMergeActor(currConfProc.(*config.Merge))
		case config.StdoutType:
			actbl = newStdoutActor(f.Stdout)
		}
		f.Actables[id] = actbl
	}
	return actbl, err
}

func (f *Formation) newCommandActor(conf *config.Command) (*actor.Command, error) {
	c, err := actor.NewCommandFromConf(actorInboxSize, conf)
	if err != nil {
		return nil, err
	}
	logReader, logWriter := io.Pipe()
	f.Logger.Add(logReader)
	c.Logging(logWriter)
	return c, nil
}

func (f *Formation) newForkActor(conf *config.Fork) *actor.Fork {
	return actor.NewFork(actorInboxSize, conf.Ident, conf.Count)
}

func (f *Formation) newMergeActor(conf *config.Merge) *actor.Merge {
	return actor.NewMerge(actorInboxSize, conf.Ident, conf.Count)
}

func (f *Formation) connectActables() error {
	return f.conf.FlowTree.Iterate(func(
		_ []*config.FlowNode,
		curr *config.FlowNode,
		next []*config.FlowNode,
	) error {
		nextActables, err := f.getActables(getNodeIDs(next))
		if err != nil {
			return err
		}
		currID := curr.Processor().ID() // TODO: Rename Processor
		currActbl, ok := f.Actables[currID]
		if !ok {
			return fmt.Errorf("Processor to connect to not found: %v", currID)
		}

		if len(nextActables) > 0 {
			currActbl.Next(nextActables...)
		}
		return nil
	})
}

func getNodeIDs(nodes []*config.FlowNode) []string {
	acc := make([]string, len(nodes))
	for i, n := range nodes {
		acc[i] = n.Processor().ID()
	}
	return acc
}

func (f *Formation) getActables(ids []string) ([]actor.Actable, error) {
	acc := make([]actor.Actable, len(ids))
	for i, id := range ids {
		actbl, ok := f.Actables[id]
		if !ok {
			return nil, fmt.Errorf("Could not find \"%v\"", id)
		}
		acc[i] = actbl
	}
	return acc, nil
}

// Start ...
func (f *Formation) Start() {
	for _, actbl := range f.Actables {
		actbl.Start()
	}
	go f.Logger.Start()
}
