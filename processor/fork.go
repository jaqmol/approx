package processor

import (
	"fmt"
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
func NewFork(conf *configuration.Fork /*, input io.Reader TODO: REMOVE */) (*Fork, error) {
	// if input == nil { TODO: REMOVE
	// 	return nil, fmt.Errorf("Fork processor %v requires 1 input", conf.ID())
	// }
	if conf.Count < 2 {
		return nil, fmt.Errorf("Fork processor %v requires more than 1 output", conf.ID())
	}

	f := Fork{
		conf: conf,
		in:   nil,
		outs: make([]procPipe, conf.Count),
		err:  newProcPipe(),
	}

	for i := range f.outs {
		pp := newProcPipe()
		f.outs[i] = *pp
	}

	return &f, nil
}

// Connect ...
func (f *Fork) Connect(inputs ...io.Reader) error {
	if f.in != nil {
		return fmt.Errorf("Fork can only be connected once")
	}
	f.in = inputs[0]
	return nil
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

// Out ...
func (f *Fork) Out() io.Reader {
	return f.outs[0].reader()
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
