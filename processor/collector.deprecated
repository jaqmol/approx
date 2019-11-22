package processor

import (
	"bytes"
	"fmt"
	"io"

	"github.com/jaqmol/approx/event"
)

// Processors are best implemented pull-based
// A collector reads events from a source as fast as it can provide

// Collector ...
type Collector struct {
	serializer chan []byte
	readers    []io.Reader
}

// About closing the serializer channel:

// The Channel Closing Principle
// Don't close a channel from the receiver side and
// don't close a channel if the channel has multiple concurrent senders.

// NewCollector ...
func NewCollector(inputs ...io.Reader) (*Collector, error) {
	c := &Collector{
		serializer: make(chan []byte),
		readers:    make([]io.Reader, 0),
	}
	if len(inputs) == 0 {
		return c, nil
	}
	err := c.Connect(inputs...)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Connect ...
func (c *Collector) Connect(inputs ...io.Reader) error {
	if readerSliceIsUsable(inputs) {
		c.readers = append(c.readers, inputs...)
		return nil
	}
	return fmt.Errorf("Input(s) not usable: %v", inputs)
}

// Events ...
func (c *Collector) Events() <-chan []byte {
	return c.serializer
}

// Start ...
func (c *Collector) Start() {
	for idx := range c.readers {
		go c.startWithIndex(idx)
	}
}

func (c *Collector) startWithIndex(index int) {
	scanner := event.NewScanner(c.readers[index])
	for scanner.Scan() {
		raw := scanner.Bytes()
		original := bytes.Trim(raw, "\x00")
		toPassOn := make([]byte, len(original))
		copy(toPassOn, original)
		c.serializer <- toPassOn
	}
}
