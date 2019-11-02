package configuration

// Merge ...
type Merge struct {
	Ident string
	Count int
	// NextProc Processor // TODO: REMOVE
}

// Type ...
func (m *Merge) Type() ProcessorType {
	return MergeType
}

// ID ...
func (m *Merge) ID() string {
	return m.Ident
}

// TODO: REMOVE
// // Next ...
// func (m *Merge) Next() []Processor {
// 	return []Processor{m.NextProc}
// }

// TODO: REMOVE
// // SetNext ...
// func (m *Merge) SetNext(next ...Processor) {
// 	m.NextProc = next[0]
// }
