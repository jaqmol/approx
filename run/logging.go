package run

import (
	"log"
	"os"

	"github.com/jaqmol/approx/channel"
	"github.com/jaqmol/approx/message"
)

// Logging ...
func Logging(stdErrs map[string]channel.Pipe) {
	logChan := make(chan *message.LogEntry)
	for procName, errPipe := range stdErrs {
		go listenForLogEntry(logChan, procName, errPipe)
	}
	for entry := range logChan {
		entry.LogTo(os.Stderr)
	}
}

func listenForLogEntry(logChan chan<- *message.LogEntry, procName string, errReader channel.Reader) {
	envReader := message.NewEnvelopeBuffer(errReader.Read())
	for env := range envReader.Envelopes() {
		msg, err := message.ParseMessage(env.MessageBytes())
		if err != nil {
			log.Fatal(err)
		}
		entry := message.NewLogEntry(procName, msg)
		logChan <- entry
	}
}
