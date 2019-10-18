package test

import (
	"bufio"
	"bytes"
	"io"
	"testing"

	"github.com/jaqmol/approx/processor"

	"github.com/jaqmol/approx/event"

	"github.com/jaqmol/approx/configuration"
)

// TestCommand ...
func TestCommand(t *testing.T) {
	originals := loadTestData()[:10]
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)

	reader := bytes.NewReader(originalCombined)
	config := configuration.Command{
		Cmd: "node node-procs/cmd-2.js",
	}

	command := processor.NewCommand(&config, reader)
	command.Start()
	command.Wait()

	// transformer := newReaderToChannelTransformer(reader)
	// go transformer.start()
	// for data := range transformer.lines {
	// 	log.Println(string(data))
	// }
}

type readerToChannelTransformer struct {
	scanner *bufio.Scanner
	lines   chan []byte
}

func newReaderToChannelTransformer(reader io.Reader) *readerToChannelTransformer {
	return &readerToChannelTransformer{
		scanner: event.NewScanner(reader),
		lines:   make(chan []byte),
	}
}

func (t *readerToChannelTransformer) start() {
	for t.scanner.Scan() {
		raw := t.scanner.Bytes()
		data := bytes.Trim(raw, "\x00")
		cp := make([]byte, len(data))
		copy(cp, data)
		t.lines <- cp
	}

	close(t.lines)
}
