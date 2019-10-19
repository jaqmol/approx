package processor

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/jaqmol/approx/event"

	"github.com/jaqmol/approx/configuration"
)

// Command ...
type Command struct {
	conf      *configuration.Command
	waitGroup sync.WaitGroup
	scanner   *bufio.Scanner
	cmd       *exec.Cmd
	cmdIn     *procPipe
	cmdOut    *procPipe
	cmdErr    *procPipe
}

// NewCommand ...
func NewCommand(conf *configuration.Command, input io.Reader) *Command {
	cmd, args := cmdAndArgs(conf.Cmd)

	c := Command{
		conf:      conf,
		waitGroup: sync.WaitGroup{},
		scanner:   event.NewScanner(input),
		cmd:       exec.Command(cmd, args...),
		cmdIn:     newProcPipe(),
		cmdOut:    newProcPipe(),
		cmdErr:    newProcPipe(),
	}

	if c.conf.Env != nil && len(c.conf.Env) > 0 {
		c.cmd.Env = append(os.Environ(), c.conf.Env...)
	}

	c.cmd.Stdin = c.cmdIn.reader()
	c.cmd.Stdout = c.cmdOut.writer()
	c.cmd.Stderr = c.cmdErr.writer()

	/* Last working state:
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr
	*/

	return &c
}

// Start ...
func (c *Command) Start() {
	go c.startReadingInput()
	go c.startCmd()
}

// Conf ...
func (c *Command) Conf() configuration.Processor {
	return c.conf
}

// Outs ...
func (c *Command) Outs() []io.Reader {
	return []io.Reader{c.cmdOut.reader()}
}

// Err ...
func (c *Command) Err() io.Reader {
	return c.cmdErr.reader()
}

// Wait ...
func (c *Command) Wait() {
	c.waitGroup.Wait()
}

func (c *Command) startCmd() {
	c.waitGroup.Add(1)
	var err error
	err = c.cmd.Start()
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = c.cmd.Wait()
	if err != nil {
		log.Fatalln(err.Error())
	}
	c.waitGroup.Done()
}

func (c *Command) startReadingInput() {
	c.waitGroup.Add(1)
	for c.scanner.Scan() {
		raw := bytes.Trim(c.scanner.Bytes(), "\x00")
		data := evntEndedCopy(raw)
		// log.Println(string(data))
		n, err := c.cmdIn.writer().Write(data)
		if err != nil {
			log.Fatalln(err.Error())
		}
		if n != len(data) {
			log.Fatalln("Command couldn't write complete event")
		}
	}
	c.waitGroup.Done()
}

func cmdAndArgs(cmdPlusArgs string) (string, []string) {
	acc := make([]string, 0)
	comps := strings.Split(cmdPlusArgs, " ")
	for _, cmp := range comps {
		if len(cmp) > 0 {
			acc = append(acc, cmp)
		}
	}
	return acc[0], acc[1:]
}