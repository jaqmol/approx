package message

// Reader ...
type Reader struct {
	messagesChan chan *Message
	envBuff      *EnvelopeBuffer
}

// NewReader ...
func NewReader(chunks <-chan []byte) *Reader {
	r := Reader{
		messagesChan: make(chan *Message),
		envBuff:      NewEnvelopeBuffer(chunks),
	}
	go r.start()
	return &r
}

// Messages ...
func (r *Reader) Messages() <-chan *Message {
	return r.messagesChan
}

func (r *Reader) start() {
	for env := range r.envBuff.Envelopes() {
		msg, err := ParseMessage(env.MessageBytes())
		catch(err)
		r.messagesChan <- msg
	}
}
