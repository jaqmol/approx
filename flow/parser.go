package flow

import (
	"strings"
)

// Parse ...
func Parse(rawFormation map[interface{}]interface{}) map[string][]string {
	acc := make(map[string][]string)
	rawFlow := findRawFlow(rawFormation)
	for _, line := range rawFlow {
		var from string
		for i, name := range line {
			if i == 0 {
				from = name
			} else {
				tos, ok := acc[from]
				if !ok {
					tos = make([]string, 0)
				}
				tos = append(tos, name)
				acc[from] = tos
				from = name
			}
		}
	}
	return acc
}

func findRawFlow(rawFormation map[interface{}]interface{}) [][]string {
	for key, value := range rawFormation {
		if key == "Flow" {
			interfaceSlice := value.([]interface{})
			rawFlow := make([][]string, 0) // ?
			for _, interfaceLine := range interfaceSlice {
				line := interfaceLine.(string)
				names := strings.Split(line, "->")
				for i, name := range names {
					names[i] = strings.TrimSpace(name) // ?
				}
				rawFlow = append(rawFlow, names)
			}
			return rawFlow
		}
	}
	return nil
}
