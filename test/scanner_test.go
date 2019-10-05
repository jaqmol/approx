package testpackage

import (
	"bytes"
	"testing"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/message"
)

// TestScanner ...
func TestScanner(t *testing.T) {
	originals := loadTestData() // [:10]
	originalForID := makePersonForIDMap(originals)
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.MsgEndBytes)
	originalCombined = append(originalCombined, configuration.MsgEndBytes...)
	reader := bytes.NewReader(originalCombined)

	scanner := message.NewScanner(reader)
	count := 0

	for scanner.Scan() {
		b := scanner.Bytes()
		checkTestSet(t, originalForID, b)
		count++
	}

	if len(originals) != count {
		t.FailNow()
	}
}
