package project

// Merge ...
type Merge struct {
	name string
}

// NewMerge ...
func NewMerge(name string, originalData interface{}) *Merge {
	// data := originalData.(map[string]interface{})
	f := Merge{
		name: name,
	}
	return &f
}

// Type ...
func (f *Merge) Type() DefinitionType {
	return MergeType
}

// Name ...
func (f *Merge) Name() string {
	return f.name
}
