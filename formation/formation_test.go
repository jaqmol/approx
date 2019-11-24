package formation

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/jaqmol/approx/logging"

	"github.com/jaqmol/approx/config"
	"github.com/jaqmol/approx/test"
)

// TestSimpleActorFormation ...
func TestSimpleActorFormation(t *testing.T) {
	originals := test.LoadTestData() // [:100]
	// originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)
	originalCombined := bytes.Join(originalBytes, config.EvntEndBytes)
	originalCombined = append(originalCombined, config.EvntEndBytes...)

	// producer := actor.NewThrottledProducer(10, 5000)

	inputReader := bytes.NewReader(originalCombined)
	// stdin := processor.NewStdin()
	// err := stdin.Connect(inputReader)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	outputWriter := test.NewWriter()
	// stdout := processor.NewStdout(outputWriter)

	// errorsWriter := newTestWriter()

	projDir, err := filepath.Abs("../test/gamma-test-proj")
	if err != nil {
		t.Fatal(err)
	}

	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	testArgs := []string{origArgs[0], projDir}
	os.Args = testArgs

	logChannel := make(chan []byte)
	logger := logging.NewChannelLog(logChannel)
	form, err := NewFormation(inputReader, outputWriter, logger)
	if err != nil {
		t.Fatal(err)
	}

	finished := form.Start()

	// goal := len(originals) * 2
	// businessIndex := 0
	// loggingCounter := 0

	// form.Start()

	loop := true
	for loop {
		select {
		case outMsg := <-outputWriter.Lines:
			log.Println(string(outMsg))
			// err = test.CheckOutEvent(ob, originalForID)
			// catchToFatal(t, err)
			// businessIndex++
			// loop = businessIndex != goal || loggingCounter != goal
		case logMsg := <-logChannel:
			log.Println(string(logMsg))
		case <-finished:
			log.Println("FINISHED")
			loop = false
			// counter, err := checkErrorEvent(eb)
			// catchToFatal(t, err)
			// loggingCounter += counter
			// loop = businessIndex != goal || loggingCounter != goal
		}
	}
}

// TestComplexActorFormation ...
// func TestComplexActorFormation(t *testing.T) {
// 	t.SkipNow()
// 	originals := LoadTestData()[:10]
// 	originalForID := MakePersonForIDMap(originals)
// 	originalBytes := MarshalPeople(originals)
// 	originalCombined := bytes.Join(originalBytes, config.EvntEndBytes)
// 	originalCombined = append(originalCombined, config.EvntEndBytes...)

// 	inputReader := bytes.NewReader(originalCombined)
// 	stdin := processor.NewStdin()
// 	err := stdin.Connect(inputReader)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	outputWriter := newTestWriter()
// 	stdout := processor.NewStdout(outputWriter)

// 	errorsWriter := newTestWriter()

// 	origArgs := os.Args
// 	defer func() { os.Args = origArgs }()
// 	testArgs := []string{origArgs[0], "beta-test-proj"}
// 	os.Args = testArgs

// 	form, err := processor.NewFormation(stdin, stdout, errorsWriter)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	goal := len(originals) * 2
// 	businessIndex := 0
// 	loggingCounter := 0

// 	form.Start()

// 	loop := true
// 	for loop {
// 		select {
// 		case ob := <-outputWriter.lines:
// 			err = checkOutEvent(ob, originalForID)
// 			catchToFatal(t, err)
// 			businessIndex++
// 			loop = businessIndex != goal || loggingCounter != goal
// 		case eb := <-errorsWriter.lines:
// 			counter, err := checkErrorEvent(eb)
// 			catchToFatal(t, err)
// 			loggingCounter += counter
// 			loop = businessIndex != goal || loggingCounter != goal
// 		}
// 	}
// }
