package actor

import "log"

// Collector ...
type Collector struct {
	inbox chan Message
}

// NewCollector ...
func NewCollector(inboxSize int) *Collector {
	return &Collector{
		inbox: make(chan Message, inboxSize),
	}
}

// Next ...
func (c *Collector) Next(next ...Actable) {
	for _, na := range next {
		if na != nil {
			log.Fatalln("Collector cannot be chained to a next actor")
		}
	}
}

// Receive ...
func (c *Collector) Receive(message Message) {
	c.inbox <- message
}

// Collect ...
func (c *Collector) Collect(collect func([]byte) error) error {
	for msg := range c.inbox {
		var err error
		switch msg.Type {
		case DataMessage:
			err = collect(msg.Data)
		case CloseInbox:
			close(c.inbox)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Start ...
func (c *Collector) Start() {}
