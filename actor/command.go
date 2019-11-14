package actor

import (
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/jaqmol/approx/event"
)

// Command ...
type Command struct {
	Actor
	ident        string
	cmd          *exec.Cmd
	logging      io.Writer
	inputWriter  io.WriteCloser
	outputReader io.ReadCloser
}

// NewCommand ...
func NewCommand(inboxSize int, ident string, cmd string, args ...string) *Command {
	c := &Command{
		ident:   ident,
		cmd:     exec.Command(cmd, args...),
		logging: os.Stderr,
	}
	c.init(inboxSize)
	return c
}

// Logging ...
func (c *Command) Logging(writer io.Writer) {
	c.logging = writer
}

// Start ...
func (c *Command) Start() {
	if len(c.next) != 1 {
		log.Fatalf(
			"Command \"%v\" is connected to %v next, 1 expected\n",
			c.ident,
			len(c.next),
		)
	}
	c.cmd.Stderr = c.logging
	c.initInputOutput()
	go c.startReceiving()
	go c.startSending()
	go c.startCommand()
}

func (c *Command) initInputOutput() {
	inputWriter, err := c.cmd.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}
	c.inputWriter = inputWriter
	outputReader, err := c.cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}
	c.outputReader = outputReader
}

func (c *Command) startReceiving() {
	for message := range c.inbox {
		switch message.Type {
		case DataMessage:
			termMsg := event.TerminatedBytesCopy(message.Data)
			n, err := c.inputWriter.Write(termMsg)
			if err != nil {
				log.Fatalln(err)
			}
			if n != len(message.Data) {
				log.Fatalf(
					"Command \"%v\" couldn't write complete event",
					c.ident,
				)
			}
		case CloseInbox:
			c.next[0].Receive(message)
			close(c.inbox)
			c.inputWriter.Close()
			c.outputReader.Close()
		}
	}
}

func (c *Command) startSending() {
	scanner := event.NewScanner(c.outputReader)
	for scanner.Scan() {
		raw := scanner.Bytes()
		data := event.ScannedBytesCopy(raw)
		msg := NewDataMessage(data)
		c.next[0].Receive(msg)
	}
}

func (c *Command) startCommand() {
	err := c.cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
