package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jaqmol/approx/assign"
	"github.com/jaqmol/approx/check"
	"github.com/jaqmol/approx/definition"
	"github.com/jaqmol/approx/env"
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

	definitions := definition.Parse(rawFormation)
	env.AugmentMissing(definitions)                    // 1. order is important
	assign.ResolveVariables(rawFormation, definitions) // 2. order is important

	flows := flow.Parse(rawFormation)
	check.Check(definitions, flows)

	processors := run.MakeProcessors(definitions, flows)
	pipes := run.MakePipes(processors, flows)

	run.Connect(processors, flows, pipes)
	run.Start(processors) // TODO
}

func formationFilePath() string {
	formationPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	return formationPath
}
