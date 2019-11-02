package configuration

// Merge ...
type Merge struct {
	Ident string
	Count int
}

// Type ...
func (m *Merge) Type() ProcessorType {
	return MergeType
}

// ID ...
func (m *Merge) ID() string {
	return m.Ident
}
