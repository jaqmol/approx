package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jaqmol/approx/assign"
	"github.com/jaqmol/approx/check"
	"github.com/jaqmol/approx/definition"
	"github.com/jaqmol/approx/env"
	"github.com/jaqmol/approx/flow"
	"github.com/jaqmol/approx/message"
	"github.com/jaqmol/approx/run"
	"gopkg.in/yaml.v2"
)

// Definition ...
type Definition struct {
	Type    DefinitionType
	Name    string
	Assign  map[string]string
	Env     map[string]string
	Command string
}

// DefinitionType ...
type DefinitionType int

// DefinitionTypes
const (
	DefTypeHTTPServer DefinitionType = iota
	DefTypeFork
	DefTypeMerge
	DefTypeProcess
)

func main() {
	formationPath := formationFilePath()
	formationBytes, err := ioutil.ReadFile(formationPath)
	if err != nil {
		log.Fatal(err)
	}

	rawFormation := make(map[interface{}]interface{})

	err = yaml.Unmarshal(formationBytes, &rawFormation)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	definitions := definition.Parse(rawFormation)
	env.AugmentMissing(definitions)                    // 1. order is important
	assign.ResolveVariables(rawFormation, definitions) // 2. order is important

	flows := flow.Parse(rawFormation)
	check.Check(definitions, flows)

	processors := run.MakeProcessors(definitions, flows)
	pipes := run.MakePipes(definitions, flows)
	stdErrs := run.MakeStderrs(definitions)

	run.Connect(processors, flows, pipes, stdErrs)
	run.Start(processors)
	listenForErrorMessages(stdErrs)
}

func formationFilePath() string {
	formationPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	return formationPath
}

func listenForErrorMessages(stdErrs map[string]run.Pipe) {
	errChan := make(chan message.SourcedErrorMessage)
	for procName, errPipe := range stdErrs {
		go listenForErrorMessage(errChan, procName, errPipe.Reader)
	}
	for errMsg := range errChan {
		errType := message.ErrorTypeForString[errMsg.Cmd]
		errMsg.WriteTo(os.Stderr)
		if errType == message.Exit {
			os.Exit(-1)
		}
	}
}

func listenForErrorMessage(errChan chan<- message.SourcedErrorMessage, procName string, errReader io.Reader) {
	scanner := bufio.NewScanner(errReader)
	for scanner.Scan() {
		errBytes := scanner.Bytes()
		var msg *message.Message
		err := json.Unmarshal(errBytes, msg)
		if err != nil {
			errStr := fmt.Sprintf("Error parsing error-message from processor %v: %v", procName, err.Error())
			msg = message.NewError(message.Exit, "", errStr)
		}
		sourcedMsg := msg.ToSourcedErrorMessage(procName)
		errChan <- *sourcedMsg
	}
}
