package test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jaqmol/approx/processor"

	"github.com/jaqmol/approx/configuration"
)

// TestCommand ...
func TestCommand(t *testing.T) {
	originals := loadTestData()
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)

	reader := bytes.NewReader(originalCombined)
	config := configuration.Command{
		Cmd: "node node-procs/test-buffer-processing.js",
	}

	command := processor.NewCommand(&config, reader)

	serialize := make(chan []byte)
	r := command.Outs()[0]
	go readFromReader(serialize, r)
	command.Start()

	goal := len(originals)
	index := 0

	for b := range serialize {
		parsed, err := unmarshallPerson(b)
		if err != nil {
			t.Fatalf("Couldn't unmarshall person from: \"%v\" -> %v\n", string(b), err.Error())
		}

		original := originals[index]
		checkFirstAndLastNames(t, &original, parsed)

		index++
		if index == goal {
			close(serialize)
		}
	}
}

func checkFirstAndLastNames(t *testing.T, original, parsed *TestPerson) {
	upperOrigFirstName := strings.ToUpper(original.FirstName)
	if upperOrigFirstName != parsed.FirstName {
		t.Fatalf("Expected uppercase first name %v, but got: %v", upperOrigFirstName, parsed.FirstName)
	}

	upperOrigLastName := strings.ToUpper(original.LastName)
	if upperOrigLastName != parsed.LastName {
		t.Fatalf("Expected uppercase last name %v, but got: %v", upperOrigLastName, parsed.LastName)
	}
}

// type readerToChannelTransformer struct {
// 	scanner *bufio.Scanner
// 	lines   chan []byte
// }

// func newReaderToChannelTransformer(reader io.Reader) *readerToChannelTransformer {
// 	return &readerToChannelTransformer{
// 		scanner: event.NewScanner(reader),
// 		lines:   make(chan []byte),
// 	}
// }

// func (t *readerToChannelTransformer) start() {
// 	for t.scanner.Scan() {
// 		raw := t.scanner.Bytes()
// 		data := bytes.Trim(raw, "\x00")
// 		cp := make([]byte, len(data))
// 		copy(cp, data)
// 		t.lines <- cp
// 	}

// 	close(t.lines)
// }
