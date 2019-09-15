package channel

// Pipe ...
type Pipe interface {
	Read() <-chan []byte
	Write() chan<- []byte
	IsTapped() bool
	SetStderr(Writer)
}

type pipeImpl struct {
	channel chan []byte
}

// NewPipe ...
func NewPipe() Pipe {
	p := pipeImpl{make(chan []byte)}
	// go p.start()
	return &p
}

// func (p *pipeImpl) start() {

// }

func (p *pipeImpl) Read() <-chan []byte {
	return p.channel
}

func (p *pipeImpl) Write() chan<- []byte {
	return p.channel
}

func (p *pipeImpl) IsTapped() bool {
	return false
}

func (p *pipeImpl) SetStderr(Writer) {
	panic("A standard pipe cannot be tapped")
}
