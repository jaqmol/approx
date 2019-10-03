package configuration

// Merge ...
type Merge struct {
	Ident    string
	SubProcs []Processor
}

// Type ...
func (f *Merge) Type() ProcessorType {
	return MergeType
}

// ID ...
func (f *Merge) ID() string {
	return f.Ident
}

// Subs ...
func (f *Merge) Subs() []Processor {
	return f.SubProcs
}
