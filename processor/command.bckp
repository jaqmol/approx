package processor

import (
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/jaqmol/approx/configuration"
)

// Command ...
type Command struct {
	conf *configuration.Command
	cmd  *exec.Cmd
	out  procPipe
	err  procPipe
}

// NewCommand ...
func NewCommand(conf *configuration.Command, input io.Reader) *Command {
	c := Command{
		conf: conf,
		out:  newProcPipe(),
		err:  newProcPipe(),
	}

	c.cmd = &exec.Cmd{
		Path:   conf.Path,
		Args:   conf.Args,
		Env:    conf.Env,
		Dir:    conf.Dir,
		Stdin:  input,
		Stdout: c.out.writer(),
		Stderr: c.err.writer(),
	}

	return &c
}

// Start ...
func (c *Command) Start() {
	go c.start()
}

// Conf ...
func (c *Command) Conf() configuration.Processor {
	return c.conf
}

// Outs ...
func (c *Command) Outs() []io.Reader {
	return []io.Reader{c.out.reader()}
}

// Err ...
func (c *Command) Err() io.Reader {
	return c.err.reader()
}

func (c *Command) start() {
	var err error
	err = c.cmd.Start()
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = c.cmd.Wait()
	if err != nil {
		log.Fatalln(err.Error())
	}
	c.stop()
}

func (c *Command) stop() {
	errs := c.out.close()
	if len(errs) > 0 {
		s := strings.Join(errsToStrs(errs), ", ")
		log.Fatalf("Errors closing pipe: %s\n", s)
	}
}
