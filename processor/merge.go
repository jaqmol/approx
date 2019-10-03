package processor

import (
	"io"
	"log"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/message"
)

// Merge ...
type Merge struct {
	conf       *configuration.Merge
	ins        []io.Reader
	out        procPipe
	err        procPipe
	serializer chan []byte
}

// NewMerge ...
func NewMerge(conf *configuration.Merge, inputs []io.Reader) *Merge {
	m := Merge{
		conf:       conf,
		ins:        inputs,
		out:        newProcPipe(),
		err:        newProcPipe(),
		serializer: make(chan []byte),
	}
	return &m
}

// Start ...
func (m *Merge) Start() {
	for _, r := range m.ins {
		go m.readFrom(r)
	}

	msgEnd := []byte(configuration.MessageEnd)

	for msg := range m.serializer {
		line := append(msg, msgEnd...)
		_, err := m.out.writer.Write(line)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
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
		m.serializer <- scanner.Bytes()
	}
	close(m.serializer)
}
