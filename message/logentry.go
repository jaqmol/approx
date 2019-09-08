package message

import (
	"fmt"
	"io"
)

// LogEntry ...
type LogEntry struct {
	Source  string
	Message string
}

// WriteTo ...
func (e *LogEntry) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w, "%v: %v\n", e.Source, e.Message)
	return int64(n), err
}
