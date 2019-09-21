package channel

// IoWriteToChan ...
type IoWriteToChan struct {
	writeChan chan<- []byte
}

// NewIoWriteToChan ...
func NewIoWriteToChan(w Writer) *IoWriteToChan {
	ww := IoWriteToChan{writeChan: w.Write()}
	return &ww
}

func (w *IoWriteToChan) Write(p []byte) (n int, err error) {
	w.writeChan <- p
	return len(p), nil
}
