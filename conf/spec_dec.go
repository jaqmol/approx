package conf

func newSpecDec(data map[string]interface{}) *specDec {
	return &specDec{data: data}
}

type specDec struct {
	data map[string]interface{}
}

func (c *specDec) string(name string) (value string, ok bool) {
	var undef interface{}
	if undef, ok = c.data[name]; ok {
		value, ok = undef.(string)
	}
	return
}

func (c *specDec) integer(name string) (value int, ok bool) {
	var undef interface{}
	if undef, ok = c.data[name]; ok {
		value, ok = undef.(int)
	}
	return
}

func (c *specDec) stringSlice(name string) (values []string, ok bool) {
	var undef interface{}
	if undef, ok = c.data[name]; ok {
		var undefSlice []interface{}
		if undefSlice, ok = undef.([]interface{}); ok {
			values = make([]string, 0)
			for _, undefValue := range undefSlice {
				var value string
				if value, ok = undefValue.(string); ok {
					values = append(values, value)
				}
			}
		}
	}
	return
}

func (c *specDec) stringStringMap(name string) (valuesMap map[string]string, ok bool) {
	var undef interface{}
	if undef, ok = c.data[name]; ok {
		var undefMap map[string]interface{}
		if undefMap, ok = undef.(map[string]interface{}); ok {
			valuesMap = make(map[string]string)
			for key, undefValue := range undefMap {
				var value string
				if value, ok = undefValue.(string); ok {
					valuesMap[key] = value
				}
			}
		}
	}
	return
}
