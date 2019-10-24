package project

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// "gopkg.in/yaml.v2"

// Formation ...
type Formation struct {
	Definitions []Definition // `yaml:"definition,omitempty"`
	Flow        []Flow       // `yaml:"flow,omitempty"`
}

// LoadFormation ...
func LoadFormation(projectDirectory string) (*Formation, error) {
	formationFilepath := filepath.Join(projectDirectory, "formation.yaml")
	if _, err := os.Stat(formationFilepath); !os.IsNotExist(err) {
		return loadFormationFromPath(formationFilepath)
	}
	errsAcc := make([]error, 0)
	def, err := LoadDefinition(projectDirectory)
	if err != nil {
		errsAcc = append(errsAcc, err)
	}
	flow, err := LoadFlow(projectDirectory)
	if err != nil {
		errsAcc = append(errsAcc, err)
	}
	if len(errsAcc) > 0 {
		return nil, noFormationError(errsAcc)
	}
	f := Formation{
		Definitions: def,
		Flow:        flow,
	}
	return &f, nil
}

func loadFormationFromPath(formationFilepath string) (*Formation, error) {
	data, err := ioutil.ReadFile(formationFilepath)
	if err != nil {
		return nil, err
	}
	var forMap map[string]interface{}
	err = yaml.Unmarshal(data, &forMap)
	if err != nil {
		return nil, err
	}
	defData := forMap["definition"].([]map[string]interface{})
	defs, err := interpreteDefinition(defData)
	if err != nil {
		return nil, err
	}
	flowData := forMap["flow"].([][]string)
	flows, err := interpreteFlow(flowData)
	if err != nil {
		return nil, err
	}
	return &Formation{defs, flows}, nil
}

func noFormationError(causes []error) error {
	buff := strings.Builder{}
	for i, e := range causes {
		if i > 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(e.Error())
	}
	return fmt.Errorf("No formation.yaml found, expected definition.yaml and flow.yaml, but found error(s): %v", buff.String())
}
