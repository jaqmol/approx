package run

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jaqmol/approx/axmsg"

	"github.com/jaqmol/approx/conf"
	"github.com/jaqmol/approx/flow"
)

// ProcItem ...
type ProcItem struct {
	errMsg   *axmsg.Errors
	FlowProc *flow.ProcItem
	dataProv DataProvider
	cmd      *exec.Cmd
}

// NewProcItem ...
func NewProcItem(errMsg *axmsg.Errors, dataProv DataProvider, flowProc *flow.ProcItem) *ProcItem {
	return &ProcItem{
		errMsg:   errMsg,
		FlowProc: flowProc,
		dataProv: dataProv,
	}
}

// Start ...
func (p *ProcItem) Start() {
	workDir, cmdStr, args := p.commandAndArgs()
	p.cmd = exec.Command(cmdStr, args...)
	p.cmd.Dir = workDir
	envs := make([]string, 0)
	envs = append(envs, os.Environ()...)
	envs = p.attachInEnvs(envs)
	envs = p.attachOutEnvs(envs)
	log.Printf("%v\n", envs)
	p.cmd.Env = envs
}

func (p *ProcItem) attachInEnvs(envs []string) []string {
	if p.inputIsStdin() {
		inFlowConn := p.FlowProc.Prev[0]
		runConn := p.dataProv.Connection(inFlowConn.Hash)
		p.cmd.Stdin = runConn.Reader()
	} else {
		envs = append(envs, p.envInCount())
		for i, inFlowConn := range p.FlowProc.Prev {
			rc := p.dataProv.Connection(inFlowConn.Hash)
			envs = append(envs, envInIdxPath(i, rc.PipePath()))
		}
	}
	return envs
}

func (p *ProcItem) attachOutEnvs(envs []string) []string {
	if p.outputIsStdout() {
		outFlowConn := p.FlowProc.Next[0]
		runConn := p.dataProv.Connection(outFlowConn.Hash)
		p.cmd.Stdout = runConn.Writer()
	} else {
		envs = append(envs, p.envOutCount())
		for i, outFlowConn := range p.FlowProc.Next {
			rc := p.dataProv.Connection(outFlowConn.Hash)
			envs = append(envs, envOutIdxPath(i, rc.PipePath()))
		}
	}
	return envs
}

func (p *ProcItem) envInCount() string {
	return fmt.Sprintf("IN_COUNT=%v", len(p.FlowProc.Prev))
}

func envInIdxPath(idx int, path string) string {
	return fmt.Sprintf("IN_%v=%v", idx, path)
}

func (p *ProcItem) envOutCount() string {
	return fmt.Sprintf("OUT_COUNT=%v", len(p.FlowProc.Next))
}

func envOutIdxPath(idx int, path string) string {
	return fmt.Sprintf("OUT_%v=%v", idx, path)
}

func (p *ProcItem) inputIsStdin() bool {
	if len(p.FlowProc.Prev) == 1 {
		return "stdin" == p.FlowProc.Conf.Inputs()[0]
	}
	return false
}

func (p *ProcItem) outputIsStdout() bool {
	if len(p.FlowProc.Next) == 1 {
		return "stdout" == p.FlowProc.Conf.Outputs()[0]
	}
	return false
}

func (p *ProcItem) commandAndArgs() (string, string, []string) {
	var workingDirectory string
	var command string
	var args []string
	switch p.FlowProc.Conf.Type() {
	case conf.TypeProcess:
		workingDirectory = p.dataProv.FormationBasePath()
		processorConf := p.FlowProc.Conf.(*conf.ProcessConf)
		command = processorConf.Command()
		args = processorConf.Arguments()
	case conf.TypeHTTPServer:
		workingDirectory = p.parentDirectoryPath()
		command = "lib/http_server"
	case conf.TypeFork:
		workingDirectory = p.parentDirectoryPath()
		command = "lib/fork"
	case conf.TypeMerge:
		workingDirectory = p.parentDirectoryPath()
		command = "lib/merge"
	case conf.TypeCheck:
		workingDirectory = p.parentDirectoryPath()
		command = "lib/check"
	}
	return workingDirectory, command, args
}

func (p *ProcItem) parentDirectoryPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		p.errMsg.LogFatal(nil, "Unable to retrieve approx parent directory: %v", err.Error())
	}
	return dir
}
