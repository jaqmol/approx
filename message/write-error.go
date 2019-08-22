package message

import (
	"encoding/json"
	"io"
)

// WriteError ...
func WriteError(w io.Writer, id string, errStr string) {
	msg := Message{
		ID:      id,
		Role:    "error",
		Payload: []byte(errStr),
	}
	bytes, err := json.Marshal(&msg)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(bytes)
	if err != nil {
		panic(err)
	}
}
