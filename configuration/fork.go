package configuration

// Fork ...
type Fork struct {
	Ident string
	Count int
}

// Type ...
func (f *Fork) Type() ProcessorType {
	return ForkType
}

// ID ...
func (f *Fork) ID() string {
	return f.Ident
}
