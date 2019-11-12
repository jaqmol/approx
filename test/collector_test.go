package test

import (
	"bytes"
	"testing"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/processor"
)

// TestSingleCollector ...
func TestSingleCollector(t *testing.T) {
	originals := loadTestData() // [:10]
	originalForID := makePersonForIDMap(originals)
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)

	reader := bytes.NewReader(originalCombined)

	collector, err := processor.NewCollector(reader)
	if err != nil {
		t.Fatal(err)
	}
	collector.Start()

	goal := len(originals)
	counter := 0

	for b := range collector.Events() {
		checkTestSet(t, originalForID, b)
		counter++
		if counter == goal {
			break
		}
	}
}

// TestMultipleCollectors ...
func TestMultipleCollectors(t *testing.T) {
	originals := loadTestData()[:10]
	originalForID := makePersonForIDMap(originals)
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)

	alphaReader := bytes.NewReader(originalCombined)
	betaReader := bytes.NewReader(originalCombined)

	collector, err := processor.NewCollector(alphaReader, betaReader)
	if err != nil {
		t.Fatal(err)
	}
	collector.Start()

	goal := len(originals)
	counter := 0

	for b := range collector.Events() {
		checkTestSet(t, originalForID, b)
		counter++
		if counter == goal {
			break
		}
	}
}
