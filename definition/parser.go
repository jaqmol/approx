package definition

import "strings"

// Parse ...
func Parse(rawFormation map[interface{}]interface{}) []Definition {
	acc := make([]Definition, 0)
	for interfaceDifinitionHead, interfaceDifinitionBody := range rawFormation {
		difinitionHead := interfaceDifinitionHead.(string)
		if difinitionHead != "Flow" {
			difinitionBodyInterfaceMap := interfaceDifinitionBody.(map[interface{}]interface{})
			defType, defName := definitionTypeAndName(difinitionHead)
			definition := Definition{Type: defType, Name: defName}
			for interfaceKey, interfaceValue := range difinitionBodyInterfaceMap {
				key := interfaceKey.(string)
				switch key {
				case "ASSIGN":
					definition.Assign = toStringMap(interfaceValue)
				case "ENV":
					definition.Env = toStringMap(interfaceValue)
				case "COMMAND":
					definition.Command = interfaceValue.(string)
				}
			}
			// fmt.Printf("definition: %v\n", definition)
			acc = append(acc, definition)
		}
	}
	return acc
}

func toStringMap(interfaceValue interface{}) map[string]string {
	interfaceMap := interfaceValue.(map[interface{}]interface{})
	acc := make(map[string]string)
	for interfaceKey, interfaceValue := range interfaceMap {
		key := interfaceKey.(string)
		value := interfaceValue.(string)
		acc[key] = value
	}
	return acc
}

func definitionTypeAndName(definition string) (Type, string) {
	firstSpaceIdx := strings.Index(definition, " ")
	typeStr := definition[:firstSpaceIdx]
	var defType Type
	switch typeStr {
	case "HttpServer":
		defType = TypeHTTPServer
	case "Fork":
		defType = TypeFork
	case "Merge":
		defType = TypeMerge
	case "Process":
		defType = TypeProcess
	}
	return defType, strings.TrimSpace(definition[firstSpaceIdx:])
}
