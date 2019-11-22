package actor

import (
	"bytes"
	"io"
	"testing"

	"github.com/jaqmol/approx/config"
	"github.com/jaqmol/approx/event"
	"github.com/jaqmol/approx/test"
)

func TestSingleProducer(t *testing.T) {
	originals := test.LoadTestData()
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, config.EvntEndBytes)
	originalCombined = append(originalCombined, config.EvntEndBytes...)

	producer := NewProducer(10)
	collector := NewCollector(10)
	receiver := make(chan []byte, 10)

	producer.Next(collector)

	startCollectingTestMessages(t, collector, receiver, func() {
		close(receiver)
	})
	startProducingTestMessages(t, producer, originalCombined)

	counter := 0

	for message := range receiver {
		test.CheckTestSet(t, originalForID, message)
		counter++
	}

	if counter != len(originals) {
		t.Fatalf("%v messages expected to be produced, but got %v", len(originals), counter)
	}
}

func startProducingTestMessages(t *testing.T, producer *Producer, sourceData []byte) {
	go func() {
		reader := bytes.NewReader(sourceData)
		scanner := event.NewScanner(reader)
		err := producer.Produce(func() ([]byte, error) {
			if scanner.Scan() {
				raw := scanner.Bytes()
				return event.ScannedBytesCopy(raw), nil
			}
			return nil, io.EOF
		})
		if err != nil {
			t.Fatal(err)
		}
	}()
}
