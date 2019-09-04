package assign

import "strings"

func parse(rawFormation map[interface{}]interface{}) map[string]string {
	for key, value := range rawFormation {
		KEY := strings.ToUpper(key.(string))
		if KEY == "ASSIGN" {
			interfaceMap := value.(map[interface{}]interface{})
			assign := make(map[string]string) // ?
			for interfaceName, interfaceValue := range interfaceMap {
				name := interfaceName.(string)
				value := interfaceValue.(string)
				assign[name] = value
			}
			return assign
		}
	}
	return nil
}
