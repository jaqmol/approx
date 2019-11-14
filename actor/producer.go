package actor

import (
	"io"
	"log"
)

// Producer ...
type Producer struct {
	next []Actable
}

// NewProducer ...
func NewProducer(inboxSize int) *Producer {
	return &Producer{
		next: make([]Actable, 0),
	}
}

// Produce ...
func (p *Producer) Produce(produce func() ([]byte, error)) error {
	var data []byte
	var err error
	for {
		data, err = produce()
		if err != nil {
			break
		}
		msg := NewDataMessage(data)
		for _, na := range p.next {
			na.Receive(msg)
		}
	}
	p.sendCloseMessage()
	if err == io.EOF {
		return nil
	}
	return err
}

// Next ...
func (p *Producer) Next(next ...Actable) {
	for _, na := range next {
		if na != nil {
			p.next = append(p.next, na)
		}
	}
}

// Receive ...
func (p *Producer) Receive(message []byte) {
	log.Fatalln("Producer cannot receive messages")
}

func (p *Producer) sendCloseMessage() {
	for _, na := range p.next {
		na.Receive(NewCloseMessage())
	}
}
