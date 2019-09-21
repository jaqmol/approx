package channel

// IoReadFromChan ...
type IoReadFromChan struct {
	readChan <-chan []byte
	buffer   []byte
}

// NewIoReadFromChan ...
func NewIoReadFromChan(r Reader) *IoReadFromChan {
	rw := IoReadFromChan{readChan: r.Read()}
	return &rw
}

// Read ...
func (r *IoReadFromChan) Read(dst []byte) (n int, err error) {
	if len(r.buffer) > 0 {
		return r.readFromBuffer(dst)
	}
	return r.readFromChannel(dst)
}

func (r *IoReadFromChan) readFromBuffer(dst []byte) (n int, err error) {
	copiedCount := copy(dst, r.buffer)
	if copiedCount < len(r.buffer) {
		r.buffer = r.buffer[copiedCount:]
	} else {
		r.buffer = nil
	}
	return copiedCount, nil
}

func (r *IoReadFromChan) readFromChannel(dst []byte) (n int, err error) {
	src := <-r.readChan
	copiedCount := copy(dst, src)
	if copiedCount < len(src) {
		r.buffer = nil
		r.buffer = append(r.buffer, src[copiedCount:]...)
	}
	return copiedCount, nil
}
