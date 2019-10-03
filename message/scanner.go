package message

import (
	"bufio"
	"bytes"
	"io"
	"log"

	"github.com/jaqmol/approx/configuration"
)

var incompleteMsgEndSuffixes [][]byte
var msgEnd []byte
var msgEndLength int

func init() {
	incompleteMsgEndSuffixes = make([][]byte, 0)
	msgEnd = []byte(configuration.MessageEnd)
	msgEndLength = len(msgEnd)
	for i := len(msgEnd) - 1; i > 0; i-- {
		suffix := msgEnd[:i]
		log.Printf("incompleteSeparatorSuffix: %v\n", string(suffix))
		incompleteMsgEndSuffixes = append(incompleteMsgEndSuffixes, suffix)
	}
}

// NewScanner ...
func NewScanner(r io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	scanner.Split(splitFn)
	return scanner
}

func splitFn(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		return 0, nil, io.EOF
	}

	sepIdx := bytes.Index(data, msgEnd)
	if sepIdx > -1 {
		msg := data[:sepIdx]
		advance = len(msg) + msgEndLength
		token = msg
		return
	}

	if containsIncompleteSeparatorSuffix(data) {
		return 0, nil, nil
	}

	advance = len(data)
	token = data
	return
}

func containsIncompleteSeparatorSuffix(data []byte) bool {
	for _, suffix := range incompleteMsgEndSuffixes {
		if bytes.HasSuffix(data, suffix) {
			return true
		}
	}
	return false
}
