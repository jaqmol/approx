package assign

import (
	"strings"

	"github.com/jaqmol/approx/definition"
)

var variables = make(map[string]string)

// // Variable ...
// func Variable(name string, value string) {
// 	variables[name] = value
// }

// GetVariable ...
func GetVariable(name string) string {
	return variables[name]
}

// Variables ...
func Variables(defs []definition.Definition) {
	for _, def := range defs {
		for name, value := range def.Assign {
			if strings.HasPrefix(name, "$") {
				variables[name] = value
			}
		}
	}
}
