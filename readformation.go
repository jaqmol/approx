package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// ReadFormation ...
func ReadFormation() {
	log.SetFlags(0)
	configFileData := readConfigFile()

	var configMap map[string]interface{}
	err := json.Unmarshal(configFileData, &configMap)
	if err != nil {
		log.Fatalf("Error parsing formation file: %v\n", err.Error())
	}

	privateProcs, publicProcs := createProcs(configMap)
	privateProcForName, inputPipePathForName, outputPipePathForName := prepareProcsOrExitOnUnspecified(
		privateProcs, publicProcs,
	)

	for _, p := range publicProcs {
		log.Printf("public proc: %v\n", *p.Name())
	}
	for n := range privateProcForName {
		log.Printf("private proc: %v\n", n)
	}
	for n, p := range inputPipePathForName {
		log.Printf("input pipe %v = %v\n", n, p)
	}
	for n, p := range outputPipePathForName {
		log.Printf("output pipe %v = %v\n", n, p)
	}
}

func prepareProcsOrExitOnUnspecified(privateProcs []Proc, publicProcs []Proc) (
	privateProcForName map[string]Proc, inputPipePathForName map[string]string, outputPipePathForName map[string]string,
) {
	privateProcForName = make(map[string]Proc)
	requiredPrivateInputNames := make(map[string]bool)
	requiredPrivateOutputNames := make(map[string]bool)

	for _, p := range privateProcs {
		pn := *p.Name()
		privateProcForName[pn] = p
		addPrivateProcNames(requiredPrivateInputNames, p.Inputs())
		addPrivateProcNames(requiredPrivateOutputNames, p.Outputs())
	}
	for _, p := range publicProcs {
		addPrivateProcNames(requiredPrivateInputNames, p.Inputs())
		addPrivateProcNames(requiredPrivateOutputNames, p.Outputs())
	}

	undefinedPrivateInputNames := findUnspecifiedProcs(privateProcForName, requiredPrivateInputNames)
	undefinedPrivateOutputNames := findUnspecifiedProcs(privateProcForName, requiredPrivateOutputNames)

	inputPipePathForName, undefinedPrivateInputNames = findEnvSpecifiedPipePaths("in", undefinedPrivateInputNames)
	outputPipePathForName, undefinedPrivateOutputNames = findEnvSpecifiedPipePaths("out", undefinedPrivateOutputNames)

	hasUnspecInputs := logUnspecifiedNames("private input", undefinedPrivateInputNames)
	hasUnspecOutputs := logUnspecifiedNames("private output", undefinedPrivateOutputNames)
	if hasUnspecInputs || hasUnspecOutputs {
		os.Exit(1)
	}

	return privateProcForName, inputPipePathForName, outputPipePathForName
}

func logUnspecifiedNames(info string, names []string) bool {
	if len(names) > 0 {
		log.Printf("Please define %v processes: %v\n", info, strings.Join(names, ", "))
		return true
	}
	return false
}

func findEnvSpecifiedPipePaths(prefix string, procNames []string) (map[string]string, []string) {
	specNames := make(map[string]string)
	unspecNames := make([]string, 0)
	for _, pn := range procNames {
		envName := strings.ToUpper(fmt.Sprintf("%v_%v", prefix, pn))
		if value, ok := os.LookupEnv(envName); ok {
			specNames[pn] = value
		} else {
			unspecNames = append(unspecNames, pn)
		}
	}
	return specNames, unspecNames
}

func findUnspecifiedProcs(definedProcs map[string]Proc, requiredNames map[string]bool) []string {
	acc := make([]string, 0)
	for rpn := range requiredNames {
		if _, ok := definedProcs[rpn]; !ok {
			acc = append(acc, rpn)
		}
	}
	return acc
}

func addPrivateProcNames(acc map[string]bool, procNames []string) {
	for _, n := range procNames {
		if strings.HasPrefix(n, "_") {
			acc[n] = true
		}
	}
}

func readConfigFile() []byte {
	if len(os.Args) < 2 {
		log.Fatalln("No formation file argument provided")
	}
	formationFilePath := os.Args[1]
	log.Printf("Loading app formation: %v", formationFilePath)
	configFileData, err := ioutil.ReadFile(formationFilePath)
	if err != nil {
		log.Fatalf("Error reading formation file: %v\n", err.Error())
	}
	return configFileData
}

func createProcs(configMap map[string]interface{}) ([]Proc, []Proc) {
	privateProcs := make([]Proc, 0)
	publicProcs := make([]Proc, 0)
	for specName, untyped := range configMap {
		isPrivate := strings.HasPrefix(specName, "_")
		if spec, ok := untyped.(map[string]interface{}); ok {
			dec := NewSpecDecoder(spec)
			specType := dec.String("type")
			if specType == nil {
				log.Fatalf("Spec \"%v\" has no type\n", specName)
			}
			var proc Proc
			switch *specType {
			case "http":
				proc = NewHTTPProc(specName, dec)
			case "fork":
				proc = NewForkProc(specName, dec)
			case "merge":
				proc = NewMergeProc(specName, dec)
			case "process":
				var err error
				proc, err = NewExtProc(specName, dec)
				if err != nil {
					log.Fatalln(err.Error())
				}
			}
			if proc != nil {
				if isPrivate {
					privateProcs = append(privateProcs, proc)
				} else {
					publicProcs = append(publicProcs, proc)
				}
			} else {

			}
		}
	}
	return privateProcs, publicProcs
}

// Proc ...
type Proc interface {
	RequiredProps() []string
	Outputs() []string
	Inputs() []string
	Name() *string
}
