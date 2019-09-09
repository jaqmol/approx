package builtin

import (
	"bufio"

	"github.com/jaqmol/approx/definition"
	"github.com/jaqmol/approx/pipe"
)

// Merge ...
type Merge struct {
	def     definition.Definition
	stdins  []pipe.Reader
	stdout  *pipe.Writer
	stderr  *pipe.Writer
	running bool
}

// SetStdin ...
func (m *Merge) SetStdin(r *pipe.Reader) {
	m.stdins = append(m.stdins, *r)
}

// SetStdout ...
func (m *Merge) SetStdout(w *pipe.Writer) {
	m.stdout = w
}

// SetStderr ...
func (m *Merge) SetStderr(w *pipe.Writer) {
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
			go m.startReading(&stdin)
		}
		m.running = true
	}
}

// MakeMerge ...
func MakeMerge(def *definition.Definition) *Merge {
	return &Merge{
		def:    *def,
		stdins: make([]pipe.Reader, 0),
	}
}

func (m *Merge) startReading(aStdin *pipe.Reader) {
	scanner := bufio.NewScanner(aStdin)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		m.stdout.Channel() <- bytes
	}
}
