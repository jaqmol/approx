package project

// Merge ...
type Merge struct {
	ident string
}

// NewMerge ...
func NewMerge(ident string, originalData interface{}) *Merge {
	// data := originalData.(map[string]interface{}) // TODO: REMOVE
	f := Merge{
		ident: ident,
	}
	return &f
}

// Type ...
func (f *Merge) Type() DefinitionType {
	return MergeType
}

// Ident ...
func (f *Merge) Ident() string {
	return f.ident
}
