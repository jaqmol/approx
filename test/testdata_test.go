package testpackage

import (
	"testing"
)

// TestTestData ...
func TestTestData(t *testing.T) {
	original := loadTestData()
	originalForID := makePersonForIDMap(original)
	originalBytes := marshallPeople(original)
	parsed := unmarshallPeople(originalBytes)
	parsedForID := makePersonForIDMap(parsed)
	for id, person := range originalForID {
		readPerson, ok := parsedForID[id]
		if !ok {
			t.FailNow()
		}
		if !person.Equals(&readPerson) {
			t.FailNow()
		}
	}
}
