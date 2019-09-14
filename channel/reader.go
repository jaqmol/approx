package channel

// Reader ...
type Reader interface {
	Read() <-chan []byte
}

// // Reader ...
// type Reader struct {
// 	outputChannel chan []byte
// 	buffer        []byte
// }

// // NewReader ...
// func NewReader() *Reader {
// 	r := Reader{
// 		outputChannel: make(chan []byte),
// 	}
// 	return &r
// }

// // Read ...
// func (r *Reader) Read(dst []byte) (n int, err error) {
// 	if len(r.buffer) > 0 {
// 		return r.readFromBuffer(dst)
// 	}
// 	return r.readFromChannel(dst)
// }

// func (r *Reader) readFromBuffer(dst []byte) (n int, err error) {
// 	copiedCount := copy(dst, r.buffer)
// 	if copiedCount < len(r.buffer) {
// 		r.buffer = r.buffer[copiedCount:]
// 	} else {
// 		r.buffer = nil
// 	}
// 	return copiedCount, nil
// }

// func (r *Reader) readFromChannel(dst []byte) (n int, err error) {
// 	src := <-r.outputChannel
// 	copiedCount := copy(dst, src)
// 	if copiedCount < len(src) {
// 		r.buffer = nil
// 		r.buffer = append(r.buffer, src[copiedCount:]...)
// 	}
// 	return copiedCount, nil
// }

// // Channel ...
// func (r *Reader) Channel() <-chan []byte {
// 	return r.outputChannel
// }
