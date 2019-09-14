package channel

// WriterWrap ...
type WriterWrap struct {
	writer Writer
}

// NewWriterWrap ...
func NewWriterWrap(w Writer) *WriterWrap {
	ww := WriterWrap{writer: w}
	return &ww
}

func (w *WriterWrap) Write(p []byte) (n int, err error) {
	w.writer.Write() <- p
	return len(p), nil
}
