package testpackage

import (
	"bytes"
	"testing"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/processor"
)

// TestMerge ...
func TestMerge(t *testing.T) {
	// TODO: implement
	prevProcsCount := 5
	originals := loadTestData()[:10]
	originalForID := makePersonForIDMap(originals)
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.MsgEndBytes)
	originalCombined = append(originalCombined, configuration.MsgEndBytes...)
	reader := bytes.NewReader(originalCombined)

	conf := configuration.Fork{
		Ident:     "test-fork",
		NextProcs: makeTestProcs(prevProcsCount),
	}
	fork := processor.NewFork(&conf, reader)

	serialize := make(chan []byte)
	for _, r := range fork.Outs() {
		go readFromReader(serialize, r)
	}
	fork.Start()

	totalCount := 0
	countForID := make(map[string]int, 0)
	goal := prevProcsCount * len(originals)

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
		if count != prevProcsCount {
			t.Fatalf("Expected to receive %v data sets, but got %v\n", prevProcsCount, count)
		}
	}
}
