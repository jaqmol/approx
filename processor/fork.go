package processor

import (
	"io"
	"log"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/message"
)

// Fork ...
type Fork struct {
	conf       *configuration.Fork
	in         io.Reader
	outs       []procPipe
	err        procPipe
	serializer chan []byte
}

// NewFork ...
func NewFork(conf *configuration.Fork, input io.Reader, outsCount int) *Fork {
	f := Fork{
		conf: conf,
		in:   input,
		outs: make([]procPipe, outsCount),
		err:  newProcPipe(),
	}

	for i := range f.outs {
		f.outs[i] = newProcPipe()
	}

	return &f
}

// Start ...
func (f *Fork) Start() {
	go f.readAndDistribute(f.in)
}

// Conf ...
func (f *Fork) Conf() configuration.Processor {
	return f.conf
}

// Outs ...
func (f *Fork) Outs() []io.Reader {
	acc := make([]io.Reader, len(f.outs))
	for i, p := range f.outs {
		acc[i] = p.reader
	}
	return acc
}

// Err ...
func (f *Fork) Err() io.Reader {
	return f.err.reader
}

func (f *Fork) readAndDistribute(r io.Reader) {
	msgEnd := []byte(configuration.MessageEnd)
	scanner := message.NewScanner(r)
	for scanner.Scan() {
		msg := scanner.Bytes()
		line := append(msg, msgEnd...)

		for _, p := range f.outs {
			_, err := p.writer.Write(line)
			if err != nil {
				log.Fatalln(err.Error())
			}
		}
	}
}
