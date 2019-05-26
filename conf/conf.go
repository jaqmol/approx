package conf

// Conf ...
type Conf interface {
	Type() Type
	Name() string
	Inputs() []string
	Outputs() []string
	Assign() map[string]string
	Required() map[string]RequiredType
}

// Type ...
type Type int

// ConfTypes ...
const (
	TypeProcess Type = iota
	TypeHTTP
	TypeFork
	TypeMerge
)

// RequiredType ...
type RequiredType int

// RequiredTypes
const (
	RequiredTypeProperty RequiredType = iota
	RequiredTypeAssign
)

func addAssignmentsToRequired(assign map[string]string, required map[string]RequiredType) {
	for _, v := range assign {
		required[v] = RequiredTypeAssign
	}
}

// func filterMapValuesForVars(mapping map[string]string) (values []string) {
// 	values = make([]string, 0)
// 	for _, v := range mapping {
// 		if strings.HasPrefix(v, "$") {
// 			values = append(values, v)
// 		}
// 	}
// 	return
// }

// func filterSliceForVars(slice []string) []string {
// 	acc := make([]string, 0)
// 	for _, v := range slice {
// 		if isVariable(v) {
// 			acc = append(acc, v)
// 		}
// 	}
// 	return acc
// }
