package processor

import "io"

type procPipe struct {
	pipeReader *io.PipeReader
	pipeWriter *io.PipeWriter
}

func newProcPipe() procPipe {
	r, w := io.Pipe()
	return procPipe{
		pipeReader: r,
		pipeWriter: w,
	}
}

func (p *procPipe) reader() io.Reader {
	return p.pipeReader
}

func (p *procPipe) writer() io.Writer {
	return p.pipeWriter
}

func (p *procPipe) close() []error {
	acc := make([]error, 0)
	var err error

	err = p.pipeWriter.Close()
	if err != nil {
		acc = append(acc, err)
	}

	err = p.pipeReader.Close()
	if err != nil {
		acc = append(acc, err)
	}
	return acc
}
