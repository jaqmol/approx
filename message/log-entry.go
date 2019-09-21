package message

import (
	"fmt"
	"io"
	"log"

	"github.com/jaqmol/approx/utils"
)

// LogEntry ...
type LogEntry struct {
	procName string
	msg      *Message
}

// NewLogEntry ...
func NewLogEntry(procName string, msg *Message) *LogEntry {
	return &LogEntry{
		procName: procName,
		msg:      msg,
	}
}

// LogTo ...
func (l *LogEntry) LogTo(w io.Writer) {
	s := fmt.Sprintf(
		"%v:%v %v\n",
		l.procName,
		l.msg.Role,
		string(utils.Truncated(l.msg.Data, 80)),
	)
	b := []byte(s)
	_, err := w.Write(b)
	if err != nil {
		log.Fatalf("Couln't log %v", b)
	}
}
