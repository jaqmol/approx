package test

import (
	"bytes"
	"io"
	"log"
	"testing"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/processor"
)

func TestSimpleCommandSequence(t *testing.T) {
	originals := loadTestData()[:10]
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)

	reader := bytes.NewReader(originalCombined)

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

	fork := processor.NewFork(&forkConf, reader)
	firstNameExtractCmd := processor.NewCommand(&firstNameExtractConf, fork.Outs()[0])
	lastNameExtractCmd := processor.NewCommand(&lastNameExtractConf, fork.Outs()[1])
	merge := processor.NewMerge(&mergeConf, []io.Reader{
		firstNameExtractCmd.Outs()[0],
		lastNameExtractCmd.Outs()[0],
	})

	// TODO: Continue with wiring the pipes and starting the processors

	log.Println(merge)

	// serializeOutput := outputSerializerChannel(command.Outs()[0])
	// serializeLogMsgs := outputSerializerChannel(command.Err())
	// command.Start()

	// goal := len(originals)
	// businessIndex := 0
	// loggingIndex := 0

	// loop := true
	// for loop {
	// 	select {
	// 	case ob := <-serializeOutput:
	// 		parsed, err := unmarshallPerson(ob)
	// 		if err != nil {
	// 			t.Fatalf("Couldn't unmarshall person from: \"%v\" -> %v\n", string(ob), err.Error())
	// 		}

	// 		original := originals[businessIndex]
	// 		checkFirstAndLastNames(t, &original, parsed)

	// 		businessIndex++
	// 		loop = businessIndex != goal || loggingIndex != goal
	// 	case eb := <-serializeLogMsgs:
	// 		msg, err := event.UnmarshalLogMsg(eb)
	// 		logMsg, cmdErr, err := msg.PayloadOrError()
	// 		if err != nil {
	// 			t.Fatal(err)
	// 		}
	// 		if logMsg != nil {
	// 			if strings.HasPrefix(*logMsg, "Did process:") {
	// 				loggingIndex++
	// 				loop = businessIndex != goal || loggingIndex != goal
	// 			}
	// 		}
	// 		if cmdErr != nil {
	// 			t.Fatal(cmdErr.Error())
	// 		}
	// 	}
	// }
}
