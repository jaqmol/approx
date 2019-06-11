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

func (r *Reader) Read() (*Action, json.RawMessage, error) {
	bytes, err := r.reader.ReadBytes('\n')
	if err != nil {
		return nil, nil, err
	}
	return ActionAndData(bytes)
}
