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
	"github.com/jaqmol/approx/pipe"
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
	showHelpAndExitIfNeeded()

	formationPath := formationFilePath()
	formationBytes, err := ioutil.ReadFile(formationPath)
	if err != nil {
		log.Fatal(err)
	}

	rawFormation := make(map[interface{}]interface{})

	err = yaml.Unmarshal(formationBytes, &rawFormation)
	if err != nil {
		log.Fatalln(err.Error())
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
	listenForLogEntries(stdErrs)
}

func showHelpAndExitIfNeeded() {
	args := os.Args[1:]
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("APPROX HELP:")
			fmt.Println("--help        | -h      Show this help.")
			fmt.Println("--json-output | -jo     Output log messages as JSON.")
			os.Exit(0)
		}
	}
}

func formationFilePath() string {
	formationPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	return formationPath
}

func listenForLogEntries(stdErrs map[string]pipe.Pipe) {
	logChan := make(chan message.SourcedLogEntry)
	for procName, errPipe := range stdErrs {
		go listenForLogEntry(logChan, procName, errPipe.Reader)
	}
	for logMsg := range logChan {
		logType := message.LogEntryTypeForString[logMsg.Cmd]
		logMsg.WriteTo(os.Stderr)
		if logType == message.Exit {
			os.Exit(-1)
		}
	}
}

func listenForLogEntry(logChan chan<- message.SourcedLogEntry, procName string, errReader io.Reader) {
	scanner := bufio.NewScanner(errReader)
	for scanner.Scan() {
		errBytes := scanner.Bytes()
		var msg message.Message
		err := json.Unmarshal(errBytes, &msg)
		if err != nil {
			strErrMsg := message.MakeSourcedLogEntry(procName, "", message.Fail, string(errBytes))
			logChan <- *strErrMsg
		} else {
			sourcedMsg := msg.ToSourcedLogEntry(procName)
			logChan <- *sourcedMsg
		}
	}
}
