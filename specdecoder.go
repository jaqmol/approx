package main

// NewSpecDecoder ...
func NewSpecDecoder(data map[string]interface{}) *SpecDecoder {
	return &SpecDecoder{data: data}
}

// SpecDecoder ...
type SpecDecoder struct {
	data map[string]interface{}
}

// String ...
func (c *SpecDecoder) String(name string) (retVal *string) {
	if undef, ok := c.data[name]; ok {
		if value, ok := undef.(string); ok {
			retVal = &value
		}
	}
	return
}

// Int ...
func (c *SpecDecoder) Int(name string) (retVal *int) {
	if undef, ok := c.data[name]; ok {
		if value, ok := undef.(int); ok {
			retVal = &value
		}
	}
	return
}

// StringSlice ...
func (c *SpecDecoder) StringSlice(name string) (retVal []string) {
	if undef, ok := c.data[name]; ok {
		if value, ok := undef.([]string); ok {
			retVal = value
		}
	}
	return
}

// StringStringMap ...
func (c *SpecDecoder) StringStringMap(name string) (retVal map[string]string) {
	if undef, ok := c.data[name]; ok {
		if value, ok := undef.(map[string]string); ok {
			retVal = value
		}
	}
	return
}
