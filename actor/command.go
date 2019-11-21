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
	ident  string
	cmd    *exec.Cmd
	input  io.WriteCloser
	output io.ReadCloser
}

// NewCommand ...
func NewCommand(inboxSize int, ident string, cmd string, args ...string) *Command {
	c := &Command{
		ident: ident,
		cmd:   exec.Command(cmd, args...),
	}
	c.cmd.Stderr = os.Stderr

	input, err := c.cmd.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}
	output, err := c.cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}
	c.input = input
	c.output = output

	c.init(inboxSize)
	return c
}

// Logging ...
func (c *Command) Logging(writer io.Writer) {
	c.cmd.Stderr = writer
}

// Directory ...
func (c *Command) Directory(dir string) {
	c.cmd.Dir = dir
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

	go c.startDispatchingInboxToCmd()
	go c.startReceivingCmdOutput()
	go c.startCommand()
}

func (c *Command) startDispatchingInboxToCmd() {
	for message := range c.inbox {
		switch message.Type {
		case DataMessage:
			termMsg := event.TerminatedBytesCopy(message.Data)
			n, err := c.input.Write(termMsg)
			if err != nil {
				log.Fatalln(err)
			}
			if n < len(message.Data) {
				log.Fatalf(
					"Command \"%v\" couldn't write all data, only %v/%v\n",
					c.ident,
					n,
					len(message.Data),
				)
			}
		case CloseInbox:
			close(c.inbox)
			err := c.input.Close() // This triggers graceful termination
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func (c *Command) startReceivingCmdOutput() {
	scanner := event.NewScanner(c.output)
	for scanner.Scan() {
		raw := scanner.Bytes()
		data := event.ScannedBytesCopy(raw)
		msg := NewDataMessage(data)
		c.next[0].Receive(msg)
	}
}

func (c *Command) startCommand() {
	err := c.cmd.Start()
	if err != nil {
		log.Fatalf("Error starting command \"%v\": %v\n", c.ident, err.Error())
	}
	err = c.cmd.Wait()
	if err != nil {
		log.Fatalf("Error completing command \"%v\": %v\n", c.ident, err.Error())
	}
	c.next[0].Receive(NewCloseMessage())
}
