package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jaqmol/approx/assign"
	"github.com/jaqmol/approx/check"
	"github.com/jaqmol/approx/definition"
	"github.com/jaqmol/approx/flow"
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

	flows := flow.Parse(rawFormation)
	fmt.Printf("flows: %v\n", flows)

	definitions := definition.Parse(rawFormation)
	fmt.Printf("definitions: %v\n", definitions)

	processors := run.MakeProcessors(definitions, flows)
	pipes := run.MakePipes(processors, flows)

	assign.Variables(definitions)
	run.Connect(processors, flows, pipes)

	// TODO: Should start

	check.Check(definitions, flows)
}

func formationFilePath() string {
	formationPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	return formationPath
}
