package test

import "testing"

func checkTestSet(t *testing.T, originalForID map[string]TestPerson, b []byte) *TestPerson {
	parsed, err := unmarshallPerson(b)
	if err != nil {
		t.Fatalf("Couldn't unmarshall person from: \"%v\" -> %v\n", string(b), err.Error())
	}
	original := originalForID[parsed.ID]
	if !original.Equals(parsed) {
		t.Fatal("Parsed data doesn't conform to original")
	}
	return parsed
}
