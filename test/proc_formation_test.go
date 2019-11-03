package test

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/processor"
)

// TestProcessorFormation ...
func TestProcessorFormation(t *testing.T) {
	// t.SkipNow()
	originals := loadTestData()[:10]
	originalForID := makePersonForIDMap(originals)
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)

	inputReader := bytes.NewReader(originalCombined)
	outputWriter := newTestWriter()
	errorsWriter := newTestWriter()

	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	testArgs := []string{origArgs[0], "complx-test-proj"}
	os.Args = testArgs

	form, err := processor.NewFormation()
	if err != nil {
		t.Fatal(err)
	}

	err = os.Setenv("APPROX_ENV", "development")
	if err != nil {
		t.Fatal(err)
	}

	err = processor.ChangeInterface(inputReader, outputWriter, errorsWriter)
	if err != nil {
		t.Fatal(err)
	}

	t.SkipNow()

	// TODO: This is not working.
	// Suggestion: Stdin and Stdout processors have not been tested in a sequence.

	form.Start()

	goal := len(originals) * 2
	businessIndex := 0
	loggingCounter := 0

	log.Println("Sarting loop")

	loop := true
	for loop {
		select {
		case ob := <-outputWriter.lines:
			err = checkOutoutEvent(ob, originalForID)
			catchToFatal(t, err)
			businessIndex++
			loop = businessIndex != goal || loggingCounter != goal
			log.Println("businessIndex:", businessIndex)
		case eb := <-errorsWriter.lines:
			counter, err := checkErrorEvent(eb)
			catchToFatal(t, err)
			loggingCounter += counter
			loop = businessIndex != goal || loggingCounter != goal
			log.Println("loggingCounter:", loggingCounter)
		}
	}

	log.Println("Finished loop")
}
