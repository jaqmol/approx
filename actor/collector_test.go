package actor

import (
	"bytes"
	"testing"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/test"
)

func TestSingleCollector(t *testing.T) {
	// Single collector is being tested in producer_test.go
}

func TestMultipleCollectors(t *testing.T) {
	originals := test.LoadTestData()
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)

	producer := NewProducer(10)
	collectorAlpha := NewCollector(10)
	collectorBeta := NewCollector(10)
	receiver := make(chan []byte, 10)

	producer.Next(collectorAlpha, collectorBeta)

	startCollectingTestMessages(t, collectorAlpha, receiver, func() {})
	startCollectingTestMessages(t, collectorBeta, receiver, func() {})
	startProducingTestMessages(t, producer, originalCombined)

	counter := 0
	expectedLen := len(originals) * 2

	for message := range receiver {
		test.CheckTestSet(t, originalForID, message)
		counter++
		if counter == expectedLen {
			close(receiver)
		}
	}

	if counter != expectedLen {
		t.Fatalf("%v messages expected to be produced, but got %v", expectedLen, counter)
	}
}

func startCollectingTestMessages(
	t *testing.T,
	collector *Collector,
	receiver chan<- []byte,
	finished func(),
) {
	go func() {
		err := collector.Collect(func(message []byte) error {
			receiver <- message
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		finished()
	}()
}
