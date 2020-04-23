package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"io"
)

// MsgDelimiter ...
var MsgDelimiter = []byte{'\n', '-', '-', '-', '\n'}

// MsgScanner ...
type MsgScanner struct {
	reader           *bufio.Reader
	message          []byte
	decodedMessage   []byte
	delimitedMessage []byte
	rest             []byte
	err              error
}

// NewMsgScanner ...
func NewMsgScanner(reader io.Reader) *MsgScanner {
	return &MsgScanner{
		reader: bufio.NewReader(reader),
	}
}

// Scan ...
func (m *MsgScanner) Scan() bool {
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
		idx := bytes.Index(buffer, MsgDelimiter)
		if idx > -1 {
			m.message = buffer[:idx]
			m.rest = buffer[idx+len(MsgDelimiter):]
			return true
		}
	}
}

// Message ...
func (m *MsgScanner) Message() []byte {
	return m.message
}

// DecodedMessage ...
func (m *MsgScanner) DecodedMessage() ([]byte, error) {
	if m.decodedMessage == nil {
		m.decodedMessage = make([]byte, base64.StdEncoding.DecodedLen(len(m.message)))
		_, err := base64.StdEncoding.Decode(m.decodedMessage, m.message)
		if err != nil {
			return nil, err
		}
	}
	return m.decodedMessage, nil
}

// DelimitedMessage ...
func (m *MsgScanner) DelimitedMessage() []byte {
	if m.delimitedMessage == nil {
		m.delimitedMessage = append(m.message, MsgDelimiter...)
	}
	return m.delimitedMessage
}

// Err ...
func (m *MsgScanner) Err() error {
	return m.err
}
