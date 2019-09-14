package builtin

import (
	"bufio"

	"github.com/jaqmol/approx/channel"
	"github.com/jaqmol/approx/definition"
)

// Fork ...
type Fork struct {
	def        definition.Definition
	stdin      channel.Reader
	stdouts    []channel.Writer
	stderr     channel.Writer
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
func (f *Fork) SetStdin(r channel.Reader) {
	f.stdin = r
}

// SetStdout ...
func (f *Fork) SetStdout(w channel.Writer) {
	f.stdouts = append(f.stdouts, w)
}

// SetStderr ...
func (f *Fork) SetStderr(w channel.Writer) {
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
		stdouts:    make([]channel.Writer, 0),
		distribute: distribute,
	}
}

func (f *Fork) start() {
	wrap := channel.NewReaderWrap(f.stdin)
	scanner := bufio.NewScanner(wrap)
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
			stdout.Write() <- msgBytes
		}
	case distributeCycle:
		stdout := f.stdouts[f.cycleIndex]
		stdout.Write() <- msgBytes
		f.cycleIndex++
		if f.cycleIndex >= len(f.stdouts) {
			f.cycleIndex = 0
		}
	}
}
