package run

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/jaqmol/approx/definition"
	"github.com/jaqmol/approx/processor"
)

// Process ...
type Process struct {
	Cmd exec.Cmd
	Def definition.Definition
}

// SetStdin ...
func (p *Process) SetStdin(r io.Reader) {
	p.Cmd.Stdin = r
}

// SetStdout ...
func (p *Process) SetStdout(w io.Writer) {
	p.Cmd.Stdout = w
}

// Definition ...
func (p *Process) Definition() definition.Definition {
	return p.Def
}

// MakeProcess ...
func MakeProcess(def *definition.Definition) *Process {
	proc := Process{
		Cmd: *exec.Command(def.Command),
		Def: *def,
	}
	proc.Cmd.Env = envSliceFromMap(def.Env)
	return &proc
}

// MakeProcessors ...
func MakeProcessors(definitions []definition.Definition) []processor.Processor {
	processors := make([]processor.Processor, len(definitions))

	idx := 0
	for _, def := range definitions {
		var proc processor.Processor
		switch def.Type {
		case definition.TypeHTTPServer:
			proc = builtin.MakeHTTPServer(&def)
		case definition.TypeFork:
			proc = builtin.MakeFork(&def)
		case definition.TypeMerge:
			proc = builtin.MakeMerge(&def)
		case definition.TypeProcess:
			proc = MakeProcess(&def)
		}
		// ^^ append(os.Environ(), ...)
		processors[idx] = proc
		idx++
	}

	return processors
}

func envSliceFromMap(envMap map[string]string) []string {
	acc := make([]string, len(envMap))
	idx := 0
	for key, value := range envMap {
		acc[idx] = fmt.Sprintf("%v=%v", key, value)
		idx++
	}
	return acc
}
