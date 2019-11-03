package configuration

// ProcessorType ...
type ProcessorType int

// ProcessorType ...
const (
	StdinType ProcessorType = iota
	CommandType
	ForkType
	MergeType
	StdoutType
)

// Processor ...
type Processor interface {
	Type() ProcessorType
	ID() string
}
