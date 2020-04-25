package main

import (
	"bufio"
	"bytes"
	"io"
)

// // MsgDelimiter ...
// var MsgDelimiter = []byte{'\n', '-', '-', '-', '\n'}

// HeavyDutyScanner ...
type HeavyDutyScanner struct {
	reader           *bufio.Reader
	delimiter        []byte
	message          []byte
	decodedMessage   []byte
	delimitedMessage []byte
	rest             []byte
	err              error
}

// NewHeavyDutyScanner ...
func NewHeavyDutyScanner(reader io.Reader, delimiter []byte) *HeavyDutyScanner {
	return &HeavyDutyScanner{
		reader: bufio.NewReader(reader),
	}
}

// Scan ...
func (m *HeavyDutyScanner) Scan() bool {
	m.message = nil
	m.decodedMessage = nil
	m.delimitedMessage = nil

	var buffer []byte
	if len(m.rest) > 0 {
		buffer = append(buffer, m.rest...)
		m.rest = nil
	}

	for {
		line, err := m.reader.ReadBytes('\n')
		if err != nil {
			m.err = err
			return false
		}
		buffer = append(buffer, line...)
		idx := bytes.Index(buffer, m.delimiter)
		if idx > -1 {
			m.message = buffer[:idx]
			m.rest = buffer[idx+len(m.delimiter):]
			return true
		}
	}
}

// Message ...
func (m *HeavyDutyScanner) Message() []byte {
	return m.message
}

// DecodedMessage ...
func (m *HeavyDutyScanner) DecodedMessage(
	decodedLen func(n int) int, // TODO: TEST
	decode func(dst, src []byte) (n int, err error), // TODO: TEST
) ([]byte, error) {
	if m.decodedMessage == nil {
		// m.decodedMessage = make([]byte, base64.StdEncoding.DecodedLen(len(m.message)))
		// _, err := base64.StdEncoding.Decode(m.decodedMessage, m.message)
		m.decodedMessage = make([]byte, decodedLen(len(m.message)))
		_, err := decode(m.decodedMessage, m.message)
		if err != nil {
			return nil, err
		}
	}
	return m.decodedMessage, nil
}

// DelimitedMessage ...
func (m *HeavyDutyScanner) DelimitedMessage() []byte {
	if m.delimitedMessage == nil {
		m.delimitedMessage = append(m.message, m.delimiter...)
	}
	return m.delimitedMessage
}

// Err ...
func (m *HeavyDutyScanner) Err() error {
	return m.err
}
