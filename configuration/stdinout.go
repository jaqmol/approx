package configuration

// Stdinout ...
type Stdinout struct {
	defType ProcessorType
	ident   string
}

// Stdin ...
var Stdin Stdinout

// Stdout ...
var Stdout Stdinout

func init() {
	Stdin = Stdinout{StdinType, "<stdin>"}
	Stdout = Stdinout{StdoutType, "<stdout>"}
}

// Type ...
func (sio *Stdinout) Type() ProcessorType {
	return sio.defType
}

// ID ...
func (sio *Stdinout) ID() string {
	return sio.ident
}
