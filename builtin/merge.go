package builtin

import (
	"bufio"
	"encoding/json"
	"io"

	"github.com/jaqmol/approx/definition"
	"github.com/jaqmol/approx/message"
)

// Merge ...
type Merge struct {
	def          definition.Definition
	stdins       []io.Reader
	stdout       io.Writer
	stderr       io.Writer
	running      bool
	stdoutWriter chan []byte
}

// SetStdin ...
func (m *Merge) SetStdin(r io.Reader) {
	m.stdins = append(m.stdins, r)
}

// SetStdout ...
func (m *Merge) SetStdout(w io.Writer) {
	m.stdout = w
}

// SetStderr ...
func (m *Merge) SetStderr(w io.Writer) {
	m.stderr = w
}

// Definition ...
func (m *Merge) Definition() *definition.Definition {
	return &m.def
}

// Start ...
func (m *Merge) Start() {
	if !m.running {
		for _, stdin := range m.stdins {
			go m.startReading(stdin)
		}
		go m.startWriting()
		m.running = true
	}
}

// MakeMerge ...
func MakeMerge(def *definition.Definition) *Merge {
	return &Merge{
		def:          *def,
		stdins:       make([]io.Reader, 0),
		stdoutWriter: make(chan []byte),
	}
}

func (m *Merge) startReading(aStdin io.Reader) {
	scanner := bufio.NewScanner(aStdin)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		m.stdoutWriter <- bytes
	}
}

func (m *Merge) startWriting() {
	for bytes := range m.stdoutWriter {
		_, err := m.stdout.Write(bytes)
		if err != nil {
			var msg message.Message
			err2 := json.Unmarshal(bytes, &msg)
			if err2 != nil {
				message.WriteLogEntry(m.stderr, message.Fail, "", err.Error())
			} else {
				message.WriteLogEntry(m.stderr, message.Fail, msg.ID, err.Error())
			}
		}
	}
}
