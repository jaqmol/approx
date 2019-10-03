package testpackage

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"testing"

	"github.com/jaqmol/approx/logger"
)

// TestLogger ...
func TestLogger(t *testing.T) {
	data := readTestData()
	readers := readersFromData(data, 5)

	writer := logWriter{collector: data}
	l := logger.NewLogger(&writer)

	for _, r := range readers {
		l.Add(&r)
	}

	l.Start()

	log.Println(len(readers))
}

func readersFromData(data []TestPerson, readersCount int) []logReader {
	acc := make([]logReader, readersCount)
	multi := len(data) / readersCount

	for i := 0; i < readersCount; i++ {
		acc[i] = logReader{testData: make([]TestPerson, 0)}
	}

	for i, person := range data {
		rdrIdx := i / multi
		rdr := acc[rdrIdx]
		rdr.testData = append(rdr.testData, person)
		acc[rdrIdx] = rdr
	}

	return acc
}

func readTestData() []TestPerson {
	dat, err := ioutil.ReadFile("./test-data.json")
	var result []TestPerson
	err = json.Unmarshal(dat, &result)
	check(err)
	return result
}

type logReader struct {
	testData []TestPerson
	index    int
}

func (l *logReader) Read(b []byte) (int, error) {
	if l.index == len(l.testData) {
		return 0, io.EOF
	}
	person := l.testData[l.index]
	bytes, err := json.Marshal(person)
	check(err)
	l.index++
	copy(b, bytes)
	return len(b), nil
}

type logWriter struct {
	collector []TestPerson
}

func (w *logWriter) Write(b []byte) (int, error) {
	var p TestPerson
	err := json.Unmarshal(b, &p)
	check(err)
	w.collector = append(w.collector, p)
	return len(b), nil
}

// TestPerson ...
type TestPerson struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
