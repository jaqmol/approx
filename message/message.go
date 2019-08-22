package message

import "encoding/json"

// Message ...
type Message struct {
	ID      string          `json:"id"`
	Role    string          `json:"role,omitempty"`
	Cmd     string          `json:"cmd,omitempty"`
	Index   int             `json:"index,omitempty"`
	Payload json.RawMessage `json:"payload"`
}
