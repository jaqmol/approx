package builtin

import (
	"bufio"

	"github.com/jaqmol/approx/definition"
	"github.com/jaqmol/approx/pipe"
)

// Fork ...
type Fork struct {
	def        definition.Definition
	stdin      *pipe.Reader
	stdouts    []pipe.Writer
	stderr     *pipe.Writer
	running    bool
	cycleIndex int
	distribute int
}

// ForkDistributeCopy ...
const ForkDistributeCopy = "copy"

// ForkDistributeCycle ...
const ForkDistributeCycle = "cycle"

const (
	distributeCopy = iota
	distributeCycle
)

// SetStdin ...
func (f *Fork) SetStdin(r *pipe.Reader) {
	f.stdin = r
}

// SetStdout ...
func (f *Fork) SetStdout(w *pipe.Writer) {
	f.stdouts = append(f.stdouts, *w)
}

// SetStderr ...
func (f *Fork) SetStderr(w *pipe.Writer) {
	f.stderr = w
}

// Definition ...
func (f *Fork) Definition() *definition.Definition {
	return &f.def
}

// Start ...
func (f *Fork) Start() {
	if !f.running {
		go f.start()
		f.running = true
	}
}

// MakeFork ...
func MakeFork(def *definition.Definition) *Fork {
	distributeStrPtr, ok := def.Env["DISTRIBUTE"]
	var distributeStr string
	if ok {
		distributeStr = *distributeStrPtr
	} else {
		distributeStr = ForkDistributeCopy
	}
	var distribute int
	switch distributeStr {
	case ForkDistributeCopy:
		distribute = distributeCopy
	case ForkDistributeCycle:
		distribute = distributeCycle
	}
	return &Fork{
		def:        *def,
		stdouts:    make([]pipe.Writer, 0),
		distribute: distribute,
	}
}

func (f *Fork) start() {
	scanner := bufio.NewScanner(f.stdin)
	for scanner.Scan() {
		msgBytes := scanner.Bytes()
		msgBytes = append(msgBytes, []byte("\n")...)
		f.writeDistribute(msgBytes)
	}
}

func (f *Fork) writeDistribute(msgBytes []byte) {
	switch f.distribute {
	case distributeCopy:
		for _, stdout := range f.stdouts {
			stdout.Channel() <- msgBytes
		}
	case distributeCycle:
		stdout := f.stdouts[f.cycleIndex]
		stdout.Channel() <- msgBytes
		f.cycleIndex++
		if f.cycleIndex >= len(f.stdouts) {
			f.cycleIndex = 0
		}
	}
}
