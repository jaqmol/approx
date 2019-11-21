package actor

import (
	"bytes"
	"testing"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/test"
)

func TestSimpleFork(t *testing.T) {
	originals := test.LoadTestData()
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)

	fork := NewFork(10, "fork", 2)
	producer := NewProducer(10)
	collectorAlpha := NewCollector(10)
	collectorBeta := NewCollector(10)
	receiver := make(chan []byte, 10)

	producer.Next(fork)
	fork.Next(collectorAlpha, collectorBeta)

	alphaDone, betaDone := false, false
	closeReceiverIfDone := func() {
		if alphaDone && betaDone {
			close(receiver)
		}
	}

	startCollectingTestMessages(t, collectorAlpha, receiver, func() {
		alphaDone = true
		closeReceiverIfDone()
	})
	startCollectingTestMessages(t, collectorBeta, receiver, func() {
		betaDone = true
		closeReceiverIfDone()
	})
	fork.Start()
	startProducingTestMessages(t, producer, originalCombined)

	counter := 0
	expectedLen := len(originals) * 2

	for message := range receiver {
		test.CheckTestSet(t, originalForID, message)
		counter++
	}

	if counter != expectedLen {
		t.Fatalf("%v messages expected to be produced, but got %v", expectedLen, counter)
	}
}

func TestMultipleFork(t *testing.T) {
	originals := test.LoadTestData()
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)

	forkAlpha := NewFork(10, "fork-alpha", 2)
	forkBeta := NewFork(10, "fork-beta", 3)
	forkGamma := NewFork(10, "fork-gamma", 2)

	producer := NewProducer(10)
	collectorAlpha := NewCollector(10)
	collectorBeta := NewCollector(10)
	collectorGamma := NewCollector(10)
	collectorDelta := NewCollector(10)
	collectorEpsilon := NewCollector(10)
	receiver := make(chan []byte, 10)

	producer.Next(forkAlpha)
	forkAlpha.Next(forkBeta, forkGamma)
	forkBeta.Next(collectorAlpha, collectorBeta, collectorGamma)
	forkGamma.Next(collectorDelta, collectorEpsilon)

	alphaDone, betaDone, gammaDone, deltaDone, epsilonDone := false, false, false, false, false
	closeReceiverIfDone := func() {
		if alphaDone && betaDone && gammaDone && deltaDone && epsilonDone {
			close(receiver)
		}
	}

	startCollectingTestMessages(t, collectorAlpha, receiver, func() {
		alphaDone = true
		closeReceiverIfDone()
	})
	startCollectingTestMessages(t, collectorBeta, receiver, func() {
		betaDone = true
		closeReceiverIfDone()
	})
	startCollectingTestMessages(t, collectorGamma, receiver, func() {
		gammaDone = true
		closeReceiverIfDone()
	})
	startCollectingTestMessages(t, collectorDelta, receiver, func() {
		deltaDone = true
		closeReceiverIfDone()
	})
	startCollectingTestMessages(t, collectorEpsilon, receiver, func() {
		epsilonDone = true
		closeReceiverIfDone()
	})
	forkAlpha.Start()
	forkBeta.Start()
	forkGamma.Start()
	startProducingTestMessages(t, producer, originalCombined)

	counter := 0
	expectedLen := len(originals) * 5

	for message := range receiver {
		test.CheckTestSet(t, originalForID, message)
		counter++
	}

	if counter != expectedLen {
		t.Fatalf("%v messages expected to be produced, but got %v", expectedLen, counter)
	}
}
