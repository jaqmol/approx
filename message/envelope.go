package message

import "fmt"

// Envelope ...
type Envelope struct {
	MessageLength int
	Bytes         []byte
}

// NewEnvelope creates envelope from header (must contain seperator) and data
func NewEnvelope(header []byte, data []byte) *Envelope {
	msgLen := len(header) + len(data)
	msgLenStr := fmt.Sprintf("%v:", msgLen)
	bytes := make([]byte, 0)
	bytes = append(bytes, []byte(msgLenStr)...)
	bytes = append(bytes, header...)
	bytes = append(bytes, data...)
	return &Envelope{
		MessageLength: msgLen,
		Bytes:         bytes,
	}
}

// MessageBytes ...
func (e *Envelope) MessageBytes() []byte {
	length := len(e.Bytes)
	start := length - e.MessageLength
	return e.Bytes[start:length]
}
