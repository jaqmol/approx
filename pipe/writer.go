package pipe

// Writer ...
type Writer struct {
	inputChannel chan []byte
}

// NewWriter ...
func NewWriter() *Writer {
	w := Writer{
		inputChannel: make(chan []byte),
	}
	return &w
}

func (w *Writer) Write(b []byte) (n int, err error) {
	w.inputChannel <- b
	return len(b), nil
}

// Channel ...
func (w *Writer) Channel() chan<- []byte {
	return w.inputChannel
}
