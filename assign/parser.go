package assign

// Parse ...
func Parse(rawFormation map[interface{}]interface{}) map[string]string {
	for key, value := range rawFormation {
		if key == "Assign" {
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
