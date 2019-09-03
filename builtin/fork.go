package builtin

import (
	"bufio"
	"encoding/json"
	"io"

	"github.com/jaqmol/approx/definition"
	"github.com/jaqmol/approx/message"
)

// Fork ...
type Fork struct {
	def        definition.Definition
	stdin      io.Reader
	stdouts    []io.Writer
	stderr     io.Writer
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
func (f *Fork) SetStdin(r io.Reader) {
	f.stdin = r
}

// SetStdout ...
func (f *Fork) SetStdout(w io.Writer) {
	f.stdouts = append(f.stdouts, w)
}

// SetStderr ...
func (f *Fork) SetStderr(w io.Writer) {
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
		stdouts:    make([]io.Writer, 0),
		distribute: distribute,
	}
}

func (f *Fork) start() {
	scanner := bufio.NewScanner(f.stdin)
	for scanner.Scan() {
		inputeBytes := scanner.Bytes()

		var msg message.Message
		err := json.Unmarshal(inputeBytes, &msg)

		if err != nil {
			message.WriteLogEntry(f.stderr, message.Fail, "", err.Error())
		} else {
			f.writeDistribute(&msg)
		}
	}
}

func (f *Fork) writeDistribute(msg *message.Message) {
	switch f.distribute {
	case distributeCopy:
		for i, stdout := range f.stdouts {
			msg.Index = &i
			f.write(stdout, msg)
		}
	case distributeCycle:
		stdout := f.stdouts[f.cycleIndex]
		msg.Index = &f.cycleIndex
		f.write(stdout, msg)
		f.cycleIndex++
		if f.cycleIndex >= len(f.stdouts) {
			f.cycleIndex = 0
		}
	}
}

func (f *Fork) write(stdout io.Writer, msg *message.Message) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		message.WriteLogEntry(f.stderr, message.Fail, msg.ID, err.Error())
	} else {
		bytes = append(bytes, []byte("\n")...)
		_, err = stdout.Write(bytes)
		if err != nil {
			message.WriteLogEntry(f.stderr, message.Fail, msg.ID, err.Error())
		}
	}
}
