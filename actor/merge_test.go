package actor

import (
	"bytes"
	"testing"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/test"
)

func TestSimpleMerge(t *testing.T) {
	originals := test.LoadTestData()
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)

	producerAlpha := NewProducer(10)
	producerBeta := NewProducer(10)
	merge := NewMerge(10, "merge", 2)
	collector := NewCollector(10)
	receiver := make(chan []byte, 10)

	producerAlpha.Next(merge)
	producerBeta.Next(merge)
	merge.Next(collector)

	counter := 0
	expectedLen := len(originals) * 2

	startCollectingTestMessages(t, collector, receiver, func() {
		close(receiver)
	})
	merge.Start()
	startProducingTestMessages(t, producerAlpha, originalCombined)
	startProducingTestMessages(t, producerBeta, originalCombined)

	for message := range receiver {
		test.CheckTestSet(t, originalForID, message)
		counter++
	}

	if counter != expectedLen {
		t.Fatalf("%v messages expected to be produced, but got %v", expectedLen, counter)
	}
}

func TestMultipleMerge(t *testing.T) {
	originals := test.LoadTestData()
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.EvntEndBytes)
	originalCombined = append(originalCombined, configuration.EvntEndBytes...)

	mergeAlpha := NewMerge(10, "merge-alpha", 2)
	mergeBeta := NewMerge(10, "merge-beta", 3)
	mergeGamma := NewMerge(10, "merge-gamma", 2)
	producerAlpha := NewProducer(10)
	producerBeta := NewProducer(10)
	producerGamma := NewProducer(10)
	producerDelta := NewProducer(10)
	producerEpsilon := NewProducer(10)
	collector := NewCollector(10)
	receiver := make(chan []byte, 10)

	producerAlpha.Next(mergeAlpha)
	producerBeta.Next(mergeAlpha)
	producerGamma.Next(mergeBeta)
	producerDelta.Next(mergeBeta)
	producerEpsilon.Next(mergeBeta)
	mergeAlpha.Next(mergeGamma)
	mergeBeta.Next(mergeGamma)
	mergeGamma.Next(collector)

	startCollectingTestMessages(t, collector, receiver, func() {
		close(receiver)
	})
	mergeAlpha.Start()
	mergeBeta.Start()
	mergeGamma.Start()
	startProducingTestMessages(t, producerAlpha, originalCombined)
	startProducingTestMessages(t, producerBeta, originalCombined)
	startProducingTestMessages(t, producerGamma, originalCombined)
	startProducingTestMessages(t, producerDelta, originalCombined)
	startProducingTestMessages(t, producerEpsilon, originalCombined)

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
