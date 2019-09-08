package pipe

// Reader ...
type Reader struct {
	outputChannel chan []byte
	buffer        []byte
}

// NewReader ...
func NewReader() *Reader {
	r := Reader{
		outputChannel: make(chan []byte),
	}
	return &r
}

// Read ...
func (r *Reader) Read(dst []byte) (n int, err error) {
	var src []byte
	if len(r.buffer) > 0 {
		src = r.buffer
		r.buffer = []byte{}
	} else {
		src = <-r.outputChannel
	}
	copiedCount := copy(dst, src)
	if copiedCount < len(src) {
		r.buffer = src[copiedCount:]
	}
	return copiedCount, nil
}

// Channel ...
func (r *Reader) Channel() <-chan []byte {
	return r.outputChannel
}
