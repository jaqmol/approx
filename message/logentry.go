package message

import (
	"encoding/json"
	"fmt"
	"io"
)

// Log entries are send via stderr and are used for error and other stage logging.
// To indicate this, log entry types are passed on via the "cmd" property.
// The type "inform" indicates an info-level logging message.
// The type "warn" indicates an warning-level logging message.
// The type "fail" indicates a runtime error that's expected / a recoverable.
// The type "exit" indicates a non-recoverable exception, that's not expected.
//   Approx will exit itself if a log message with command "exit" ist received.

// SourcedLogEntry ...
type SourcedLogEntry struct {
	// ID helps to identify relations between messages, for instance if a response (role) has the same ID as a request (role)
	ID string `json:"id"`
	// Role is the context in which the message is supposed to be interpreted, like request, error or response
	Role string `json:"role,omitempty"`
	// Cmd can be used to indicate which action is expected on behalf of a message
	Cmd string `json:"cmd,omitempty"`
	// Index might indicate the order of a message in a stream of related chunks
	Index *int `json:"index,omitempty"`
	// If fork is used for parallelizing, sequence indicates on which parallel code path (sequence) a message is running
	Sequence *int `json:"sequence,omitempty"`
	// Represents the transport payload of a message, binary data is represented as a base-64 string
	Payload SourcedLogEntryPayload `json:"payload"`
}

// WriteTo ...
func (sem *SourcedLogEntry) WriteTo(w io.Writer) (int64, error) {
	bytes, err := json.Marshal(sem)
	if err != nil {
		panic(err)
	}
	toWrite := append(bytes, []byte("\n")...)
	i, err := w.Write(toWrite)
	return int64(i), err
}

// ToSourcedLogEntry ...
func (m *Message) ToSourcedLogEntry(source string) *SourcedLogEntry {
	return &SourcedLogEntry{
		ID:       m.ID,
		Role:     m.Role,
		Cmd:      m.Cmd,
		Index:    m.Index,
		Sequence: m.Sequence,
		Payload: SourcedLogEntryPayload{
			Processor: source,
			Message:   string(*m.Payload),
		},
	}
}

// MakeSourcedLogEntry ...
func MakeSourcedLogEntry(processor string, id string, eType LogEntryType, message string) *SourcedLogEntry {
	return &SourcedLogEntry{
		ID:   id,
		Role: "log",
		Cmd:  StringForLogEntryType[eType],
		Payload: SourcedLogEntryPayload{
			Processor: processor,
			Message:   message,
		},
	}
}

// SourcedLogEntryPayload ...
type SourcedLogEntryPayload struct {
	Processor string `json:"processor"`
	Message   string `json:"message"`
}

// LogEntryType ...
type LogEntryType int

// Error-Levels, also useful in logging
const (
	Inform LogEntryType = iota
	Warn
	Fail
	Exit
)

// StringForLogEntryType ...
var StringForLogEntryType = map[LogEntryType]string{
	Inform: "inform",
	Warn:   "warn",
	Fail:   "fail",
	Exit:   "exit",
}

// LogEntryTypeForString ...
var LogEntryTypeForString = map[string]LogEntryType{
	"inform": Inform,
	"warn":   Warn,
	"fail":   Fail,
	"exit":   Exit,
}

// NewLogEntry ...
func NewLogEntry(eType LogEntryType, id string, msgStr string) *Message {
	payloadBytes := append([]byte("\""), append([]byte(msgStr), []byte("\"")...)...)
	payload := json.RawMessage(payloadBytes)
	return &Message{
		ID:      id,
		Role:    "log",
		Cmd:     StringForLogEntryType[eType],
		Payload: &payload,
	}
}

// WriteLogEntry ...
func WriteLogEntry(w io.Writer, eType LogEntryType, id string, msgStr string) {
	msg := NewLogEntry(eType, id, msgStr)
	bytes, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	toWrite := append(bytes, []byte("\n")...)
	_, err = w.Write(toWrite)
	if err != nil {
		panic(err)
	}
}

// WriteLogEntryF ...
func WriteLogEntryF(w io.Writer, eType LogEntryType, id string, format string, values ...interface{}) {
	msgStr := fmt.Sprintf(format, values...)
	WriteLogEntry(w, eType, id, msgStr)
}

// WriteSourcedLogEntry ...
func WriteSourcedLogEntry(w io.Writer, eType LogEntryType, processor string, id string, msgStr string) {
	msg := NewLogEntry(eType, id, msgStr)
	sourcedMsg := msg.ToSourcedLogEntry(processor)
	_, err := sourcedMsg.WriteTo(w)
	if err != nil {
		panic(err)
	}
}
