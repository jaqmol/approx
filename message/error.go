package message

import (
	"encoding/json"
	"io"
)

// NewError ...
func NewError(id string, errStr string) *Message {
	return &Message{
		ID:      id,
		Role:    "error",
		Payload: []byte(errStr),
	}
}

// WriteError ...
func WriteError(w io.Writer, id string, errStr string) {
	msg := NewError(id, errStr)
	bytes, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(bytes)
	if err != nil {
		panic(err)
	}
}
