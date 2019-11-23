package actor

import (
	"bytes"
	"log"
	"path/filepath"
	"testing"

	"github.com/jaqmol/approx/config"
	"github.com/jaqmol/approx/logging"
	"github.com/jaqmol/approx/test"
)

// TestCommandWithBufferProcessing ...
func TestCommandWithBufferProcessing(t *testing.T) {
	performTestWithIdentCmdAndArgs(t, "buffer-cmd", "node", "node-procs/test-buffer-processing.js")
}

// TestCommandWithJSONProcessing ...
func TestCommandWithJSONProcessing(t *testing.T) {
	performTestWithIdentCmdAndArgs(t, "json-cmd", "node", "node-procs/test-json-processing.js")
}

// TestCommandWithBufferProcessingWithLogging ...
func TestCommandWithBufferProcessingWithLogging(t *testing.T) {
	// t.SkipNow()
	performTestWithIdentCmdAndArgsAndLogging(t, "buffer-cmd", "node", "node-procs/test-buffer-processing.js")
}

func performTestWithIdentCmdAndArgs(t *testing.T, ident, cmd, arg string) {
	originals := test.LoadTestData()
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, config.EvntEndBytes)
	originalCombined = append(originalCombined, config.EvntEndBytes...)

	producer := NewThrottledProducer(10, 5000)
	command := NewCommand(10, ident, cmd, arg)

	testDir, err := filepath.Abs("../test")
	if err != nil {
		t.Fatal(err)
	}
	command.Directory(testDir)

	receiver := make(chan unifiedMessage, 10)
	collector := NewCollector(10)

	producer.Next(command)
	command.Next(collector)
	dumpCommandLogMessages(command)

	startCollectingUnifiedDataMessages(t, collector, receiver, func() {
		close(receiver)
	})
	command.Start()
	startProducingTestMessages(t, producer, originalCombined)

	counter := 0
	for unimsg := range receiver {
		if unimsg.messageType == unifiedMsgDataType {
			parsed, err := test.UnmarshalPerson(unimsg.data)
			if err != nil {
				t.Fatal(err)
			}

			original := originalForID[parsed.ID]

			test.CheckUpperFirstAndLastNames(t, &original, parsed)
			counter++
		}
	}

	if counter != len(originals) {
		t.Fatalf("%v messages expected to be produced, but got %v", len(originals), counter)
	}
}

func performTestWithIdentCmdAndArgsAndLogging(t *testing.T, ident, cmd, arg string) {
	originals := test.LoadTestData()
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, config.EvntEndBytes)
	originalCombined = append(originalCombined, config.EvntEndBytes...)

	producer := NewThrottledProducer(10, 5000)
	command := NewCommand(10, ident, cmd, arg)

	testDir, err := filepath.Abs("../test")
	if err != nil {
		t.Fatal(err)
	}
	command.Directory(testDir)

	receiver := make(chan unifiedMessage, 10)
	collector := NewCollector(10)

	producer.Next(command)
	command.Next(collector)

	// Logger
	logReceiver := make(chan []byte, 10)
	logger := logging.NewChannelLog(logReceiver)
	funnelIntoUnifiedLogMessages(logReceiver, receiver)
	logger.Add(command.Logging())
	// /Logger

	startCollectingUnifiedDataMessages(t, collector, receiver, func() {
		close(receiver)
	})
	command.Start()
	go logger.Start() // Logger
	startProducingTestMessages(t, producer, originalCombined)

	counter := 0
	for unimsg := range receiver {
		if unimsg.messageType == unifiedMsgDataType {
			parsed, err := test.UnmarshalPerson(unimsg.data)
			if err != nil {
				t.Fatal(err)
			}

			original := originalForID[parsed.ID]

			test.CheckUpperFirstAndLastNames(t, &original, parsed)
			counter++
		} else if unimsg.messageType == unifiedMsgLogType {
			log.Println(string(unimsg.data))
		}
	}

	if counter != len(originals) {
		t.Fatalf("%v messages expected to be produced, but got %v", len(originals), counter)
	}
}

const unifiedMsgDataType = 1
const unifiedMsgLogType = 2

type unifiedMessage struct {
	messageType int
	data        []byte
}

func startCollectingUnifiedDataMessages(
	t *testing.T,
	collector *Collector,
	receiver chan<- unifiedMessage,
	finished func(),
) {
	go func() {
		err := collector.Collect(func(message []byte) error {
			receiver <- unifiedMessage{
				messageType: unifiedMsgDataType,
				data:        message,
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		finished()
	}()
}

func funnelIntoUnifiedLogMessages(
	logReceiver <-chan []byte,
	receiver chan<- unifiedMessage,
) {
	go func() {
		for message := range logReceiver {
			receiver <- unifiedMessage{
				messageType: unifiedMsgLogType,
				data:        message,
			}
		}
	}()
}

func dumpLogMessagesOfCommand(c *Command) {
	rec := make(chan []byte, 10)
	l := logging.NewChannelLog(rec)
	l.Add(c.Logging())
	go func() {
		for len(rec) > 0 {
			<-rec
		}
	}()
	go l.Start()
}

func dumpCommandLogMessages(cmd *Command) {
	logReceiver := make(chan []byte, 10)
	logger := logging.NewChannelLog(logReceiver)
	dumpByteChannelEvents(logReceiver)
	logger.Add(cmd.Logging())
	go logger.Start()
}

func dumpByteChannelEvents(channel <-chan []byte) {
	go func() {
		ok := true
		for ok {
			_, ok = <-channel
		}
	}()
}