package pipe

// Pipe ...
type Pipe struct {
	Reader *Reader
	Writer *Writer
}

// NewPipe ...
func NewPipe() *Pipe {
	p := Pipe{
		Reader: NewReader(),
		Writer: NewWriter(),
	}
	go p.start()
	return &p
}

func (p *Pipe) start() {
	for b := range p.Writer.inputChannel {
		select {
		case p.Reader.outputChannel <- b:
		default:
		}
	}
}
