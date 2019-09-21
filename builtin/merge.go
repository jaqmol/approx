package builtin

import (
	"github.com/jaqmol/approx/channel"
	"github.com/jaqmol/approx/message"

	"github.com/jaqmol/approx/definition"
)

// Merge ...
type Merge struct {
	def     definition.Definition
	stdins  []channel.Reader
	stdout  channel.Writer
	stderr  channel.Writer
	running bool
}

// SetStdin ...
func (m *Merge) SetStdin(r channel.Reader) {
	m.stdins = append(m.stdins, r)
}

// SetStdout ...
func (m *Merge) SetStdout(w channel.Writer) {
	m.stdout = w
}

// SetStderr ...
func (m *Merge) SetStderr(w channel.Writer) {
	m.stderr = w
}

// Definition ...
func (m *Merge) Definition() *definition.Definition {
	return &m.def
}

// Start ...
func (m *Merge) Start() {
	if !m.running {
		for stdinIndex := range m.stdins {
			go m.startReading(stdinIndex)
		}
		m.running = true
	}
}

// MakeMerge ...
func MakeMerge(def *definition.Definition) *Merge {
	return &Merge{
		def:    *def,
		stdins: make([]channel.Reader, 0),
	}
}

func (m *Merge) startReading(stdinIndex int) {
	stdin := m.stdins[stdinIndex]
	envBuff := message.NewEnvelopeBuffer(stdin.Read())
	for env := range envBuff.Envelopes() {
		m.stdout.Write() <- env.Bytes
	}
}
