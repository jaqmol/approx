package assign

import (
	"log"
	"strings"

	"github.com/jaqmol/approx/definition"
)

// ResolveVariables ...
func ResolveVariables(rawFormation map[interface{}]interface{}, definitions []definition.Definition) {
	defForName := defForNameMap(definitions)
	assigns := parse(rawFormation)
	vars := variables(assigns, defForName)
	for _, d := range definitions {
		for envName, envValuePtr := range d.Env {
			envValue := *envValuePtr
			if strings.HasPrefix(envValue, "$") {
				realValue, ok := vars[envValue]
				if !ok {
					log.Fatalf("Variable %v in processor %v could not be resolved", envValue, d.Name)
				}
				d.Env[envName] = &realValue
			}
		}
	}
}

func defForNameMap(ds []definition.Definition) map[string]definition.Definition {
	acc := make(map[string]definition.Definition)
	for _, d := range ds {
		acc[d.Name] = d
	}
	return acc
}
