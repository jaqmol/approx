package test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jaqmol/approx/event"

	"github.com/jaqmol/approx/processor"

	"github.com/jaqmol/approx/configuration"
)

// TestCommandWithBufferProcessing ...
func TestCommandWithBufferProcessing(t *testing.T) {
	performTestWithCmd(t, "node node-procs/test-buffer-processing.js")
}

// TestCommandWithJSONProcessing ...
func TestCommandWithJSONProcessing(t *testing.T) {
	performTestWithCmd(t, "node node-procs/test-json-processing.js")
}

func performTestWithCmd(t *testing.T, commandString string) {
	originals := loadTestData()
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)

	reader := bytes.NewReader(originalCombined)
	config := configuration.Command{
		Cmd: commandString,
	}

	command, err := processor.NewCommand(&config)
	if err != nil {
		t.Fatal(err)
	}
	err = command.Connect(reader)
	if err != nil {
		t.Fatal(err)
	}

	outputCollector, err := processor.NewCollector(command.Out())
	if err != nil {
		t.Fatal(err)
	}
	logMsgsCollector, err := processor.NewCollector(command.Err())
	if err != nil {
		t.Fatal(err)
	}

	command.Start()
	outputCollector.Start()
	logMsgsCollector.Start()

	goal := len(originals)
	businessIndex := 0
	loggingIndex := 0

	loop := true
	for loop {
		select {
		case ob := <-outputCollector.Events():
			parsed, err := unmarshallPerson(ob)
			if err != nil {
				t.Fatalf("Couldn't unmarshall person from: \"%v\" -> %v\n", string(ob), err.Error())
			}

			original := originals[businessIndex]
			checkFirstAndLastNames(t, &original, parsed)

			businessIndex++
			loop = businessIndex != goal || loggingIndex != goal
		case eb := <-logMsgsCollector.Events():
			msg, err := event.UnmarshalLogMsg(eb)
			logMsg, cmdErr, err := msg.PayloadOrError()
			if err != nil {
				t.Fatal(err)
			}
			if logMsg != nil {
				if strings.HasPrefix(*logMsg, "Did process:") {
					loggingIndex++
					loop = businessIndex != goal || loggingIndex != goal
				}
			}
			if cmdErr != nil {
				t.Fatal(cmdErr.Error())
			}
		}
	}
}

func checkFirstAndLastNames(t *testing.T, original, parsed *Person) {
	upperOrigFirstName := strings.ToUpper(original.FirstName)
	if upperOrigFirstName != parsed.FirstName {
		t.Fatalf("Expected uppercase first name %v, but got: %v", upperOrigFirstName, parsed.FirstName)
	}

	upperOrigLastName := strings.ToUpper(original.LastName)
	if upperOrigLastName != parsed.LastName {
		t.Fatalf("Expected uppercase last name %v, but got: %v", upperOrigLastName, parsed.LastName)
	}
}
