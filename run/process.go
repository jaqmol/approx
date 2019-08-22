package run

import (
	"io"
	"os"
	"os/exec"

	"github.com/jaqmol/approx/builtin"
	"github.com/jaqmol/approx/definition"
	"github.com/jaqmol/approx/processor"
)

// Process ...
type Process struct {
	cmd exec.Cmd
	def definition.Definition
}

// SetStdin ...
func (p *Process) SetStdin(r io.Reader) {
	p.cmd.Stdin = r
}

// SetStdout ...
func (p *Process) SetStdout(w io.Writer) {
	p.cmd.Stdout = w
}

// SetStderr ...
func (p *Process) SetStderr(w io.Writer) {
	p.cmd.Stderr = w
}

// Definition ...
func (p *Process) Definition() *definition.Definition {
	return &p.def
}

// Start ...
func (p *Process) Start() {

}

// MakeProcess ...
func MakeProcess(def *definition.Definition) *Process {
	proc := Process{
		cmd: *exec.Command(def.Command),
		def: *def,
	}
	proc.cmd.Env = def.EnvSlice()
	return &proc
}

// MakeProcessors ...
func MakeProcessors(definitions []definition.Definition, flows map[string][]string) []processor.Processor {
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
		proc.SetStderr(os.Stderr)
		// ^^ append(os.Environ(), ...)
		processors[idx] = proc
		idx++
	}

	return processors
}