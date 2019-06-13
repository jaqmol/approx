package axmsg

import "bufio"

// Writer ...
type Writer struct {
	writer *bufio.Writer
}

// NewWriter ...
func NewWriter(writer *bufio.Writer) *Writer {
	return &Writer{
		writer: writer,
	}
}

func (w *Writer) Write(action *Action) error {
	ab, err := action.Bytes()
	if err != nil {
		return err
	}
	_, err = w.writer.Write(ab.NewlineSuffix())
	if err != nil {
		return err
	}
	return w.writer.Flush()
}

// WriteBytes ...
func (w *Writer) WriteBytes(bytes []byte) error {
	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}
	return w.writer.Flush()
}