package axmsg

import (
	"bufio"
	"encoding/json"
)

// Reader ...
type Reader struct {
	reader *bufio.Reader
}

// NewReader ...
func NewReader(reader *bufio.Reader) *Reader {
	return &Reader{
		reader: reader,
	}
}

// NewReaders ...
func NewReaders(readers []bufio.Reader) []Reader {
	acc := make([]Reader, 0)
	for _, r := range readers {
		msgReader := NewReader(&r)
		acc = append(acc, *msgReader)
	}
	return acc
}

func (r *Reader) Read() (*Action, json.RawMessage, error) {
	bytes, err := r.ReadBytes()
	if err != nil {
		return nil, nil, err
	}
	return ActionAndData(bytes)
}

// ReadBytes ...
func (r *Reader) ReadBytes() ([]byte, error) {
	return r.reader.ReadBytes('\n')
}
