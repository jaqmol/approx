package configuration

// Command ...
type Command struct {
	Ident    string
	Path     string
	Args     []string
	Env      []string
	Dir      string
	NextProc Processor
}

// Type ...
func (c *Command) Type() ProcessorType {
	return CommandType
}

// ID ...
func (c *Command) ID() string {
	return c.Ident
}

// Next ...
func (c *Command) Next() []Processor {
	return []Processor{c.NextProc}
}
