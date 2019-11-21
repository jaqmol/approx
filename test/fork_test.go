package test

import (
	"bytes"
	"testing"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/processor"
)

// TestFork ...
func TestFork(t *testing.T) {
	// t.SkipNow()
	nextProcsCount := 5
	originals := LoadTestData()
	originalForID := MakePersonForIDMap(originals)
	originalBytes := MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)
	reader := bytes.NewReader(originalCombined)

	conf := configuration.Fork{
		Ident: "test-fork",
		Count: nextProcsCount,
	}
	fork, err := processor.NewFork(&conf)
	if err != nil {
		t.Fatal(err)
	}
	err = fork.Connect(reader)
	if err != nil {
		t.Fatal(err)
	}

	collector, err := processor.NewCollector(fork.Outs()...)
	if err != nil {
		t.Fatal(err)
	}
	collector.Start()
	// serialize := make(chan []byte)
	// for _, r := range fork.Outs() {
	// 	go readFromReader(serialize, r)
	// }
	fork.Start()

	totalCount := 0
	countForID := make(map[string]int, 0)
	goal := nextProcsCount * len(originals)

	for b := range collector.Events() {
		parsed := CheckTestSet(t, originalForID, b)
		totalCount++
		countForID[parsed.ID]++
		if totalCount == goal {
			break
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
