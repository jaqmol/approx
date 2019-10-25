package project

// Fork ...
type Fork struct {
	name string
}

// NewFork ...
func NewFork(name string, originalData interface{}) *Fork {
	// data := originalData.(map[string]interface{})
	f := Fork{
		name: name,
	}
	return &f
}

// Type ...
func (f *Fork) Type() DefinitionType {
	return ForkType
}

// Name ...
func (f *Fork) Name() string {
	return f.name
}
