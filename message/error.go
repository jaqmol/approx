package message

import (
	"encoding/json"
	"io"
)

// Error messages are not only used for errors but for logging as well.
// To indicate this, error types are passed on via the "cmd" property.
// The type "inform" indicates an info-level logging message.
// The type "warn" indicates an warning-level logging message.
// The type "fail" indicates a runtime error that's expected / a recoverable.
// The type "exit" indicates a non-recoverable exception, that's not expected.

// ErrorType ...
type ErrorType int

// Error-Levels, also useful in logging
const (
	Inform ErrorType = iota
	Warn
	Fail
	Exit
)

var stringForErrorType = map[ErrorType]string{
	Inform: "inform",
	Warn:   "warn",
	Fail:   "fail",
	Exit:   "exit",
}

var typeForErrorString = map[string]ErrorType{
	"inform": Inform,
	"warn":   Warn,
	"fail":   Fail,
	"exit":   Exit,
}

// NewError ...
func NewError(eType ErrorType, id string, errStr string) *Message {
	return &Message{
		ID:      id,
		Role:    "error",
		Cmd:     stringForErrorType[eType],
		Payload: []byte(errStr),
	}
}

// WriteError ...
func WriteError(w io.Writer, eType ErrorType, id string, errStr string) {
	msg := NewError(eType, id, errStr)
	bytes, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(bytes)
	if err != nil {
		panic(err)
	}
}

func stringFromErrorType(eType ErrorType) {

}
