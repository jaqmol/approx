package run

import (
	"fmt"
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
	flowProc *flow.ProcItem
	dataProv DataProvider
	cmd      *exec.Cmd
}

// NewProcItem ...
func NewProcItem(errMsg *axmsg.Errors, dataProv DataProvider, flowProc *flow.ProcItem) *ProcItem {
	return &ProcItem{
		errMsg:   errMsg,
		flowProc: flowProc,
		dataProv: dataProv,
	}
}

// Init ...
func (p *ProcItem) Init(check *ViabilityCheck) {
	workDir, cmdStr, args := p.commandAndArgs()
	p.cmd = exec.Command(cmdStr, args...)
	p.cmd.Dir = workDir
	p.cmd.Stderr = os.Stderr
	envs := make([]string, 0)
	envs = append(envs, os.Environ()...)
	envs = p.attachInEnvs(envs, check)
	envs = p.attachOutEnvs(envs, check)
	envs = append(envs, os.Environ()...)
	envs = append(envs, p.flowProc.Conf.Environment()...)
	p.cmd.Env = envs
}

// Start ...
func (p *ProcItem) Start(prematureExitChan chan<- string, errChan chan<- error) {
	err := p.cmd.Start()
	if err != nil {
		errChan <- err
	} else {
		prematureExitChan <- p.flowProc.Conf.Name()
	}
	go func() {
		err := p.cmd.Wait()
		if err != nil {
			errChan <- err
		} else {
			prematureExitChan <- p.flowProc.Conf.Name()
		}
	}()
}

func (p *ProcItem) attachInEnvs(envs []string, check *ViabilityCheck) []string {
	if p.inputIsStdin() {
		inFlowConn := p.flowProc.Prev[0]
		runConn := p.dataProv.Connection(inFlowConn.Hash)
		p.cmd.Stdin = runConn.Reader()
		check.AddInType(inFlowConn.Hash, ConnectionTypeStdin)
		// log.Printf("Using reader for pipe »%v« as output\n", inFlowConn.Hash)
	} else {
		envs = append(envs, p.envInCount())
		for i, inFlowConn := range p.flowProc.Prev {
			rc := p.dataProv.Connection(inFlowConn.Hash)
			envs = append(envs, envInIdxPath(i, rc.PipePath()))
			check.AddInType(inFlowConn.Hash, ConnectionTypeEnvInPipe)
		}
	}
	return envs
}

func (p *ProcItem) attachOutEnvs(envs []string, check *ViabilityCheck) []string {
	if p.outputIsStdout() {
		outFlowConn := p.flowProc.Next[0]
		runConn := p.dataProv.Connection(outFlowConn.Hash)
		p.cmd.Stdout = runConn.Writer()
		check.AddOutType(outFlowConn.Hash, ConnectionTypeStdout)
		// log.Printf("Using writer for pipe »%v« as output\n", outFlowConn.Hash)
	} else {
		envs = append(envs, p.envOutCount())
		for i, outFlowConn := range p.flowProc.Next {
			rc := p.dataProv.Connection(outFlowConn.Hash)
			envs = append(envs, envOutIdxPath(i, rc.PipePath()))
			check.AddOutType(outFlowConn.Hash, ConnectionTypeEnvOutPipe)
		}
	}
	return envs
}

func (p *ProcItem) envInCount() string {
	return fmt.Sprintf("IN_COUNT=%v", len(p.flowProc.Prev))
}

func envInIdxPath(idx int, path string) string {
	return fmt.Sprintf("IN_%v=%v", idx, path)
}

func (p *ProcItem) envOutCount() string {
	return fmt.Sprintf("OUT_COUNT=%v", len(p.flowProc.Next))
}

func envOutIdxPath(idx int, path string) string {
	return fmt.Sprintf("OUT_%v=%v", idx, path)
}

func (p *ProcItem) inputIsStdin() bool {
	if len(p.flowProc.Prev) == 1 {
		return "stdin" == p.flowProc.Conf.Inputs()[0]
	}
	return false
}

func (p *ProcItem) outputIsStdout() bool {
	if len(p.flowProc.Next) == 1 {
		return "stdout" == p.flowProc.Conf.Outputs()[0]
	}
	return false
}

func (p *ProcItem) commandAndArgs() (string, string, []string) {
	var workingDirectory string
	var command string
	var args []string
	switch p.flowProc.Conf.Type() {
	case conf.TypeProcess:
		workingDirectory = p.dataProv.FormationBasePath()
		processorConf := p.flowProc.Conf.(*conf.ProcessConf)
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
