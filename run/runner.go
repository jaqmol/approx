package run

import (
	"io/ioutil"
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
}

// NewRunner ...
func NewRunner(errMsg *axmsg.Errors, fl *flow.Flow) *Runner {
	pipeBasePath := preparePipeBasePath(errMsg, fl)
	r := &Runner{
		flow:         fl,
		errMsg:       errMsg,
		pipeBasePath: pipeBasePath,
		connForHash:  createConnections(errMsg, pipeBasePath, fl),
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
		os.Mkdir(basePath, os.ModePerm)
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

func createProcessors(errMsg *axmsg.Errors, DataProv DataProvider, fl *flow.Flow) map[string]*ProcItem {
	acc := make(map[string]*ProcItem, 0)
	fl.IterateProcs(func(row []*flow.ProcItem) {
		for _, flowProc := range row {
			runProc := NewProcItem(errMsg, DataProv, flowProc)
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

// Start ...
func (r *Runner) Start() {
	for _, proc := range r.procForName {
		proc.Start()
	}
}
