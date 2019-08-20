package flow

import "strings"

// Parse ...
func Parse(rawFormation map[interface{}]interface{}) map[string][]string {
	acc := make(map[string][]string)
	// TODO: 1. findRawFlow, 2. transform to map[string][]string
	return acc
}

func findRawFlow(rawFormation map[interface{}]interface{}) [][]string {
	for key, value := range rawFormation {
		if key == "Flow" {
			interfaceSlice := value.([]interface{})
			rawFlow := make([][]string, len(interfaceSlice))
			for _, interfaceLine := range interfaceSlice {
				line := interfaceLine.(string)
				names := strings.Split(line, "->")
				for i, name := range names {
					names[i] = strings.TrimSpace(name)
				}
				rawFlow = append(rawFlow, names)
			}
			return rawFlow
		}
	}
	return nil
}
