package message

import "encoding/json"

// Message ...
type Message struct {
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
	Payload json.RawMessage `json:"payload,omitempty"`
}
