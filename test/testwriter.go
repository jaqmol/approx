package test

import (
	"bytes"
	"log"
)

type testWriter struct {
	lines   chan []byte
	running bool
}

func newTestWriter() *testWriter {
	return &testWriter{
		lines:   make(chan []byte),
		running: true,
	}
}

func (w *testWriter) Write(raw []byte) (int, error) {
	// log.Println("Test writer was written to", len(raw), "bytes") TODO: REMOVE
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

func (w *testWriter) stop(doStop bool) {
	if doStop {
		w.running = false
		close(w.lines)
	}
}
