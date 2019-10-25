package project

// Command ...
type Command struct {
	name string            // `yaml:"name,omitempty"`
	cmd  string            // `yaml:"cmd,omitempty"`
	env  map[string]string // `yaml:"env,omitempty"`
}

// NewCommand ...
func NewCommand(name string, originalData interface{}) *Command {
	data := originalData.(map[string]interface{})
	c := Command{
		name: name,
		cmd:  data["cmd"].(string),
		env:  toStringMapString(data["env"]),
	}
	return &c
}

// Type ...
func (c *Command) Type() DefinitionType {
	return CommandType
}

// Name ...
func (c *Command) Name() string {
	return c.name
}

// Cmd ...
func (c *Command) Cmd() string {
	return c.cmd
}

// Env ...
func (c *Command) Env() map[string]string {
	return c.env
}

func toStringMapString(originalData interface{}) map[string]string {
	acc := make(map[string]string)
	for keyI, valueI := range originalData.(map[interface{}]interface{}) {
		acc[keyI.(string)] = valueI.(string)
	}
	return acc
}
