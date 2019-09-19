package message

// Reader ...
type Reader struct {
	messagesChan chan *Message
	cache        *chunkCache
}

// NewReader ...
func NewReader(bytes <-chan []byte) *Reader {
	r := Reader{
		messagesChan: make(chan *Message),
		cache:        newChunkCache(bytes),
	}
	go r.start()
	return &r
}

// Messages ...
func (r *Reader) Messages() <-chan *Message {
	return r.messagesChan
}

func (r *Reader) start() {
	for msgBytes := range r.cache.messageBytes() {
		msg, err := ParseNessage(msgBytes)
		catch(err)
		r.messagesChan <- msg
	}
}
