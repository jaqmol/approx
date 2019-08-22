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
					definition.Env = toStringPointerMap(interfaceValue)
				case "COMMAND":
					definition.Command = interfaceValue.(string)
				default:
					if strings.HasPrefix(key, "ENV") {
						definition.Env = parseKeyPathEnvAssignment(key, interfaceValue)
					}
				}
			}
			acc = append(acc, definition)
		}
	}
	return acc
}

func parseKeyPathEnvAssignment(keyPathEnv string, interfaceValue interface{}) map[string]*string {
	env := make(map[string]*string)
	keyIdx := strings.Index(keyPathEnv, ".") + 1
	key := keyPathEnv[keyIdx:]
	value := interfaceValue.(string)
	env[key] = &value
	return env
}

func toStringSlice(interfaceValue interface{}) []string {
	interfaceSlice := interfaceValue.([]interface{})
	acc := make([]string, 0)
	for _, interfaceValue := range interfaceSlice {
		value := interfaceValue.(string)
		acc = append(acc, value)
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

func toStringPointerMap(interfaceValue interface{}) map[string]*string {
	interfaceMap := interfaceValue.(map[interface{}]interface{})
	acc := make(map[string]*string)
	for interfaceKey, interfaceValue := range interfaceMap {
		key := interfaceKey.(string)
		if interfaceValue == nil {
			acc[key] = nil
		} else {
			value := interfaceValue.(*string)
			acc[key] = value
		}
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
