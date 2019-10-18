package processor

import (
	"bytes"
	"io"
	"log"
	"strings"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/event"
)

// Merge ...
type Merge struct {
	conf               *configuration.Merge
	ins                []io.Reader
	out                *procPipe
	err                *procPipe
	serialize          chan []byte
	changeScannerCount chan int
	scannerCount       int
}

// NewMerge ...
func NewMerge(conf *configuration.Merge, inputs []io.Reader) *Merge {
	m := Merge{
		conf:               conf,
		ins:                inputs,
		out:                newProcPipe(),
		err:                newProcPipe(),
		serialize:          make(chan []byte),
		changeScannerCount: make(chan int),
	}
	return &m
}

// Start ...
func (m *Merge) Start() {
	for _, r := range m.ins {
		go m.readAndSynchronize(r)
	}
	go m.start()
}

// Conf ...
func (m *Merge) Conf() configuration.Processor {
	return m.conf
}

// Outs ...
func (m *Merge) Outs() []io.Reader {
	return []io.Reader{m.out.reader()}
}

// Err ...
func (m *Merge) Err() io.Reader {
	return m.err.reader()
}

func (m *Merge) readAndSynchronize(r io.Reader) {
	m.changeScannerCount <- 1
	scanner := event.NewScanner(r)
	for scanner.Scan() {
		msg := evntEndedCopy(scanner.Bytes())
		m.serialize <- msg
	}
	m.changeScannerCount <- -1
}

func (m *Merge) start() {
	loop := true
	for loop {
		select {
		case raw := <-m.serialize:
			msg := bytes.Trim(raw, "\x00")
			n, err := m.out.writer().Write(msg)
			if err != nil {
				log.Fatalln(err.Error())
			}
			if n != len(msg) {
				log.Fatalln("Merge couldn't write complete event")
			}
		case amount := <-m.changeScannerCount:
			m.scannerCount += amount
			loop = m.scannerCount > 0
		}
	}
	m.stop()
}

func (m *Merge) stop() {
	errs := m.out.close()
	if len(errs) > 0 {
		s := strings.Join(errsToStrs(errs), ", ")
		log.Fatalf("Errors closing pipe: %s\n", s)
	}
}
