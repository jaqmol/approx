package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// ReadFormation ...
func ReadFormation() *Formation {
	log.SetFlags(0)
	configFileData := readConfigFile()

	var configMap map[string]interface{}
	err := json.Unmarshal(configFileData, &configMap)
	if err != nil {
		log.Fatalf("Error parsing formation file: %v\n", err.Error())
	}

	privateConfs, publicConfs := createConfs(configMap)
	publicInputPathForName, publicOutputPathForName :=
		findPublicPathsOrExitOnUnspecified(privateConfs, publicConfs)

	exitOnUnassignedVariables(privateConfs, publicConfs)

	return &Formation{
		PrivateConfs:            privateConfs,
		PublicConfs:             publicConfs,
		PublicInputPathForName:  publicInputPathForName,
		PublicOutputPathForName: publicOutputPathForName,
	}
}

// Formation ...
type Formation struct {
	PrivateConfs            []Conf
	PublicConfs             []Conf
	PublicInputPathForName  map[string]string
	PublicOutputPathForName map[string]string
}

func exitOnUnassignedVariables(privateConfs []Conf, publicConfs []Conf) {
	varsAssignedTo, varsAssignedFrom := make(map[string]bool), make(map[string]bool)
	collectVarsFromConfs(privateConfs, varsAssignedTo, varsAssignedFrom)
	collectVarsFromConfs(publicConfs, varsAssignedTo, varsAssignedFrom)
	unassignedVars := findUnassignedVars(varsAssignedTo, varsAssignedFrom)
	if len(unassignedVars) > 0 {
		acc := make([]string, 0)
		for k := range unassignedVars {
			acc = append(acc, k)
		}
		log.Fatalf("Please resolve assignment of variables: %v\n", strings.Join(acc, ", "))
	}
}

func collectVarsFromConfs(confs []Conf, varsAssignedTo map[string]bool, varsAssignedFrom map[string]bool) {
	for _, c := range confs {
		for k, v := range c.Assign() {
			if strings.HasPrefix(k, "$") {
				varsAssignedTo[k] = true
			}
			if strings.HasPrefix(v, "$") {
				varsAssignedFrom[v] = true
			}
		}
	}
	return
}

func findUnassignedVars(varsAssignedTo map[string]bool, varsAssignedFrom map[string]bool) (
	unassignedVars map[string]bool,
) {
	unassignedVars = make(map[string]bool)
	for to := range varsAssignedTo {
		didFind := false
		for from := range varsAssignedFrom {
			if from == to {
				didFind = true
				break
			}
		}
		if !didFind {
			unassignedVars[to] = true
		}
	}
	for from := range varsAssignedFrom {
		didFind := false
		for to := range varsAssignedTo {
			if to == from {
				didFind = true
				break
			}
		}
		if !didFind {
			unassignedVars[from] = true
		}
	}
	return
}

func findPublicPathsOrExitOnUnspecified(privateConfs []Conf, publicConfs []Conf) (
	publicInputPathForName map[string]string,
	publicOutputPathForName map[string]string,
) {
	privateConfForName := make(map[string]Conf)
	requiredPrivateInputNames := make(map[string]bool)
	requiredPrivateOutputNames := make(map[string]bool)

	for _, p := range privateConfs {
		pn := p.Name()
		privateConfForName[pn] = p
		addPrivateConfNames(requiredPrivateInputNames, p.Inputs())
		addPrivateConfNames(requiredPrivateOutputNames, p.Outputs())
	}
	for _, p := range publicConfs {
		addPrivateConfNames(requiredPrivateInputNames, p.Inputs())
		addPrivateConfNames(requiredPrivateOutputNames, p.Outputs())
	}

	undefinedPrivateInputNames := findUnspecifiedPrivateInputNames(privateConfForName, requiredPrivateInputNames)
	undefinedPrivateOutputNames := findUnspecifiedPrivateInputNames(privateConfForName, requiredPrivateOutputNames)

	hasUnspecInputs := logUnspecifiedNames("private input", undefinedPrivateInputNames)
	hasUnspecOutputs := logUnspecifiedNames("private output", undefinedPrivateOutputNames)
	if hasUnspecInputs || hasUnspecOutputs {
		os.Exit(1)
	}

	return
}

func keysFromStringBoolMap(aMap map[string]bool) (keys []string) {
	keys = make([]string, len(aMap))
	i := 0
	for k := range aMap {
		keys[i] = k
		i++
	}
	return
}

func logUnspecifiedNames(info string, names []string) bool {
	if len(names) > 0 {
		log.Printf("Please define %v processes: %v\n", info, strings.Join(names, ", "))
		return true
	}
	return false
}

func findEnvSpecifiedPipePaths(prefix string, confNames []string) (map[string]string, []string) {
	specNames := make(map[string]string)
	unspecNames := make([]string, 0)
	for _, pn := range confNames {
		envName := strings.ToUpper(fmt.Sprintf("%v_%v", prefix, pn))
		if value, ok := os.LookupEnv(envName); ok {
			specNames[pn] = value
		} else {
			unspecNames = append(unspecNames, pn)
		}
	}
	return specNames, unspecNames
}

func findUnspecifiedPrivateInputNames(definedConfs map[string]Conf, requiredNames map[string]bool) []string {
	acc := make([]string, 0)
	for rpn := range requiredNames {
		if _, ok := definedConfs[rpn]; !ok {
			acc = append(acc, rpn)
		}
	}
	return acc
}

func addPrivateConfNames(acc map[string]bool, confNames []string) {
	for _, n := range confNames {
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

func createConfs(configMap map[string]interface{}) ([]Conf, []Conf) {
	privateConfs := make([]Conf, 0)
	publicConfs := make([]Conf, 0)
	for specName, untyped := range configMap {
		isPrivate := strings.HasPrefix(specName, "_")
		if spec, ok := untyped.(map[string]interface{}); ok {
			dec := newSpecDec(spec)
			specType, ok := dec.string("type")
			if !ok {
				log.Fatalf("Spec \"%v\" has no type\n", specName)
			}
			var conf Conf
			var err error
			switch specType {
			case "http":
				conf, err = NewHTTPConf(specName, dec)
			case "fork":
				conf, err = NewForkConf(specName, dec)
			case "merge":
				conf, err = NewMergeConf(specName, dec)
			case "process":
				conf, err = NewProcessConf(specName, dec)
			}
			if err != nil {
				log.Fatalln(err.Error())
			}
			if conf != nil {
				if isPrivate {
					privateConfs = append(privateConfs, conf)
				} else {
					publicConfs = append(publicConfs, conf)
				}
			} else {

			}
		}
	}
	return privateConfs, publicConfs
}
