package run

import (
	"os"
	"os/exec"
	"strings"

	"github.com/jaqmol/approx/builtin"
	"github.com/jaqmol/approx/builtin/httpserver"
	"github.com/jaqmol/approx/channel"
	"github.com/jaqmol/approx/definition"
	"github.com/jaqmol/approx/processor"
)

// Process ...
type Process struct {
	cmd     exec.Cmd
	def     definition.Definition
	running bool
}

// SetStdin ...
func (p *Process) SetStdin(r channel.Reader) {
	wrap := channel.NewReaderWrap(r)
	p.cmd.Stdin = wrap
}

// SetStdout ...
func (p *Process) SetStdout(w channel.Writer) {
	wrap := channel.NewWriterWrap(w)
	p.cmd.Stdout = wrap
}

// SetStderr ...
func (p *Process) SetStderr(w channel.Writer) {
	wrap := channel.NewWriterWrap(w)
	p.cmd.Stderr = wrap
}

// Definition ...
func (p *Process) Definition() *definition.Definition {
	return &p.def
}

// Start ...
func (p *Process) Start() {
	if !p.running {
		go p.start()
		p.running = true
	}
}

func (p *Process) start() {
	err := p.cmd.Start()
	if err != nil {
		panic(err)
	}
	err = p.cmd.Wait()
	if err != nil {
		panic(err)
	}
}

// MakeProcess ...
func MakeProcess(def *definition.Definition) *Process {
	cmd, args := commandComponents(def.Command)
	proc := Process{
		cmd: *exec.Command(cmd, args...),
		def: *def,
	}
	proc.cmd.Env = append(os.Environ(), def.EnvSlice()...)
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
			proc = httpserver.MakeHTTPServer(&def)
		case definition.TypeFork:
			proc = builtin.MakeFork(&def)
		case definition.TypeMerge:
			proc = builtin.MakeMerge(&def)
		case definition.TypeProcess:
			proc = MakeProcess(&def)
		}
		processors[idx] = proc
		idx++
	}

	return processors
}

func commandComponents(command string) (string, []string) {
	comps := make([]string, 0)
	rawComps := strings.Split(command, " ")
	for _, cmp := range rawComps {
		if len(cmp) > 0 {
			comps = append(comps, cmp)
		}
	}
	return comps[0], comps[1:]
}
