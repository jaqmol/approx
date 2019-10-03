package logger

import (
	"io"
	"log"

	"github.com/jaqmol/approx/message"
)

// Logger ...
type Logger struct {
	writer     io.Writer
	serializer chan []byte
}

// NewLogger ...
func NewLogger(w io.Writer) *Logger {
	l := Logger{
		writer:     w,
		serializer: make(chan []byte),
	}
	return &l
}

// Start ...
func (l *Logger) Start() {
	for msg := range l.serializer {
		line := append(msg, '\n')
		_, err := l.writer.Write(line)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
}

// Add ...
func (l *Logger) Add(r io.Reader) {
	go l.readFrom(r)
}

func (l *Logger) readFrom(r io.Reader) {
	scanner := message.NewScanner(r)
	for scanner.Scan() {
		l.serializer <- scanner.Bytes()
	}
	close(l.serializer)
}
