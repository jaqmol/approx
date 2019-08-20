package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jaqmol/approx/definition"
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

	rawFlow := findRawFlow(rawFormation)
	fmt.Printf("rawFlow: %v\n", rawFlow)

	definitions := definition.Parse(rawFormation)
	fmt.Printf("definitions: %v\n", definitions)
}

func formationFilePath() string {
	formationPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	return formationPath
}

func findRawFlow(rawFormation map[interface{}]interface{}) []string {
	for key, value := range rawFormation {
		if key == "Flow" {
			interfaceSlice := value.([]interface{})
			rawFlow := make([]string, len(interfaceSlice))
			for _, interfaceLine := range interfaceSlice {
				line := interfaceLine.(string)
				rawFlow = append(rawFlow, line)
			}
			return rawFlow
		}
	}
	return nil
}
