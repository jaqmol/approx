package configuration

// Fork ...
type Fork struct {
	Ident     string
	NextProcs []Processor
}

// Type ...
func (f *Fork) Type() ProcessorType {
	return ForkType
}

// ID ...
func (f *Fork) ID() string {
	return f.Ident
}

// Next ...
func (f *Fork) Next() []Processor {
	return f.NextProcs
}
