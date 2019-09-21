package httpserver

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/jaqmol/approx/message"
)

func (h *HTTPServer) startResponding() {
	messageReader := message.NewReader(h.stdin.Read())
	for resp := range messageReader.Messages() {
		rc, ok := h.uncacheResponseChannel(resp.ID)
		if ok {
			rc <- resp
		} else {
			err := fmt.Errorf("No response channel found for message ID %v", resp.ID)
			panic(err)
		}
	}
}

func (h *HTTPServer) respond(rw http.ResponseWriter, resp *message.Message) bool {
	if resp.Seq == 0 {
		rw.Header().Set("Content-Type", resp.MediaType)
		rw.WriteHeader(resp.Status)
	}

	var body []byte
	if resp.Encoding == "base64" {
		body = make([]byte, base64.StdEncoding.DecodedLen(len(resp.Data)))
		_, err := base64.StdEncoding.Decode(body, resp.Data)
		if err != nil {
			panic(err)
		}
	} else {
		body = resp.Data
	}

	writtenBytesCount := 0
	for writtenBytesCount < len(body) {
		wbc, err := rw.Write(body)
		if err != nil {
			return respond500Error(rw, resp.ID, err)
		}
		writtenBytesCount += wbc
	}

	return true
}

func (h *HTTPServer) respondWithPipelineResponseTimeout(rw http.ResponseWriter, id string) bool {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusInternalServerError)
	_, err := fmt.Fprintf(rw, "Upstream pipeline response timeout for message ID: %v", id)
	if err != nil {
		panic(err)
	}
	return true
}

func respond500Error(rw http.ResponseWriter, id string, err error) bool {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusInternalServerError)
	_, err2 := fmt.Fprintf(rw, "Internal server error for message (ID: %v): %v", id, err.Error())
	if err2 != nil {
		panic(err2)
	}
	return false
}

func (h *HTTPServer) uncacheResponseChannel(id string) (rc chan<- *message.Message, ok bool) {
	rcIf, ok := h.cache.Get(id)
	if ok {
		rc = rcIf.(chan<- *message.Message)
		h.cache.Remove(id)
	}
	return
}

func bodyFromPayloadBody(payloadBody string) ([]byte, error) {
	if len(payloadBody) == 0 {
		return nil, nil
	}
	bytes, err := base64.StdEncoding.DecodeString(payloadBody)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
