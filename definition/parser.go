package definition

import (
	"fmt"
	"strings"
)

// Parse ...
func Parse(rawFormation map[interface{}]interface{}) []Definition {
	acc := make([]Definition, 0)
	for interfaceDifinitionHead, interfaceDifinitionBody := range rawFormation {
		definitionHead := interfaceDifinitionHead.(string)
		DIFINITIONHEAD := strings.ToUpper(definitionHead)
		if DIFINITIONHEAD != "FLOW" && DIFINITIONHEAD != "ASSIGN" {
			difinitionBodyInterfaceMap := interfaceDifinitionBody.(map[interface{}]interface{})
			defType, defName := definitionTypeAndName(definitionHead)
			definition := Definition{Type: defType, Name: defName}
			for interfaceKey, interfaceValue := range difinitionBodyInterfaceMap {
				key := interfaceKey.(string)
				switch key {
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
	if interfaceValue == nil {
		env[key] = nil
	} else {
		env[key] = stringPtrFromStringOrIntValue(interfaceValue)
	}
	return env
}

func toStringPointerMap(interfaceValue interface{}) map[string]*string {
	interfaceMap := interfaceValue.(map[interface{}]interface{})
	acc := make(map[string]*string)
	for interfaceKey, interfaceValue := range interfaceMap {
		key := interfaceKey.(string)
		if interfaceValue == nil {
			acc[key] = nil
		} else {
			acc[key] = stringPtrFromStringOrIntValue(interfaceValue)
		}
	}
	return acc
}

func stringPtrFromStringOrIntValue(interfaceValue interface{}) *string {
	strValue, ok := interfaceValue.(string)
	if !ok {
		intValue := interfaceValue.(int)
		strValue = fmt.Sprintf("%v", intValue)
	}
	return &strValue
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

func definitionTypeAndName(definition string) (Type, string) {
	firstSpaceIdx := strings.Index(definition, " ")
	typeStr := strings.ToUpper(definition[:firstSpaceIdx])
	var defType Type
	switch typeStr {
	case "GATEWAY":
		fallthrough
	case "GW":
		defType = TypeHTTPServer
	case "FORK":
		defType = TypeFork
	case "MERGE":
		defType = TypeMerge
	case "PROCESS":
		fallthrough
	case "PROC":
		defType = TypeProcess
	}
	return defType, strings.TrimSpace(definition[firstSpaceIdx:])
}
