package errormsg

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"os"
)

// ErrorMsg ...
type ErrorMsg struct {
	Processor string
}

func (e *ErrorMsg) code(format string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(e.Processor))
	h.Write([]byte(format))
	return h.Sum32()
}

// LogFatal ...
func (e *ErrorMsg) LogFatal(id *int, format string, a ...interface{}) {
	e.Log(id, format, a...)
	os.Exit(1)
}

// Log ...
func (e *ErrorMsg) Log(id *int, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	msg := newJSONRPCErrorMsg(e.Processor, id, e.code(format), message)
	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Error marshalling error message: %v\n", message)
	}
	bytes = append(bytes, '\n')
	os.Stderr.Write(bytes)
}

// newJSONRPCErrorMsg ...
func newJSONRPCErrorMsg(processor string, id *int, code uint32, message string) *JSONRPCErrorMsg {
	return &JSONRPCErrorMsg{
		JSONRPC: "2.0",
		ID:      id,
		Error: Error{
			Code:    code,
			Message: message,
			Data: Data{
				Processor: processor,
			},
		},
	}
}

// JSONRPCErrorMsg ...
type JSONRPCErrorMsg struct {
	JSONRPC string `json:"jsonrpc"`
	ID      *int   `json:"id"`
	Error   Error  `json:"error"`
}

// Error ...
type Error struct {
	Code    uint32 `json:"code"`
	Message string `json:"message"`
	Data    Data   `json:"data,omitempty"`
}

// Data ...
type Data struct {
	Processor string `json:"processor"`
}
