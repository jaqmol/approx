package configuration

// Command ...
type Command struct {
	Ident string
	Cmd   string
	Env   []string
	// NextProc Processor // TODO: REMOVE
}

// Type ...
func (c *Command) Type() ProcessorType {
	return CommandType
}

// ID ...
func (c *Command) ID() string {
	return c.Ident
}

// TODO: REMOVE
// // Next ...
// func (c *Command) Next() []Processor {
// 	return []Processor{c.NextProc}
// }

// TODO: REMOVE
// // SetNext ...
// func (c *Command) SetNext(next ...Processor) {
// 	c.NextProc = next[0]
// }
