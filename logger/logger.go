package logger

import (
	"bytes"
	"io"
	"log"

	"github.com/jaqmol/approx/message"
)

// Logger ...
type Logger struct {
	writer       io.Writer
	serialize    chan []byte
	readersCount int
}

// NewLogger ...
func NewLogger(w io.Writer) *Logger {
	l := Logger{
		writer:    w,
		serialize: make(chan []byte),
	}
	return &l
}

// Start ...
func (l *Logger) Start() {
	for raw := range l.serialize {
		msg := bytes.Trim(raw, "\x00")
		line := append(msg, '\n')
		n, err := l.writer.Write(line)
		if err != nil {
			log.Fatalln(err.Error())
		}
		if n != len(line) {
			panic("Couldn't write complete line")
		}
	}
}

// Add ...
func (l *Logger) Add(r io.Reader) {
	go l.readFrom(r)
	l.readersCount++
}

func (l *Logger) readFrom(r io.Reader) {
	scanner := message.NewScanner(r)
	for scanner.Scan() {
		original := scanner.Bytes()
		toPassOn := make([]byte, len(original))
		copy(toPassOn, original)
		l.serialize <- toPassOn
	}
	l.readersCount--
	if l.readersCount == 0 {
		close(l.serialize)
	}
}
