package definition

import "strings"

// Definition ...
type Definition struct {
	Type    Type
	Name    string
	Assign  map[string]string
	Env     map[string]string
	Command string
}

// IsPublic ...
func (d *Definition) IsPublic() bool {
	first := d.Name[:1]
	return first == strings.ToUpper(first)
}

// IsPrivate ...
func (d *Definition) IsPrivate() bool {
	first := d.Name[:1]
	return first == strings.ToLower(first)
}
