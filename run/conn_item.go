package run

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/jaqmol/approx/axmsg"

	"github.com/jaqmol/approx/flow"
)

// ConnItem ...
type ConnItem struct {
	errMsg       *axmsg.Errors
	pipeBasePath string
	FlowConn     *flow.ConnItem
	PipeName     string
	_file        *os.File
	reader       *bufio.Reader
	writer       *bufio.Writer
}

// NewConnItem ...
func NewConnItem(errMsg *axmsg.Errors, pipeBasePath string, flowConn *flow.ConnItem) (*ConnItem, error) {
	c := &ConnItem{
		errMsg:       errMsg,
		pipeBasePath: pipeBasePath,
		FlowConn:     flowConn,
		PipeName:     fmt.Sprintf("%v.pipe", flowConn.Hash),
	}
	err := c.createPipe()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *ConnItem) createPipe() error {
	return syscall.Mkfifo(c.PipePath(), 0600)
}

// PipePath ...
func (c *ConnItem) PipePath() string {
	return filepath.Join(c.pipeBasePath, c.PipeName)
}

// Cleanup ...
func (c *ConnItem) Cleanup() error {
	return os.Remove(c.PipePath())
}

// Reader ...
func (c *ConnItem) Reader() *bufio.Reader {
	if c.reader == nil {
		c.reader = bufio.NewReader(c.file())
	}
	return c.reader
}

// Writer ...
func (c *ConnItem) Writer() *bufio.Writer {
	if c.writer == nil {
		c.writer = bufio.NewWriter(c.file())
	}
	return c.writer
}

func (c *ConnItem) file() *os.File {
	if c._file == nil {
		var err error
		c._file, err = os.OpenFile(c.PipePath(), os.O_RDWR, 0600)
		if err != nil {
			c.errMsg.LogFatal(
				nil,
				"Can't open pipe @ %v: %v",
				c.PipePath(),
				err.Error(),
			)
		}
	}
	return c._file
}
