package actor

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/jaqmol/approx/config"
	"github.com/jaqmol/approx/event"
)

// Command ...
type Command struct {
	Actor
	ident string
	cmd   *exec.Cmd
	// input   io.WriteCloser
	// logging io.ReadCloser
	// output  io.ReadCloser
	input   *io.PipeWriter
	logging *io.PipeReader
	output  *io.PipeReader
}

// NewCommand ...
func NewCommand(inboxSize int, ident string, cmd string, args ...string) *Command {
	c := &Command{
		ident: ident,
		cmd:   exec.Command(cmd, args...),
	}
	// c.cmd.Stderr = os.Stderr

	inputReader, inputWriter := io.Pipe()
	// loggingReader, loggingWriter := io.Pipe()
	outputReader, outputWriter := io.Pipe()

	c.cmd.Stdin = inputReader
	// c.cmd.Stderr = loggingWriter
	c.cmd.Stdout = outputWriter

	c.input = inputWriter
	// c.logging = loggingReader
	c.output = outputReader

	// logging, err := c.cmd.StderrPipe()
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// output, err := c.cmd.StdoutPipe()
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// c.input = input
	// c.logging = logging
	// c.output = output

	c.init(inboxSize)
	return c
}

// NewCommandFromConf ...
func NewCommandFromConf(inboxSize int, conf *config.Command) (*Command, error) {
	cmdAndArgs := strings.Split(conf.Cmd, " ")
	// TODO: actor.Command ENV support missing
	var c *Command
	if len(cmdAndArgs) == 1 {
		c = NewCommand(inboxSize, conf.Ident, cmdAndArgs[0])
	} else if len(cmdAndArgs) > 1 {
		c = NewCommand(inboxSize, conf.Ident, cmdAndArgs[0], cmdAndArgs[1:]...)
	} else {
		return nil, fmt.Errorf("Command definition of \"%v\" is wrong: \"%v\"", conf.Ident, conf.Cmd)
	}
	return c, nil
}

// Logging ...
func (c *Command) Logging() io.Reader {
	return c.logging
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
				log.Fatalf(
					"Error dispatching event to command \"%v\": %v -> %v\n",
					c.ident,
					err,
					string(termMsg),
				)
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
		}
	}

	err := c.input.Close() // This triggers graceful termination
	if err != nil {
		log.Fatalf(
			"Error closing <stdin> on command \"%v\": %v\n",
			c.ident,
			err,
		)
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
