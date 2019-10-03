package configuration

// Command ...
type Command struct {
	Ident    string
	Path     string
	Args     []string
	Env      []string
	Dir      string
	SubProcs []Processor
}

// Type ...
func (c *Command) Type() ProcessorType {
	return CommandType
}

// ID ...
func (c *Command) ID() string {
	return c.Ident
}

// Subs ...
func (c *Command) Subs() []Processor {
	return c.SubProcs
}
