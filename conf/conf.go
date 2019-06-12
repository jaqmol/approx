package conf

// Conf ...
type Conf interface {
	Type() Type
	Name() string
	Inputs() []string
	Outputs() []string
	Environment() []string
	Assign() map[string]string
	Required() map[string]RequiredType
}

// Type ...
type Type int

// ConfTypes ...
const (
	TypeProcess Type = iota
	TypeHTTPServer
	TypeFork
	TypeMerge
	TypeCheck
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
