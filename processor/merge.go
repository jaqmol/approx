package processor

import (
	"io"
	"log"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/message"
)

// Merge ...
type Merge struct {
	conf      *configuration.Merge
	ins       []io.Reader
	out       procPipe
	err       procPipe
	serialize chan []byte
}

// NewMerge ...
func NewMerge(conf *configuration.Merge, inputs []io.Reader) *Merge {
	m := Merge{
		conf:      conf,
		ins:       inputs,
		out:       newProcPipe(),
		err:       newProcPipe(),
		serialize: make(chan []byte),
	}
	return &m
}

// Start ...
func (m *Merge) Start() {
	for _, r := range m.ins {
		go m.readFrom(r)
	}
	go m.start()
}

// Conf ...
func (m *Merge) Conf() configuration.Processor {
	return m.conf
}

// Outs ...
func (m *Merge) Outs() []io.Reader {
	return []io.Reader{m.out.reader}
}

// Err ...
func (m *Merge) Err() io.Reader {
	return m.err.reader
}

func (m *Merge) readFrom(r io.Reader) {
	scanner := message.NewScanner(r)
	for scanner.Scan() {
		m.serialize <- msgEndedCopy(scanner.Bytes())
	}
	close(m.serialize)
}

func (m *Merge) start() {
	for msg := range m.serialize {
		n, err := m.out.writer.Write(msg)
		if err != nil {
			log.Fatalln(err.Error())
		}
		if n != len(msg) {
			log.Fatalln("Merge couldn't write complete message")
		}
	}
}
