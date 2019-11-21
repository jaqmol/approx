package actor

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/test"
)

func TestSingleCommand(t *testing.T) {
	originals := test.LoadTestData()
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)

	producer := NewProducer(10)
	command := NewCommand(10, "cmd", "node", "node-procs/test-buffer-processing.js")

	testDir, err := filepath.Abs("../test")
	if err != nil {
		t.Fatal(err)
	}
	command.Directory(testDir)

	var logBuffer bytes.Buffer
	command.Logging(&logBuffer)

	collector := NewCollector(10)
	receiver := make(chan []byte, 10)

	producer.Next(command)
	command.Next(collector)

	startCollectingTestMessages(t, collector, receiver, func() {
		close(receiver)
	})
	command.Start()
	startProducingTestMessages(t, producer, originalCombined)

	counter := 0
	for message := range receiver {
		parsed, err := test.UnmarshalPerson(message)
		if err != nil {
			t.Fatal(err)
		}

		original := originalForID[parsed.ID]

		test.CheckUpperFirstAndLastNames(t, &original, parsed)
		counter++
	}

	if counter != len(originals) {
		t.Fatalf("%v messages expected to be produced, but got %v", len(originals), counter)
	}
}
