package test

import (
	"bytes"
	"io"
	"testing"

	"github.com/jaqmol/approx/event"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/processor"
)

// TestFork ...
func TestFork(t *testing.T) {
	nextProcsCount := 5
	originals := loadTestData()[:5]
	originalForID := makePersonForIDMap(originals)
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.MsgEndBytes)
	originalCombined = append(originalCombined, configuration.MsgEndBytes...)
	reader := bytes.NewReader(originalCombined)

	conf := configuration.Fork{
		Ident:     "test-fork",
		NextProcs: makeTestProcs(nextProcsCount),
	}
	fork := processor.NewFork(&conf, reader)

	serialize := make(chan []byte)
	for _, r := range fork.Outs() {
		go readFromReader(serialize, r)
	}
	fork.Start()

	totalCount := 0
	countForID := make(map[string]int, 0)
	goal := nextProcsCount * len(originals)

	for b := range serialize {
		parsed := checkTestSet(t, originalForID, b)
		totalCount++
		countForID[parsed.ID]++
		if totalCount == goal {
			close(serialize)
		}
	}

	if goal != totalCount {
		t.Fatal("Forked count doesn't corespond to multitude of source count")
	}

	if len(originals) != len(countForID) {
		t.Fatal("Received individual data sets count doesn't corespond source count")
	}

	for _, count := range countForID {
		if count != nextProcsCount {
			t.Fatalf("Expected to receive %v data sets, but got %v\n", nextProcsCount, count)
		}
	}
}

func readFromReader(serialize chan<- []byte, reader io.Reader) {
	scanner := event.NewScanner(reader)
	for scanner.Scan() {
		raw := scanner.Bytes()
		original := bytes.Trim(raw, "\x00")
		toPassOn := make([]byte, len(original))
		copy(toPassOn, original)
		serialize <- toPassOn
	}
}
