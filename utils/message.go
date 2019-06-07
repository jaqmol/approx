package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

// Message ...
type Message struct {
	header              Header
	headerData          []byte
	headerDataReadIndex int64
	headerDataLen       int64
	bodyReadCloser      io.ReadCloser
}

// MessageFromRequest ...
func MessageFromRequest(id int, r *http.Request) (msg *Message, err error) {
	msg = &Message{
		header: Header{
			JSONRPC: "2.0",
			ID:      id,
			Method:  r.Method,
			Params: Params{
				URL:           r.URL.String(),
				Header:        r.Header,
				ContentLength: r.ContentLength,
			},
		},
		bodyReadCloser: r.Body,
	}
	err = msg.initHeaderData()
	return
}

// Header ...
func (m *Message) Header() *Header {
	return &m.header
}

func (m *Message) initHeaderData() (err error) {
	var data []byte
	data, err = json.Marshal(m.header)
	data = append(data, []byte("\n")...)
	if err != nil {
		return
	}
	m.headerData = data
	m.headerDataLen = int64(len(m.headerData))
	return
}

// Reader implementation ...
func (m *Message) Read(buf []byte) (n int, err error) {
	if m.headerDataReadIndex >= m.headerDataLen {
		return m.bodyReadCloser.Read(buf)
	}
	return m.readHeader(buf)
}

func (m *Message) readHeader(buf []byte) (n int, err error) {
	n = copy(buf, m.headerData[m.headerDataReadIndex:])
	m.headerDataReadIndex += int64(n)
	return
}

// Header ...
type Header struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
}

// Params ...
type Params struct {
	URL           string              `json:"url"`
	Header        map[string][]string `json:"header"`
	ContentLength int64               `json:"contentLength,omitempty"`
}
