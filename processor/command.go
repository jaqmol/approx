package processor

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/event"
)

// Command ...
type Command struct {
	conf   *configuration.Command
	cmd    *exec.Cmd
	in     io.Reader
	cmdIn  procPipe
	cmdOut procPipe
	cmdErr procPipe
}

// NewCommand ...
func NewCommand(conf *configuration.Command, input io.Reader) *Command {
	c := Command{
		conf:   conf,
		in:     input,
		cmdIn:  newProcPipe(),
		cmdOut: newProcPipe(),
		cmdErr: newProcPipe(),
	}

	cmd, args := cmdAndArgs(conf.Cmd)
	c.cmd = exec.Command(cmd, args...)
	c.cmd.Env = append(os.Environ(), conf.Env...)
	c.cmd.Stdin = c.cmdIn.reader()
	c.cmd.Stdout = c.cmdOut.writer()
	c.cmd.Stderr = c.cmdErr.writer()

	return &c
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

// // NewCommand ...
// func NewCommand(conf *configuration.Command, input io.Reader) *Command {
// 	c := Command{
// 		conf: conf,
// 		out:  newProcPipe(),
// 		err:  newProcPipe(),
// 	}

// 	c.cmd = &exec.Cmd{
// 		Path:   conf.Path,
// 		Args:   conf.Args,
// 		Env:    conf.Env,
// 		Dir:    conf.Dir,
// 		Stdin:  input,
// 		Stdout: c.out.writer(),
// 		Stderr: c.err.writer(),
// 	}

// 	return &c
// }

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

func (c *Command) startReadingInput() {
	scanner := event.NewScanner(c.in)

	log.Println("Command start reading input")

	for scanner.Scan() {
		msg := msgEndedCopy(scanner.Bytes())
		n, err := c.cmdIn.writer().Write(msg)
		if err != nil {
			log.Fatalln(err.Error())
		}
		if n != len(msg) {
			log.Fatalln("Command couldn't write complete event")
		}
	}

	log.Println("Command stop reading input")

	c.Stop()
}

func (c *Command) startCmd() {
	var err error
	err = c.cmd.Start()
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = c.cmd.Wait()
	if err != nil {
		log.Fatalln(err.Error())
	}
	c.Stop()
}

// Stop ...
func (c *Command) Stop() {
	errs := make([]error, 0)
	errs = append(errs, c.cmdOut.close()...)
	errs = append(errs, c.cmdErr.close()...)
	if len(errs) > 0 {
		s := strings.Join(errsToStrs(errs), ", ")
		log.Fatalf("Errors closing pipe: %s\n", s)
	}
}

// // SigInt ...
// func (c *Command) SigInt() {
// 	err := c.cmd.Process.Signal(os.Interrupt)
// 	if err != nil {
// 		// log.Println("SIGINT ERR: " + err.Error())
// 		err = c.cmd.Process.Kill()
// 		if err != nil {
// 			// log.Fatalln("KILL ERR: " + err.Error())
// 		}
// 	}
// 	c.Stop()
// }
