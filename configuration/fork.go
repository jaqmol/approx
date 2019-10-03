package configuration

// Fork ...
type Fork struct {
	Ident    string
	SubProcs []Processor
}

// Type ...
func (f *Fork) Type() ProcessorType {
	return ForkType
}

// ID ...
func (f *Fork) ID() string {
	return f.Ident
}

// Subs ...
func (f *Fork) Subs() []Processor {
	return f.SubProcs
}
