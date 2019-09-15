package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jaqmol/approx/assign"
	"github.com/jaqmol/approx/channel"
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

	procFlow, tappedPipeNames := flow.Parse(rawFormation)
	fmt.Printf("procFlow: %v\n", procFlow)
	fmt.Printf("tappedPipeNames: %v\n", tappedPipeNames)
	check.Check(definitions, procFlow)

	processors := run.MakeProcessors(definitions)
	pipes := run.MakePipes(definitions, procFlow, tappedPipeNames)
	stdErrs := run.MakeStderrs(definitions, tappedPipeNames)

	run.Connect(processors, procFlow, tappedPipeNames, pipes, stdErrs)
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

func listenForLogEntries(stdErrs map[string]channel.Pipe) {
	logChan := make(chan message.LogEntry)
	for procName, errPipe := range stdErrs {
		go listenForLogEntry(logChan, procName, errPipe)
	}
	for logMsg := range logChan {
		logMsg.WriteTo(os.Stderr)
	}
}

func listenForLogEntry(logChan chan<- message.LogEntry, procName string, errReader channel.Reader) {
	wrap := channel.NewReaderWrap(errReader)
	scanner := bufio.NewScanner(wrap)
	for scanner.Scan() {
		entryBytes := scanner.Bytes()
		entry := message.LogEntry{Source: procName, Message: string(entryBytes)}
		logChan <- entry
	}
}
