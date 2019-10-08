package test

import (
	"bufio"
	"bytes"
	"log"
	"testing"

	"github.com/jaqmol/approx/event"
	"github.com/jaqmol/approx/processor"

	"github.com/jaqmol/approx/configuration"
)

// // TestCommand ...
// func TestCommand(t *testing.T) {
// 	originals := loadTestData()
// 	originalForID := makePersonForIDMap(originals)
// 	originalBytes := marshallPeople(originals)

// 	originalCombined := bytes.Join(originalBytes, configuration.MsgEndBytes)
// 	originalCombined = append(originalCombined, configuration.MsgEndBytes...)

// 	reader := bytes.NewReader(originalCombined)
// 	conf := configuration.Command{
// 		Ident: "test-command",
// 		Cmd:   "node command.js",
// 	}
// 	command := processor.NewCommand(&conf, reader)

// 	totalCount := 0
// 	outputReader := command.Outs()[0]
// 	command.Start()
// 	scanner := event.NewScanner(outputReader)

// 	for scanner.Scan() {
// 		raw := scanner.Bytes()
// 		data := bytes.Trim(raw, "\x00")
// 		checkTestSet(t, originalForID, data)

// 		totalCount++
// 		if totalCount == len(originals) {
// 			command.SigInt()
// 		}
// 	}

// 	if len(originals) != totalCount {
// 		t.Fatal("Command dispatch count doesn't corespond to source count")
// 	}
// }

// TestCommand ...
func TestCommand(t *testing.T) {
	// t.SkipNow()
	originals := loadTestData()[:10]
	originalForID := makePersonForIDMap(originals)
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.MsgEndBytes)
	originalCombined = append(originalCombined, configuration.MsgEndBytes...)

	conf := configuration.Command{
		Ident: "test-command",
		Cmd:   "node command.js",
	}

	inReader := bytes.NewReader(originalCombined)
	command := processor.NewCommand(&conf, inReader)
	outScanner := newCommandOutScanner(command)

	go outScanner.start()
	command.Start()
	log.Println("Did start command and scanner")

	totalCount := 0
	for data := range outScanner.lines {
		parsed := checkTestSet(t, originalForID, data)
		log.Printf("parsed: %v\n", parsed)

		totalCount++
		if totalCount == len(originals) {
			// command.Stop()
			close(outScanner.lines)
		}
	}
	if totalCount != len(originals) {
		t.Fatalf("Command didn't dispatch all events: %v/%v\n", totalCount, len(originals))
	}
	log.Printf("    totalCount: %v\n", totalCount)
	log.Printf("len(originals): %v\n", len(originals))
}

// TestCommandLogging ...
func TestCommandLogging(t *testing.T) {
	t.SkipNow()
	originals := loadTestData()[:10]
	// originalForID := makePersonForIDMap(originals)
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.MsgEndBytes)
	originalCombined = append(originalCombined, configuration.MsgEndBytes...)

	conf := configuration.Command{
		Ident: "test-command",
		Cmd:   "node command-2.js",
	}

	inReader := bytes.NewReader(originalCombined)
	command := processor.NewCommand(&conf, inReader)
	// logWriter := newTestWriter()
	logScanner := bufio.NewScanner(command.Err())

	command.Start()

	for logScanner.Scan() {
		log.Printf("%v\n", logScanner.Text())
	}

	// lggr := logger.NewLogger(logWriter) // os.Stdout)

	// lggr.Add(command.Err())
	// go lggr.Start()

	// totalCount := 0

	// for data := range logWriter.lines {
	// 	log.Printf("Received via log writer: %v\n", string(data))
	// 	// checkTestSet(t, originalForID, data)

	// 	totalCount++
	// 	logWriter.stop(totalCount == len(originals))
	// }
}

/*
logWriter := newTestWriter()

	lggr := logger.NewLogger(logWriter)
	l.Add(command.Err())
	go l.Start()

	totalCount := 0
	loggedCount := 0

	command.Start()

	loop := true
	checkLoop := func() {
		loop = totalCount < len(originals) || loggedCount < len(originals)
	}
	for loop {
		select {
		case raw := <-scannedLines:
			data := bytes.Trim(raw, "\x00")
			if len(data) > 0 {
				checkTestSet(t, originalForID, data)
				totalCount++
				if totalCount == len(originals) {
					command.SigInt()
				}
			}
			checkLoop()
		case loggedLine := <-logWriter.lines:
			log.Printf("LOGGED: %v\n", string(loggedLine))
			loggedCount++
			logWriter.stop(loggedCount == len(originals))
			checkLoop()
		}
	}

	if len(originals) != totalCount {
		t.Fatal("Command dispatch count doesn't corespond to source count")
	}
*/

type commandOutScanner struct {
	scanner *bufio.Scanner
	lines   chan []byte
}

func newCommandOutScanner(cmd *processor.Command) *commandOutScanner {
	return &commandOutScanner{
		scanner: event.NewScanner(cmd.Outs()[0]),
		lines:   make(chan []byte),
	}
}

func (s *commandOutScanner) start() {
	log.Println("Start scanning")
	for s.scanner.Scan() {
		raw := s.scanner.Bytes()
		data := bytes.Trim(raw, "\x00")
		cp := make([]byte, len(data))
		copy(cp, data)
		s.lines <- cp
	}
	log.Println("End scanning")
	close(s.lines)
}
