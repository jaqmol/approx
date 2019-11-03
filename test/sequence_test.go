package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/event"
	"github.com/jaqmol/approx/processor"
)

// TODO: Add complex with Stdin and Stdout

func TestSimpleCommandSequence(t *testing.T) {
	originals := loadTestData() // [:10]
	originalForID := makePersonForIDMap(originals)
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)

	reader := bytes.NewReader(originalCombined)
	config := makeSimpleSequenceConfig()

	fork, err := processor.NewFork(&config.fork)
	catchToFatal(t, err)
	err = fork.Connect(reader)
	catchToFatal(t, err)

	firstNameExtractCmd, err := processor.NewCommand(&config.firstNameExtract)
	catchToFatal(t, err)
	err = firstNameExtractCmd.Connect(fork.Outs()[0])
	catchToFatal(t, err)

	lastNameExtractCmd, err := processor.NewCommand(&config.lastNameExtract)
	catchToFatal(t, err)
	err = lastNameExtractCmd.Connect(fork.Outs()[1])
	catchToFatal(t, err)

	merge, err := processor.NewMerge(&config.merge)
	catchToFatal(t, err)
	err = merge.Connect(firstNameExtractCmd.Out(), lastNameExtractCmd.Out())
	catchToFatal(t, err)

	serializeOutput := outputSerializerChannel(merge.Out())
	serializeLogMsgs := outputsSerializerChannel([]io.Reader{
		firstNameExtractCmd.Err(),
		lastNameExtractCmd.Err(),
	})

	fork.Start()
	firstNameExtractCmd.Start()
	lastNameExtractCmd.Start()
	merge.Start()

	goal := len(originals) * 2
	businessIndex := 0
	loggingCounter := 0

	loop := true
	for loop {
		select {
		case ob := <-serializeOutput:
			err = checkOutoutEvent(ob, originalForID)
			catchToFatal(t, err)
			businessIndex++
			loop = businessIndex != goal || loggingCounter != goal
		case eb := <-serializeLogMsgs:
			counter, err := checkErrorEvent(eb)
			catchToFatal(t, err)
			loggingCounter += counter
			loop = businessIndex != goal || loggingCounter != goal
		}
	}
}

func checkOutoutEvent(ob []byte, originalForID map[string]Person) error {
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

func checkExtractedPerson(original Person, extraction map[string]string) (err error) {
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

type simpleSequenceConfig struct {
	fork             configuration.Fork
	firstNameExtract configuration.Command
	lastNameExtract  configuration.Command
	merge            configuration.Merge
}

func makeSimpleSequenceConfig() *simpleSequenceConfig {
	mergeConf := configuration.Merge{
		Ident: "merge",
		Count: 2,
	}
	firstNameExtractConf := configuration.Command{
		Ident: "extract-first-name",
		Cmd:   "node node-procs/test-extract-prop.js",
		Env:   []string{"PROP_NAME=first_name"},
	}
	lastNameExtractConf := configuration.Command{
		Ident: "extract-last-name",
		Cmd:   "node node-procs/test-extract-prop.js",
		Env:   []string{"PROP_NAME=last_name"},
	}
	forkConf := configuration.Fork{
		Ident: "fork",
		Count: 2,
	}
	return &simpleSequenceConfig{
		fork:             forkConf,
		firstNameExtract: firstNameExtractConf,
		lastNameExtract:  lastNameExtractConf,
		merge:            mergeConf,
	}
}
