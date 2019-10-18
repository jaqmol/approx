package processor

import (
	"io"
	"log"
	"strings"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/event"
)

// Fork ...
type Fork struct {
	conf *configuration.Fork
	in   io.Reader
	outs []procPipe
	err  *procPipe
}

// NewFork ...
func NewFork(conf *configuration.Fork, input io.Reader) *Fork {
	f := Fork{
		conf: conf,
		in:   input,
		outs: make([]procPipe, len(conf.NextProcs)),
		err:  newProcPipe(),
	}

	for i := range f.outs {
		pp := newProcPipe()
		f.outs[i] = *pp
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
		acc[i] = p.reader()
	}
	return acc
}

// Err ...
func (f *Fork) Err() io.Reader {
	return f.err.reader()
}

func (f *Fork) readAndDistribute(r io.Reader) {
	scanner := event.NewScanner(r)
	for scanner.Scan() {
		msg := evntEndedCopy(scanner.Bytes())
		for _, p := range f.outs {
			n, err := p.writer().Write(msg)
			if err != nil {
				log.Fatalln(err.Error())
			}
			if n != len(msg) {
				log.Fatalln("Fork couldn't write complete event")
			}
		}
	}
	f.stop()
}

func (f *Fork) stop() {
	errs := closeProcPipes(f.outs)
	if len(errs) > 0 {
		s := strings.Join(errsToStrs(errs), ", ")
		log.Fatalf("Errors closing pipe: %s\n", s)
	}
}
