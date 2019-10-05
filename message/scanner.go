package message

import (
	"bufio"
	"bytes"
	"io"

	"github.com/jaqmol/approx/configuration"
)

var msgEndLength int

func init() {
	msgEndLength = len(configuration.MsgEndBytes)
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

	msgEndIndex := bytes.Index(data, configuration.MsgEndBytes)

	if msgEndIndex == -1 {
		return 0, nil, nil
	}

	token = data[:msgEndIndex]
	advance = len(token) + msgEndLength
	return
}
