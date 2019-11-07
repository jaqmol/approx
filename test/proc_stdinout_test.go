package test

import (
	"bytes"
	"os"
	"testing"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/processor"
)

// TestProcStdin ...
func TestProcStdin(t *testing.T) {
	originals := loadTestData()
	originalForID := makePersonForIDMap(originals)
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)
	reader := bytes.NewReader(originalCombined)

	err := os.Setenv("APPROX_ENV", "development")
	if err != nil {
		t.Fatal(err)
	}
	processor.DebugChangeStdin(reader)

	stdin := processor.Stdin

	serialize := make(chan []byte)
	go readFromReader(serialize, stdin.Out())
	stdin.Start()

	totalCount := 0
	countForID := make(map[string]int, 0)
	goal := len(originals)

	for b := range serialize {
		parsed := checkTestSet(t, originalForID, b)
		totalCount++
		countForID[parsed.ID]++
		if totalCount == goal {
			close(serialize)
		}
	}

	if goal != totalCount {
		t.Fatal("Stdin count doesn't corespond to multitude of source count")
	}

	if len(originals) != len(countForID) {
		t.Fatal("Received individual data sets count doesn't corespond source count")
	}

	processor.DebugResetStdin()
}

// TestProcStdout ...
func TestProcStdout(t *testing.T) {
	originals := loadTestData()
	originalForID := makePersonForIDMap(originals)
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)
	reader := bytes.NewReader(originalCombined)

	err := os.Setenv("APPROX_ENV", "development")
	if err != nil {
		t.Fatal(err)
	}

	stdout := processor.Stdout
	stdout.Connect(reader)
	stdout.Start()
	// log.Println("Did start stdout") TODO: REMOVE

	writer := newTestWriter()
	processor.DebugChangeStdout(writer)

	totalCount := 0
	countForID := make(map[string]int, 0)
	goal := len(originals)

	// log.Println("Starting to read writer lines") TODO: REMOVE
	for raw := range writer.lines {
		b := processor.ClearEventEnd(raw)
		// log.Println("Did read", len(b), "bytes from writer") TODO: REMOVE
		parsed := checkTestSet(t, originalForID, b)
		totalCount++
		countForID[parsed.ID]++
		writer.stop(totalCount == goal)
	}

	if goal != totalCount {
		t.Fatal("Stdin count doesn't corespond to multitude of source count")
	}

	if len(originals) != len(countForID) {
		t.Fatal("Received individual data sets count doesn't corespond source count")
	}

	processor.DebugResetStdout()
}
