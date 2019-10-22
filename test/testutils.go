package test

import (
	"bytes"
	"io"

	"github.com/jaqmol/approx/event"
)

func outputSerializerChannel(output io.Reader) <-chan []byte {
	serializer := make(chan []byte)
	go readFromReader(serializer, output)
	return serializer
}

func readFromReader(serialize chan<- []byte, reader io.Reader) {
	scanner := event.NewScanner(reader)
	for scanner.Scan() {
		raw := scanner.Bytes()
		original := bytes.Trim(raw, "\x00")
		toPassOn := make([]byte, len(original))
		copy(toPassOn, original)
		serialize <- toPassOn
	}
}

func outputsSerializerChannel(outputs []io.Reader) <-chan []byte {
	serializer := make(chan []byte)
	for _, reader := range outputs {
		go readFromReader(serializer, reader)
	}
	return serializer
}
