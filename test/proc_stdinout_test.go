package test

import (
	"bytes"
	"testing"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/processor"
)

// TestProcStdin ...
func TestProcStdin(t *testing.T) {
	t.SkipNow()
	originals := loadTestData()
	originalForID := makePersonForIDMap(originals)
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)
	reader := bytes.NewReader(originalCombined)

	stdin := processor.NewStdin()
	err := stdin.Connect(reader)
	if err != nil {
		t.Fatal(err)
	}

	collector, err := processor.NewCollector(stdin.Out())
	if err != nil {
		t.Fatal(err)
	}
	collector.Start()
	// serialize := make(chan []byte)
	// go readFromReader(serialize, stdin.Out())
	stdin.Start()

	totalCount := 0
	countForID := make(map[string]int, 0)
	goal := len(originals)

	for b := range collector.Events() {
		parsed := checkTestSet(t, originalForID, b)
		totalCount++
		countForID[parsed.ID]++
		if totalCount == goal {
			break
			// collector.Stop()
			// close(serialize)
		}
	}

	if goal != totalCount {
		t.Fatal("Stdin count doesn't corespond to multitude of source count")
	}

	if len(originals) != len(countForID) {
		t.Fatal("Received individual data sets count doesn't corespond source count")
	}
}

// TestProcStdout ...
func TestProcStdout(t *testing.T) {
	t.SkipNow()
	originals := loadTestData()
	originalForID := makePersonForIDMap(originals)
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)
	reader := bytes.NewReader(originalCombined)

	writer := newTestWriter()

	stdout := processor.NewStdout()
	err := stdout.Connect(reader)
	if err != nil {
		t.Fatal(err)
	}
	stdout.Start()

	totalCount := 0
	countForID := make(map[string]int, 0)
	goal := len(originals)

	for raw := range writer.lines {
		b := processor.ClearEventEnd(raw)
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
}
