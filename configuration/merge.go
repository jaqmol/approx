package configuration

// Merge ...
type Merge struct {
	Ident    string
	NextProc Processor
}

// Type ...
func (m *Merge) Type() ProcessorType {
	return MergeType
}

// ID ...
func (m *Merge) ID() string {
	return m.Ident
}

// Next ...
func (m *Merge) Next() []Processor {
	return []Processor{m.NextProc}
}
