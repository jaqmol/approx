package builtin

import (
	"bufio"
	"fmt"

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
	fmt.Printf("Merge stdins: %v\n", m.stdins)
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
		stdins: make([]pipe.Reader, 0),
	}
}

func (m *Merge) startReading(stdinIndex int) {
	stdin := m.stdins[stdinIndex]
	scanner := bufio.NewScanner(&stdin)
	for scanner.Scan() {
		msgBytes := scanner.Bytes()
		msgBytes = append(msgBytes, []byte("\n")...)
		// fmt.Printf("Merge did read: %v", string(msgBytes))
		m.stdout.Channel() <- msgBytes
	}
}
