package run

import (
	"fmt"

	"github.com/jaqmol/approx/definition"
	"github.com/jaqmol/approx/pipe"
)

// MakePipes ...
func MakePipes(definitions []definition.Definition, flows map[string][]string) map[string]pipe.Pipe {
	acc := make(map[string]pipe.Pipe)

	for _, fromDef := range definitions {
		fromName := fromDef.Name
		toNames := flows[fromName]

		for _, toName := range toNames {
			key := PipeKey(fromName, toName)
			// reader, writer := io.Pipe()
			acc[key] = *pipe.NewPipe() // Pipe{ Reader: reader, Writer: writer }
		}
	}

	fmt.Printf("Did make pipes: %v\n", acc)

	return acc
}

// MakeStderrs ...
func MakeStderrs(definitions []definition.Definition) map[string]pipe.Pipe {
	acc := make(map[string]pipe.Pipe)

	for _, def := range definitions {
		// reader, writer := io.Pipe()
		acc[def.Name] = *pipe.NewPipe() // Pipe{ Reader: reader, Writer: writer }
	}

	return acc
}

// Pipe ...
// type Pipe struct {
// 	running       bool
// 	readSrcRest   []byte
// 	inputChannel  chan []byte
// 	outputChannel chan []byte
// }

// NewPipe ...
// func NewPipe() Pipe {
// 	p := Pipe{}
// 	p.inputChannel = make(chan []byte, 100)
// 	p.outputChannel = make(chan []byte, 100)
// 	go p.start()
// 	return p
// }

// // Start ...
// func (p *Pipe) Start() {
// 	if !p.running {
// 		go p.start()
// 		p.running = true
// 	}
// }

// func (p *Pipe) start() {
// 	for b := range p.inputChannel {
// 		select {
// 		case p.outputChannel <- b:
// 			// fmt.Printf("Pipe sent: %v\n", string(b))
// 		default:
// 		}
// 	}
// }

// Reader ...
// func (p *Pipe) Reader() io.Reader {
// 	return p
// }

// Writer ...
// func (p *Pipe) Writer() io.Writer {
// 	return p
// }

// PipeKey ...
func PipeKey(fromName string, toName string) string {
	return fmt.Sprintf("%v->%v", fromName, toName)
}

// type reader struct {
// }

// func (p *Pipe) Read(dst []byte) (n int, err error) {
// 	var src []byte
// 	if len(p.readSrcRest) > 0 {
// 		src = p.readSrcRest
// 		p.readSrcRest = []byte{}
// 	} else {
// 		src = <-p.outputChannel
// 	}
// 	copiedCount := copy(dst, src)
// 	if copiedCount < len(src) {
// 		p.readSrcRest = src[copiedCount:]
// 	}
// 	return copiedCount, nil
// }

// type writer struct {
// }

// func (p *Pipe) Write(b []byte) (n int, err error) {
// 	p.inputChannel <- b
// 	return len(b), nil
// }
