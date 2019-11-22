package actor

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/jaqmol/approx/event"
	"github.com/jaqmol/approx/test"
)

func TestSimpleCommandSequence(t *testing.T) {
	t.SkipNow()
	// originals := test.LoadTestData()[:10]
	// originalForID := test.MakePersonForIDMap(originals)
	// originalBytes := test.MarshalPeople(originals)

	// originalCombined := bytes.Join(originalBytes, config.EvntEndBytes)
	// originalCombined = append(originalCombined, config.EvntEndBytes...)

	// conf := test.MakeSimpleSequenceConfig()

	// producer := NewProducer(10)

	// fork := NewFork(10, conf.Fork.Ident, conf.Fork.Count)
	// producer.Next(fork)

	// firstNameExtractCmd, err := NewCommandFromConf(10, &conf.FirstNameExtract)
	// catchToFatal(t, err)
	// lastNameExtractCmd, err := NewCommandFromConf(10, &conf.LastNameExtract)
	// catchToFatal(t, err)

	// var firstNameCmdLogBuffer bytes.Buffer
	// firstNameExtractCmd.Logging(&firstNameCmdLogBuffer)
	// var lastNameCmdLogBuffer bytes.Buffer
	// lastNameExtractCmd.Logging(&lastNameCmdLogBuffer)

	// fork.Next(firstNameExtractCmd, lastNameExtractCmd)

	// merge := NewMerge(10, conf.Merge.Ident, conf.Merge.Count)
	// firstNameExtractCmd.Next(merge)
	// lastNameExtractCmd.Next(merge)

	// collector := NewCollector(10)
	// merge.Next(collector)
	// receiver := make(chan []byte, 10)

	// startCollectingTestMessages(t, collector, receiver, func() {
	// 	close(receiver)
	// })

	// fork.Start()
	// firstNameExtractCmd.Start()
	// lastNameExtractCmd.Start()
	// merge.Start()

	// startProducingTestMessages(t, producer, originalCombined)

	// // goal := len(originals) * 2
	// counter := 0
	// for message := range receiver {
	// 	err = checkOutEvent(message, originalForID)
	// 	catchToFatal(t, err)
	// 	counter++
	// }

	// // TODO: CHECK COUNTER

	// // TODO: CHECK ALL LOGGING EVENTS IN BOTH BUFFERS:
	// // - firstNameCmdLogBuffer
	// // - lastNameCmdLogBuffer
	// //
	// // counter, err := checkErrorEvent(logEventBytes)
	// // catchToFatal(t, err)
	// // loggingCounter += counter
}

// func TestComplexCommandSequence(t *testing.T) {
// 	// t.SkipNow()
// 	originals := LoadTestData()[:10]
// 	originalForID := MakePersonForIDMap(originals)
// 	originalBytes := MarshalPeople(originals)

// 	originalCombined := bytes.Join(originalBytes, config.EvntEndBytes)
// 	originalCombined = append(originalCombined, config.EvntEndBytes...)

// 	reader := bytes.NewReader(originalCombined)
// 	conf := makeSimpleSequenceConfig()

// 	stdin := processor.NewStdin()
// 	err := stdin.Connect(reader)
// 	catchToFatal(t, err)

// 	fork, err := processor.NewFork(&conf.fork)
// 	catchToFatal(t, err)
// 	err = fork.Connect(stdin.Out())
// 	catchToFatal(t, err)

// 	firstNameExtractCmd, err := processor.NewCommand(&conf.firstNameExtract)
// 	catchToFatal(t, err)
// 	err = firstNameExtractCmd.Connect(fork.Outs()[0])
// 	catchToFatal(t, err)

// 	lastNameExtractCmd, err := processor.NewCommand(&conf.lastNameExtract)
// 	catchToFatal(t, err)
// 	err = lastNameExtractCmd.Connect(fork.Outs()[1])
// 	catchToFatal(t, err)

// 	merge, err := processor.NewMerge(&conf.merge)
// 	catchToFatal(t, err)
// 	err = merge.Connect(firstNameExtractCmd.Out(), lastNameExtractCmd.Out())
// 	catchToFatal(t, err)

// 	stdout := processor.NewStdout()
// 	err = stdout.Connect(merge.Out())
// 	catchToFatal(t, err)

// 	outputCollector, err := processor.NewCollector(stdout.Out())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// serializeOutput := outputSerializerChannel(stdout.Out())
// 	logMsgsCollector, err := processor.NewCollector(
// 		firstNameExtractCmd.Err(),
// 		lastNameExtractCmd.Err(),
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	outputCollector.Start()
// 	logMsgsCollector.Start()
// 	// serializeLogMsgs := outputsSerializerChannel([]io.Reader{
// 	// 	firstNameExtractCmd.Err(),
// 	// 	lastNameExtractCmd.Err(),
// 	// })

// 	fork.Start()
// 	firstNameExtractCmd.Start()
// 	lastNameExtractCmd.Start()
// 	merge.Start()

// 	goal := len(originals) * 2
// 	businessIndex := 0
// 	loggingCounter := 0

// 	loop := true
// 	for loop {
// 		select {
// 		case ob := <-outputCollector.Events():
// 			err = checkOutEvent(ob, originalForID)
// 			catchToFatal(t, err)
// 			businessIndex++
// 			loop = businessIndex != goal || loggingCounter != goal
// 		case eb := <-logMsgsCollector.Events():
// 			counter, err := checkErrorEvent(eb)
// 			catchToFatal(t, err)
// 			loggingCounter += counter
// 			loop = businessIndex != goal || loggingCounter != goal
// 		}
// 	}
// }

func checkOutEvent(ob []byte, originalForID map[string]test.Person) error {
	var extraction map[string]string
	err := json.Unmarshal(ob, &extraction)
	if err != nil {
		return fmt.Errorf("Couldn't unmarshall person from: \"%v\" -> %v", string(ob), err.Error())
	}

	original := originalForID[extraction["id"]]
	return checkExtractedPerson(original, extraction)
}

func catchToFatal(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func checkExtractedPerson(original test.Person, extraction map[string]string) (err error) {
	extractedValue, ok := extraction["first_name"]
	if ok {
		upperValue := strings.ToUpper(original.FirstName)
		if upperValue != extractedValue {
			err = fmt.Errorf("Extracted value %v not as expected: %v", extractedValue, upperValue)
		}
	} else {
		extractedValue, ok = extraction["last_name"]
		if ok {
			upperValue := strings.ToUpper(original.LastName)
			if upperValue != extractedValue {
				err = fmt.Errorf("Extracted value %v not as expected: %v", extractedValue, upperValue)
			}
		} else {
			err = fmt.Errorf("Extraction not as expected: %v", extraction)
		}
	}
	return
}

func checkErrorEvent(eb []byte) (counter int, err error) {
	msg, err := event.UnmarshalLogMsg(eb)
	if err != nil {
		return
	}
	logMsgPntr, cmdErr, err := msg.PayloadOrError()
	if err != nil {
		return
	}

	if logMsgPntr != nil {
		logMsg := *logMsgPntr
		if strings.HasPrefix(logMsg, "Did extract \"first_name\"") {
			counter++
		} else if strings.HasPrefix(logMsg, "Did extract \"last_name\"") {
			counter++
		} else {
			log.Println("Unexpected log message:", logMsg)
		}
	}
	if cmdErr != nil {
		err = cmdErr
	}
	return
}
