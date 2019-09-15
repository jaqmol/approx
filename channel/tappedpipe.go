package channel

import (
	"github.com/jaqmol/approx/utils"
)

type tappedPipe struct {
	outChan chan []byte
	inChan  chan []byte
	stderr  Writer
}

// NewTappedPipe ...
func NewTappedPipe(name string) Pipe {
	p := tappedPipe{
		outChan: make(chan []byte),
		inChan:  make(chan []byte),
	}
	go p.start()
	return &p
}

func (p *tappedPipe) start() {
	for b := range p.inChan {
		msg := append([]byte{}, b...)
		// msg = append(msg, '\n')
		p.stderr.Write() <- utils.Truncated(msg, 100)
		p.outChan <- b
	}
}

func (p *tappedPipe) Read() <-chan []byte {
	return p.outChan
}

func (p *tappedPipe) Write() chan<- []byte {
	return p.inChan
}

func (p *tappedPipe) IsTapped() bool {
	return true
}

func (p *tappedPipe) SetStderr(w Writer) {
	p.stderr = w
}
