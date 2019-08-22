package definition

import (
	"fmt"
	"strings"
)

// Definition ...
type Definition struct {
	Type    Type
	Name    string
	Env     map[string]*string
	Command string
}

// IsPublic ...
func (d *Definition) IsPublic() bool {
	first := d.Name[:1]
	return first == strings.ToUpper(first)
}

// EnvSlice ...
func (d *Definition) EnvSlice() []string {
	acc := make([]string, len(d.Env))
	idx := 0
	for key, value := range d.Env {
		acc[idx] = fmt.Sprintf("%v=%v", key, value)
		idx++
	}
	return acc
}

// Required ...
func (d *Definition) Required() []string {
	acc := make([]string, 0)
	for key, value := range d.Env {
		if value == nil {
			acc = append(acc, key)
		}
	}
	return acc
}

// // IsPrivate ...
// func (d *Definition) IsPrivate() bool {
// 	first := d.Name[:1]
// 	return first == strings.ToLower(first)
// }
