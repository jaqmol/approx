package actor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/jaqmol/approx/config"
	"github.com/jaqmol/approx/logging"
	"github.com/jaqmol/approx/test"
)

func TestSimpleCommandSequence(t *testing.T) {
	// t.SkipNow()
	originals := test.LoadTestData() // [:10]
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, config.EvntEndBytes)
	originalCombined = append(originalCombined, config.EvntEndBytes...)

	conf := test.MakeSimpleSequenceConfig()
	producer := NewThrottledProducer(10, 5000)

	fork := NewForkFromConf(10, &conf.Fork)
	producer.Next(fork)

	firstNameExtractCmd, err := NewCommandFromConf(10, &conf.FirstNameExtract)
	catchToFatal(t, err)
	lastNameExtractCmd, err := NewCommandFromConf(10, &conf.LastNameExtract)
	catchToFatal(t, err)

	receiver := make(chan unifiedMessage, 10)

	// Logger
	logReceiver := make(chan []byte, 10)
	logger := logging.NewChannelLog(logReceiver)
	funnelIntoUnifiedLogMessages(logReceiver, receiver)
	logger.Add(firstNameExtractCmd.Logging())
	logger.Add(lastNameExtractCmd.Logging())
	// /Logger

	fork.Next(firstNameExtractCmd, lastNameExtractCmd)

	merge := NewMergeFromConf(10, &conf.Merge)
	firstNameExtractCmd.Next(merge)
	lastNameExtractCmd.Next(merge)

	collector := NewCollector(10)
	merge.Next(collector)

	startCollectingUnifiedDataMessages(t, collector, receiver, func() {
		close(receiver)
	})

	fork.Start()
	firstNameExtractCmd.Start()
	lastNameExtractCmd.Start()
	merge.Start()
	go logger.Start()

	startProducingTestMessages(t, producer, originalCombined)

	counter := 0
	for message := range receiver {
		if message.messageType == unifiedMsgDataType {
			err = checkOutEvent(message.data, originalForID)
			catchToFatal(t, err)
			counter++
		} else if message.messageType == unifiedMsgLogType {
			checkCommandLogMsg(t, "Did extract", message.data)
		}
	}

	goal := len(originals) * 2
	if counter != goal {
		t.Fatalf("%v messages expected to be produced, but got %v", goal, counter)
	}
}

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

// func checkErrorEvent(eb []byte) (counter int, err error) {
// 	msg, err := event.UnmarshalLogMsg(eb)
// 	if err != nil {
// 		return
// 	}
// 	logMsgPntr, cmdErr, err := msg.PayloadOrError()
// 	if err != nil {
// 		return
// 	}
//
// 	if logMsgPntr != nil {
// 		logMsg := *logMsgPntr
// 		if strings.HasPrefix(logMsg, "Did extract \"first_name\"") {
// 			counter++
// 		} else if strings.HasPrefix(logMsg, "Did extract \"last_name\"") {
// 			counter++
// 		} else {
// 			log.Println("Unexpected log message:", logMsg)
// 		}
// 	}
// 	if cmdErr != nil {
// 		err = cmdErr
// 	}
// 	return
// }
