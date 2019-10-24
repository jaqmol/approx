package project

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// DefinitionType ...
type DefinitionType int

// DefinitionType ...
const (
	CommandType DefinitionType = iota
	ForkType
	MergeType
)

// Definition ...
type Definition interface {
	Type() DefinitionType
	Name() string
}

// LoadDefinition ...
func LoadDefinition(projectDirectory string) ([]Definition, error) {
	definitionFilepath := filepath.Join(projectDirectory, "definition.yaml")
	_, err := os.Stat(definitionFilepath)
	if !os.IsNotExist(err) {
		return loadDefinitionFromPath(definitionFilepath)
	}
	return nil, err
}

func loadDefinitionFromPath(definitionFilepath string) ([]Definition, error) {
	data, err := ioutil.ReadFile(definitionFilepath)
	if err != nil {
		return nil, err
	}
	var parsed []map[string]interface{}
	err = yaml.Unmarshal(data, &parsed)
	if err != nil {
		return nil, err
	}
	return interpreteDefinition(parsed)
}

func interpreteDefinition(dataMap []map[string]interface{}) ([]Definition, error) {
	defs := make([]Definition, len(dataMap))
	for i, d := range dataMap {
		switch d["type"] {
		case "command":
			defs[i] = NewCommand(d)
		case "fork":
			defs[i] = NewFork(d)
		case "merge":
			defs[i] = NewMerge(d)
		}
	}
	return defs, nil
}
