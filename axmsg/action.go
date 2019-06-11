package axmsg

import "encoding/json"

// NewAction ...
func NewAction(id *int, seq *int, role string, cmd *string, data interface{}) *Action {
	return &Action{
		AXMSG: 1,
		ID:    id,
		Seq:   seq,
		Role:  role,
		Cmd:   cmd,
		Data:  data,
	}
}

// ActionAndData ...
func ActionAndData(bytes []byte) (*Action, json.RawMessage, error) {
	var data json.RawMessage
	action := Action{Data: &data}
	err := json.Unmarshal(bytes, &action)
	return &action, data, err
}

// Action ...
type Action struct {
	AXMSG int         `json:"axmsg"`
	ID    *int        `json:"id,omitempty"`
	Seq   *int        `json:"seq,omitempty"`
	Role  string      `json:"role"`
	Cmd   *string     `json:"cmd,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

// ActionBytes ...
type ActionBytes []byte

// Bytes ...
func (a *Action) Bytes() (ActionBytes, error) {
	return json.Marshal(a)
}

// NewlineSuffix ...
func (ab ActionBytes) NewlineSuffix() []byte {
	return append(ab, '\n')
}
