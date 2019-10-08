package test

import (
	"bytes"
	"testing"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/logger"
)

// TestLoggerWithSingleReader ...
func TestLoggerWithSingleReader(t *testing.T) {
	originals := loadTestData()
	originalForID := makePersonForIDMap(originals)
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.MsgEndBytes)
	originalCombined = append(originalCombined, configuration.MsgEndBytes...)
	reader := bytes.NewReader(originalCombined)

	writer := newTestWriter()
	l := logger.NewLogger(writer)
	l.Add(reader)
	go l.Start()

	count := 0
	for b := range writer.lines {
		checkTestSet(t, originalForID, b)
		count++
		writer.stop(count == len(originalBytes))
	}

	if len(originals) != count {
		t.Fatal("Logged line count doesn't corespond to received ones")
	}
}

func TestLoggerWithMultipleReaders(t *testing.T) {
	originals := loadTestData()
	originalForID := makePersonForIDMap(originals)
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.MsgEndBytes)
	originalCombined = append(originalCombined, configuration.MsgEndBytes...)

	writer := newTestWriter()
	l := logger.NewLogger(writer)

	for i := 0; i < 5; i++ {
		reader := bytes.NewReader(originalCombined)
		l.Add(reader)
	}

	go l.Start()
	goal := 5 * len(originals)
	count := 0

	for b := range writer.lines {
		checkTestSet(t, originalForID, b)
		count++
		writer.stop(count == goal)
	}

	if goal != count {
		t.Fatal("Logged line count doesn't corespond to received ones")
	}
}
