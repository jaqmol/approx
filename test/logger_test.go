package testpackage

import (
	"bytes"
	"log"
	"testing"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/logger"
)

// TestLogger ...
func TestLogger(t *testing.T) {
	originals := loadTestData()
	originalForID := makePersonForIDMap(originals)
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.MsgEndBytes)
	originalCombined = append(originalCombined, configuration.MsgEndBytes...)
	reader := bytes.NewReader(originalCombined)

	writer := newLogWriter()
	l := logger.NewLogger(writer)
	l.Add(reader)
	go l.Start()

	count := 0
	for b := range writer.lines {
		parsed, err := unmarshallPerson(b)
		if err != nil {
			t.Fatalf("Couldn't unmarshall person from: \"%v\"\n", string(b))
		}
		original := originalForID[parsed.ID]
		if !original.Equals(parsed) {
			t.Fatal("Parsed data doesn't conform to original")
		}
		count++
		writer.stop(count == len(originalBytes))
	}

	if len(originals) != count {
		t.Fatal("Logged line count doesn't corespond to received ones")
	}
}

type logWriter struct {
	lines   chan []byte
	running bool
}

func newLogWriter() *logWriter {
	return &logWriter{
		lines:   make(chan []byte),
		running: true,
	}
}

func (w *logWriter) Write(raw []byte) (int, error) {
	if w.running {
		b := bytes.Trim(raw, "\n\r")
		if len(b) == 0 {
			return len(raw), nil
		}
		l := make([]byte, len(b))
		copy(l, b)
		w.lines <- l
	} else {
		log.Fatalf("Writer stopped but received data: \"%v\"\n", string(raw))
	}
	return len(raw), nil
}

func (w *logWriter) stop(doStop bool) {
	if doStop {
		w.running = false
		close(w.lines)
	}
}
