package run

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	"github.com/jaqmol/approx/errormsg"
	"github.com/jaqmol/approx/flow"
)

// State ...
type State struct {
	errMsg        *errormsg.ErrorMsg
	pipesBasePath string
	pipePaths     []string
}

// Flow ...
func Flow(errMsg *errormsg.ErrorMsg, fl *flow.Flow) *State {
	tmpDir, err := ioutil.TempDir("", "approx")
	if err != nil {
		errMsg.LogFatal(nil, "Error getting temp dir: %v", err.Error())
	}
	basePath := filepath.Join(tmpDir, fl.MainItem.Conf.Name())
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		os.Mkdir(basePath, os.ModePerm)
	}
	s := &State{
		errMsg:        errMsg,
		pipesBasePath: basePath,
		pipePaths:     make([]string, 0),
	}
	s.createPipes(fl)
	return s
}

func (s *State) createPipes(fl *flow.Flow) {
	fl.IterateConns(func(row []*flow.ConnItem) {
		for _, conn := range row {
			pipeName := fmt.Sprintf("%v.pipe", conn.Hash)
			pp := filepath.Join(s.pipesBasePath, pipeName)
			err := syscall.Mkfifo(pp, 0600)
			if err != nil {
				s.errMsg.LogFatal(nil, "Error creating pipe: %v", err.Error())
			}
			s.pipePaths = append(s.pipePaths, pp)
		}
	})
}

// Cleanup ...
func (s *State) Cleanup() []error {
	errs := make([]error, 0)
	for _, pp := range s.pipePaths {
		err := os.Remove(pp)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	err := os.Remove(s.pipesBasePath)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	return nil
}
