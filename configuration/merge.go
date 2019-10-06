package configuration

// Merge ...
type Merge struct {
	Ident    string
	NextProc Processor
}

// Type ...
func (f *Merge) Type() ProcessorType {
	return MergeType
}

// ID ...
func (f *Merge) ID() string {
	return f.Ident
}

// Next ...
func (f *Merge) Next() []Processor {
	return []Processor{f.NextProc}
}
