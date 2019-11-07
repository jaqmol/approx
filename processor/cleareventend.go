package processor

import (
	"bytes"

	"github.com/jaqmol/approx/configuration"
)

// ClearEventEnd ...
func ClearEventEnd(raw []byte) []byte {
	msg := bytes.ReplaceAll(raw, []byte("\x00"), []byte(""))
	msgEndIndex := bytes.Index(msg, configuration.EvntEndBytes)
	if msgEndIndex == -1 {
		return msg
	}
	return msg[:msgEndIndex]
}
