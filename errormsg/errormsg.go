package errormsg

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// LogFatal ...
func LogFatal(processor string, id *int, code int, format string, a ...interface{}) {
	Log(processor, id, code, format, a...)
	os.Exit(1)
}

// Log ...
func Log(processor string, id *int, code int, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	msg := NewErrorMsg(processor, id, code, message)
	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Error marshalling error message: %v\n", message)
	}
	bytes = append(bytes, '\n')
	os.Stderr.Write(bytes)
}

// NewErrorMsg ...
func NewErrorMsg(processor string, id *int, code int, message string) *ErrorMsg {
	return &ErrorMsg{
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

// ErrorMsg ...
type ErrorMsg struct {
	JSONRPC string `json:"jsonrpc"`
	ID      *int   `json:"id"`
	Error   Error  `json:"error"`
}

// Error ...
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Data   `json:"data,omitempty"`
}

// Data ...
type Data struct {
	Processor string `json:"processor"`
}
