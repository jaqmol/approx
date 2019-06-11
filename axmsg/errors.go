package axmsg

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"os"
)

// Errors ...
type Errors struct {
	Source string
}

// LogFatal ...
func (e *Errors) LogFatal(id *int, format string, a ...interface{}) {
	e.Log(id, format, a...)
	os.Exit(1)
}

// Log ...
func (e *Errors) Log(id *int, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	msg := newErrorAction(e.Source, id, e.code(format), message)
	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Error marshalling error message: %v\n", message)
	}
	bytes = append(bytes, '\n')
	os.Stderr.Write(bytes)
}

// newErrorAction ...
func newErrorAction(source string, id *int, code uint32, message string) *Action {
	return NewAction(
		id,
		nil,
		"error",
		nil,
		ErrorData{
			Source:  source,
			Code:    code,
			Message: message,
		},
	)
}
func (e *Errors) code(format string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(e.Source))
	h.Write([]byte(format))
	return h.Sum32()
}

// ErrorData ...
type ErrorData struct {
	Source  string `json:"source"`
	Code    uint32 `json:"code"`
	Message string `json:"message"`
}
