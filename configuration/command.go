package configuration

// Command ...
type Command struct {
	Ident    string
	Cmd      string
	Env      []string
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
