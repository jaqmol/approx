package channel

// ReaderWrap ...
type ReaderWrap struct {
	reader Reader
	buffer []byte
}

// NewReaderWrap ...
func NewReaderWrap(r Reader) *ReaderWrap {
	rw := ReaderWrap{reader: r}
	return &rw
}

// Read ...
func (r *ReaderWrap) Read(dst []byte) (n int, err error) {
	if len(r.buffer) > 0 {
		return r.readFromBuffer(dst)
	}
	return r.readFromChannel(dst)
}

func (r *ReaderWrap) readFromBuffer(dst []byte) (n int, err error) {
	copiedCount := copy(dst, r.buffer)
	if copiedCount < len(r.buffer) {
		r.buffer = r.buffer[copiedCount:]
	} else {
		r.buffer = nil
	}
	return copiedCount, nil
}

func (r *ReaderWrap) readFromChannel(dst []byte) (n int, err error) {
	src := <-r.reader.Read()
	copiedCount := copy(dst, src)
	if copiedCount < len(src) {
		r.buffer = nil
		r.buffer = append(r.buffer, src[copiedCount:]...)
	}
	return copiedCount, nil
}
