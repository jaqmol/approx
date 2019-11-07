package processor

import (
	"github.com/jaqmol/approx/configuration"
)

var evntEndLength int

func init() {
	evntEndLength = len(configuration.EvntEndBytes)
}

func evntEndedCopy(data []byte) []byte {
	dataCopy := make([]byte, len(data)+evntEndLength)
	copy(dataCopy, data)
	return append(dataCopy, configuration.EvntEndBytes...)
}
