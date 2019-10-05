package processor

import (
	"github.com/jaqmol/approx/configuration"
)

var msgEndLength int

func init() {
	msgEndLength = len(configuration.MsgEndBytes)
}

func msgEndedCopy(data []byte) []byte {
	dataCopy := make([]byte, len(data)+msgEndLength)
	copy(dataCopy, data)
	return append(dataCopy, configuration.MsgEndBytes...)
}
