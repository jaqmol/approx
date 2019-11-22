package test

import (
	"bytes"
	"testing"

	"github.com/jaqmol/approx/config"
	"github.com/jaqmol/approx/event"
)

// TestScanner ...
func TestScanner(t *testing.T) {
	// t.SkipNow()
	originals := LoadTestData() // [:10]
	originalForID := MakePersonForIDMap(originals)
	originalBytes := MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, config.EvntEndBytes)
	originalCombined = append(originalCombined, config.EvntEndBytes...)
	reader := bytes.NewReader(originalCombined)

	scanner := event.NewScanner(reader)
	count := 0

	for scanner.Scan() {
		b := scanner.Bytes()
		CheckTestSet(t, originalForID, b)
		count++
	}

	if len(originals) != count {
		t.FailNow()
	}
}
