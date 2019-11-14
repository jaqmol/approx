package test

import (
	"bytes"
	"io"
	"testing"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/event"
	"github.com/jaqmol/approx/processor"
)

// TestMerge ...
func TestMerge(t *testing.T) {
	// t.SkipNow()
	prevProcsCount := 5
	originals := LoadTestData()
	originalForID := MakePersonForIDMap(originals)
	originalBytes := MarshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)

	readers := make([]io.Reader, prevProcsCount)
	for i := range readers {
		readers[i] = bytes.NewReader(originalCombined)
	}

	conf := configuration.Merge{
		Ident: "test-merge",
	}
	merge, err := processor.NewMerge(&conf)
	if err != nil {
		t.Fatal(err)
	}
	err = merge.Connect(readers...)
	if err != nil {
		t.Fatal(err)
	}

	totalCount := 0
	countForID := make(map[string]int, 0)
	goal := prevProcsCount * len(originals)

	outputReader := merge.Out()
	merge.Start()
	scanner := event.NewScanner(outputReader)

	for scanner.Scan() {
		raw := scanner.Bytes()
		data := bytes.Trim(raw, "\x00")
		parsed := CheckTestSet(t, originalForID, data)
		totalCount++
		countForID[parsed.ID]++
	}

	if goal != totalCount {
		t.Fatal("Merged count doesn't corespond to multitude of source count")
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
