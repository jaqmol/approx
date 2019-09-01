package httpserver

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/jaqmol/approx/message"
)

type responsePayload struct {
	Status      int    `json:"status"`
	ContentType string `json:"contentType"`
	Body        string `json:"body"`
}

func (h *HTTPServer) startResponding() {
	scanner := bufio.NewScanner(h.stdin)
	for scanner.Scan() {
		msgBytes := scanner.Bytes()
		var msg message.Message
		err := json.Unmarshal(msgBytes, &msg)
		if err != nil {
			message.WriteLogEntry(h.stderr, message.Fail, "", err.Error())
		} else {
			rc, ok := h.uncacheResponseChannel(msg.ID)
			if ok {
				rc <- &msg
			} else {
				message.WriteLogEntry(h.stderr, message.Fail, msg.ID, "No response channel found, most likely due to timeout")
			}
		}
	}
}

func (h *HTTPServer) respond(rw http.ResponseWriter, msg *message.Message) bool {
	payloadBytes := msg.Payload
	var payload responsePayload
	err := json.Unmarshal(*payloadBytes, &payload)
	if err != nil {
		return respond500Error(rw, msg, err)
	}
	rw.WriteHeader(payload.Status)
	rw.Header().Set("Content-Type", payload.ContentType)
	body, err := bodyFromPayloadBody(payload.Body)
	if err != nil {
		return respond500Error(rw, msg, err)
	}
	if body != nil {
		_, err = rw.Write(body)
		if err != nil {
			return respond500Error(rw, msg, err)
		}
	}
	return true
}

func respond500Error(rw http.ResponseWriter, msg *message.Message, err error) bool {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusInternalServerError)
	message.WriteLogEntry(rw, message.Fail, msg.ID, err.Error())
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
