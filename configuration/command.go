package configuration

// Command ...
type Command struct {
	Ident string
	Cmd   string
	Env   []string
}

// Type ...
func (c *Command) Type() ProcessorType {
	return CommandType
}

// ID ...
func (c *Command) ID() string {
	return c.Ident
}
