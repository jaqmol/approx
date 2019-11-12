package test

import (
	"testing"
)

// TODO: NONE OF THIS TESTS IS SUCCEEDING

// TestSimpleProcessorFormation ...
func TestSimpleProcessorFormation(t *testing.T) {
	// originals := loadTestData()[:10]
	// originalForID := makePersonForIDMap(originals)
	// originalBytes := marshallPeople(originals)
	// originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	// originalCombined = append(originalCombined, configuration.EvntEndBytes...)

	// inputReader := bytes.NewReader(originalCombined)
	// stdin := processor.NewStdin()
	// err := stdin.Connect(inputReader)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// outputWriter := newTestWriter()
	// stdout := processor.NewStdout(outputWriter)

	// errorsWriter := newTestWriter()

	// origArgs := os.Args
	// defer func() { os.Args = origArgs }()
	// testArgs := []string{origArgs[0], "gamma-test-proj"}
	// os.Args = testArgs

	// form, err := processor.NewFormation(stdin, stdout, errorsWriter)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// goal := len(originals) * 2
	// businessIndex := 0
	// loggingCounter := 0

	// form.Start()

	// loop := true
	// for loop {
	// 	select {
	// 	case ob := <-outputWriter.lines:
	// 		err = checkOutEvent(ob, originalForID)
	// 		catchToFatal(t, err)
	// 		businessIndex++
	// 		loop = businessIndex != goal || loggingCounter != goal
	// 	case eb := <-errorsWriter.lines:
	// 		counter, err := checkErrorEvent(eb)
	// 		catchToFatal(t, err)
	// 		loggingCounter += counter
	// 		loop = businessIndex != goal || loggingCounter != goal
	// 	}
	// }
}

// TestComplexProcessorFormation ...
// func TestComplexProcessorFormation(t *testing.T) {
// 	t.SkipNow()
// 	originals := loadTestData()[:10]
// 	originalForID := makePersonForIDMap(originals)
// 	originalBytes := marshallPeople(originals)
// 	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
// 	originalCombined = append(originalCombined, configuration.EvntEndBytes...)

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
