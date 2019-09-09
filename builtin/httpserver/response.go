package httpserver

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/jaqmol/approx/message"
)

func (h *HTTPServer) startResponding() {
	scanner := bufio.NewScanner(h.stdin)
	for scanner.Scan() {
		msgBytes := scanner.Bytes()
		resp := message.ParseMessage(msgBytes)
		if resp == nil {
			errBytes := []byte(fmt.Sprintf("No message in: %v", string(msgBytes)))
			h.stderr.Channel() <- errBytes
			continue
		}
		rc, ok := h.uncacheResponseChannel(resp.ID)
		if ok {
			rc <- resp
		} else {
			errBytes := []byte(fmt.Sprintf("No response channel found for message ID %v, timeout likely\n", resp.ID))
			h.stderr.Channel() <- errBytes
		}
	}
}

func (h *HTTPServer) respond(rw http.ResponseWriter, resp *message.Message) bool {
	rw.Header().Set("Content-Type", resp.MediaType)
	rw.WriteHeader(resp.Status)

	_, err := rw.Write(resp.Body)
	if err != nil {
		return respond500Error(rw, resp.ID, err)
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
