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
func LoadDefinition(projectDirectory string) (map[string]Definition, error) {
	definitionFilepath := filepath.Join(projectDirectory, "definition.yaml")
	_, err := os.Stat(definitionFilepath)
	if !os.IsNotExist(err) {
		return loadDefinitionFromPath(definitionFilepath)
	}
	return nil, err
}

func loadDefinitionFromPath(definitionFilepath string) (map[string]Definition, error) {
	data, err := ioutil.ReadFile(definitionFilepath)
	if err != nil {
		return nil, err
	}
	var parsed map[string]map[string]interface{}
	err = yaml.Unmarshal(data, &parsed)
	if err != nil {
		return nil, err
	}
	return interpreteDefinition(parsed)
}

func interpreteDefinition(dataMap map[string]map[string]interface{}) (map[string]Definition, error) {
	defs := make(map[string]Definition)
	for name, data := range dataMap {
		switch data["type"] {
		case "command":
			defs[name] = NewCommand(name, data)
		case "fork":
			defs[name] = NewFork(name, data)
		case "merge":
			defs[name] = NewMerge(name, data)
		}
	}
	return defs, nil
}
