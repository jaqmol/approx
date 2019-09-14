package channel

type tappedPipe struct {
	channel chan []byte
	stderr  Writer
}

// NewTappedPipe ...
func NewTappedPipe(name string) Pipe {
	p := tappedPipe{
		channel: make(chan []byte),
	}
	// go p.start()
	return &p
}

// func (p *pipeImpl) start() {

// }

func (p *tappedPipe) Read() <-chan []byte {
	return p.channel
}

func (p *tappedPipe) Write() chan<- []byte {
	return p.channel
}

func (p *tappedPipe) IsTapped() bool {
	return true
}

func (p *tappedPipe) SetStderr(w Writer) {
	p.stderr = w
}
