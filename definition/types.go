package definition

// Type ...
type Type int

// DefinitionTypes
const (
	TypeHTTPServer Type = iota
	TypeFork
	TypeMerge
	TypeProcess
)
