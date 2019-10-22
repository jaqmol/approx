package test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/event"
	"github.com/jaqmol/approx/processor"
)

func TestSimpleCommandSequence(t *testing.T) {
	originals := loadTestData()[:10]
	originalForID := makePersonForIDMap(originals)
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)

	reader := bytes.NewReader(originalCombined)
	config := makeSimpleSequenceConfig()

	fork := processor.NewFork(&config.fork, reader)
	firstNameExtractCmd := processor.NewCommand(&config.firstNameExtract, fork.Outs()[0])
	lastNameExtractCmd := processor.NewCommand(&config.lastNameExtract, fork.Outs()[1])
	merge := processor.NewMerge(&config.merge, []io.Reader{
		firstNameExtractCmd.Out(),
		lastNameExtractCmd.Out(),
	})

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
			var extraction map[string]string
			err := json.Unmarshal(ob, &extraction)
			if err != nil {
				t.Fatalf("Couldn't unmarshall person from: \"%v\" -> %v\n", string(ob), err.Error())
			}

			original := originalForID[extraction["id"]]
			checkExtractedProp(t, original, extraction)

			businessIndex++
			loop = businessIndex != goal || loggingCounter != goal
		case eb := <-serializeLogMsgs:
			msg, err := event.UnmarshalLogMsg(eb)
			if err != nil {
				t.Fatal(err)
			}
			logMsgPntr, cmdErr, err := msg.PayloadOrError()
			if err != nil {
				t.Fatal(err)
			}

			if logMsgPntr != nil {
				logMsg := *logMsgPntr
				if strings.HasPrefix(logMsg, "Did extract \"first_name\"") {
					loggingCounter++
				} else if strings.HasPrefix(logMsg, "Did extract \"last_name\"") {
					loggingCounter++
				} else {
					log.Println("Unexpected log message:", logMsg)
				}
				loop = businessIndex != goal || loggingCounter != goal
			}
			if cmdErr != nil {
				t.Fatal(cmdErr.Error())
			}
		}
	}
}

func checkExtractedProp(t *testing.T, original Person, extraction map[string]string) {
	extractedValue, ok := extraction["first_name"]
	if ok {
		upperValue := strings.ToUpper(original.FirstName)
		if upperValue != extractedValue {
			t.Fatalf("Extracted value %v not as expected: %v", extractedValue, upperValue)
		}
	} else {
		extractedValue, ok = extraction["last_name"]
		if ok {
			upperValue := strings.ToUpper(original.LastName)
			if upperValue != extractedValue {
				t.Fatalf("Extracted value %v not as expected: %v", extractedValue, upperValue)
			}
		} else {
			t.Fatalf("Extraction not as expected: %v", extraction)
		}
	}
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
	}
	firstNameExtractConf := configuration.Command{
		Ident:    "extract-first-name",
		Cmd:      "node node-procs/test-extract-prop.js",
		Env:      []string{"PROP_NAME=first_name"},
		NextProc: &mergeConf,
	}
	lastNameExtractConf := configuration.Command{
		Ident:    "extract-last-name",
		Cmd:      "node node-procs/test-extract-prop.js",
		Env:      []string{"PROP_NAME=last_name"},
		NextProc: &mergeConf,
	}
	forkConf := configuration.Fork{
		Ident: "fork",
		NextProcs: []configuration.Processor{
			&firstNameExtractConf,
			&lastNameExtractConf,
		},
	}
	return &simpleSequenceConfig{
		fork:             forkConf,
		firstNameExtract: firstNameExtractConf,
		lastNameExtract:  lastNameExtractConf,
		merge:            mergeConf,
	}
}
