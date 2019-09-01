package message

import (
	"encoding/json"
	"fmt"
	"io"
)

// Error messages are not only used for errors but for logging as well.
// To indicate this, error types are passed on via the "cmd" property.
// The type "inform" indicates an info-level logging message.
// The type "warn" indicates an warning-level logging message.
// The type "fail" indicates a runtime error that's expected / a recoverable.
// The type "exit" indicates a non-recoverable exception, that's not expected.

// SourcedErrorMessage ...
type SourcedErrorMessage struct {
	// ID helps to identify relations between messages, for instance if a response (role) has the same ID as a request (role)
	ID string `json:"id"`
	// Role is the context in which the message is supposed to be interpreted, like request, error or response
	Role string `json:"role,omitempty"`
	// Cmd can be used to indicate which action is expected on behalf of a message
	Cmd string `json:"cmd,omitempty"`
	// Index might indicate the order of a message in a stream of related chunks
	Index int `json:"index,omitempty"`
	// If fork is used for parallelizing, sequence indicates on which parallel code path (sequence) a message is running
	Sequence int `json:"sequence,omitempty"`
	// Represents the transport payload of a message, binary data is represented as a base-64 string
	Payload SourcedErrorMessagePayload `json:"payload"`
}

// WriteTo ...
func (sem *SourcedErrorMessage) WriteTo(w io.Writer) (int64, error) {
	bytes, err := json.Marshal(sem)
	if err != nil {
		panic(err)
	}
	i, err := w.Write(bytes)
	return int64(i), err
}

// ToSourcedErrorMessage ...
func (m *Message) ToSourcedErrorMessage(source string) *SourcedErrorMessage {
	return &SourcedErrorMessage{
		ID:       m.ID,
		Role:     m.Role,
		Cmd:      m.Cmd,
		Index:    m.Index,
		Sequence: m.Sequence,
		Payload: SourcedErrorMessagePayload{
			Processor: source,
			Message:   string(*m.Payload),
		},
	}
}

// SourcedErrorMessagePayload ...
type SourcedErrorMessagePayload struct {
	Processor string `json:"processor"`
	Message   string `json:"message"`
}

// ErrorType ...
type ErrorType int

// Error-Levels, also useful in logging
const (
	Inform ErrorType = iota
	Warn
	Fail
	Exit
)

// StringForErrorType ...
var StringForErrorType = map[ErrorType]string{
	Inform: "inform",
	Warn:   "warn",
	Fail:   "fail",
	Exit:   "exit",
}

// ErrorTypeForString ...
var ErrorTypeForString = map[string]ErrorType{
	"inform": Inform,
	"warn":   Warn,
	"fail":   Fail,
	"exit":   Exit,
}

// NewError ...
func NewError(eType ErrorType, id string, errStr string) *Message {
	payloadString := fmt.Sprintf("\"%v\"", errStr)
	payload := json.RawMessage(payloadString)
	return &Message{
		ID:      id,
		Role:    "error",
		Cmd:     StringForErrorType[eType],
		Payload: &payload,
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
