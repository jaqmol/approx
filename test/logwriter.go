package testpackage

import (
	"bytes"
	"log"
)

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
