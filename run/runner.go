package run

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jaqmol/approx/axmsg"
	"github.com/jaqmol/approx/flow"
)

// Runner ...
type Runner struct {
	flow         *flow.Flow
	errMsg       *axmsg.Errors
	pipeBasePath string
	connForHash  map[uint32]*ConnItem
	procForName  map[string]*ProcItem
	// errChan      chan error
}

// NewRunner ...
func NewRunner(errMsg *axmsg.Errors, fl *flow.Flow) *Runner {
	pipeBasePath := preparePipeBasePath(errMsg, fl)
	r := &Runner{
		flow:         fl,
		errMsg:       errMsg,
		pipeBasePath: pipeBasePath,
		connForHash:  createConnections(errMsg, pipeBasePath, fl),
		// errChan:      make(chan error, 0),
	}
	r.procForName = createProcessors(errMsg, r, fl)
	return r
}

func preparePipeBasePath(errMsg *axmsg.Errors, fl *flow.Flow) string {
	tmpDir, err := ioutil.TempDir("", "approx")
	if err != nil {
		errMsg.LogFatal(nil, "Error getting temp dir: %v", err.Error())
	}
	basePath := filepath.Join(tmpDir, fl.MainItem.Conf.Name())
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		os.MkdirAll(basePath, os.ModePerm)
	}
	return basePath
}

func createConnections(errMsg *axmsg.Errors, pipeBasePath string, fl *flow.Flow) map[uint32]*ConnItem {
	acc := make(map[uint32]*ConnItem, 0)
	fl.IterateConns(func(row []*flow.ConnItem) {
		for _, flowConn := range row {
			runConn, err := NewConnItem(errMsg, pipeBasePath, flowConn)
			if err != nil {
				errMsg.LogFatal(nil, err.Error())
			}
			acc[flowConn.Hash] = runConn
		}
	})
	return acc
}

func createProcessors(errMsg *axmsg.Errors, dataProv DataProvider, fl *flow.Flow) map[string]*ProcItem {
	acc := make(map[string]*ProcItem, 0)
	fl.IterateProcs(func(row []*flow.ProcItem) {
		for _, flowProc := range row {
			runProc := NewProcItem(errMsg, dataProv, flowProc)
			acc[flowProc.Conf.Name()] = runProc
		}
	})
	return acc
}

// Cleanup ...
func (r *Runner) Cleanup() []error {
	errs := make([]error, 0)
	for _, c := range r.connForHash {
		err := c.Cleanup()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	err := os.Remove(r.pipeBasePath)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	return nil
}

// Connection ...
func (r *Runner) Connection(hash uint32) *ConnItem {
	return r.connForHash[hash]
}

// FormationBasePath ...
func (r *Runner) FormationBasePath() string {
	return r.flow.FormationBasePath
}

// CheckViability ...
func (r *Runner) CheckViability() {

}

// InitProcessors ...
func (r *Runner) InitProcessors() {
	check := NewViabilityCheck()
	for _, proc := range r.procForName {
		proc.Init(check)
	}
	if check.InsAndOutsAreBalanced() {
		log.Println("Flow is operable, inputs and outputs are balanced")
	} else {
		log.Println("Flow is not operable, inputs and outputs are not balanced")
	}
}

// Start ...
func (r *Runner) Start() {
	prematureExitChan := make(chan string, 0)
	errChan := make(chan error, 0)
	go func() {
		for _, proc := range r.procForName {
			proc.Start(prematureExitChan, errChan)
		}
	}()
	select {
	case procName := <-prematureExitChan:
		r.errMsg.Log(nil, "Processor %v exited prematurely, please make sure the process is long running", procName)
	case err := <-errChan:
		r.errMsg.LogFatal(nil, err.Error())
	}
}
