package configuration

// ProcessorType ...
type ProcessorType int

// ProcessorType ...
const (
	CommandType ProcessorType = iota
	ForkType
	MergeType
)

// Processor ...
type Processor interface {
	Type() ProcessorType
	ID() string
	Subs() []Processor
}
