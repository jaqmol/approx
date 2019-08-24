package httpserver

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jaqmol/approx/message"
)

func (h *HTTPServer) startResponding() {
	scanner := bufio.NewScanner(h.stdin)
	for scanner.Scan() {
		msgBytes := scanner.Bytes()
		var msg message.Message
		err := json.Unmarshal(msgBytes, &msg)
		if err != nil {
			message.WriteError(h.stderr, "", err.Error())
		} else {
			h.respond(&msg)
		}
	}
}

func (h *HTTPServer) respond(msg *message.Message) {
	// {}{
	// 	"method":  r.Method,
	// 	"url":     urlMap(r),
	// 	"headers": headerToMap(r.Header),
	// 	"body":    *body,
	// }
	rw, ok := h.uncacheResponseWriter(msg.ID)
	if ok {
		payloadBytes := msg.Payload
		var payload map[string]interface{}
		err := json.Unmarshal(payloadBytes, &payload)
		if err != nil {
			message.WriteError(rw, msg.ID, err.Error())
			rw.WriteHeader(500)
			return
		}
		status, err := statusFromResponsePayload(payload)
		if err != nil {
			message.WriteError(rw, msg.ID, err.Error())
			rw.WriteHeader(500)
			return
		}
		contentType, err := contentTypeFromResponsePayload(payload)
		if err != nil {
			message.WriteError(rw, msg.ID, err.Error())
			rw.WriteHeader(500)
			return
		}
		rw.Header().Set("Content-Type", contentType)
		body, err := bodyFromResponsePayload(payload)
		if err != nil {
			message.WriteError(rw, msg.ID, err.Error())
			rw.WriteHeader(500)
			return
		}
		_, err = rw.Write(body)
		if err != nil {
			message.WriteError(rw, msg.ID, err.Error())
			rw.WriteHeader(500)
			return
		}
		rw.WriteHeader(status)
	} else {
		message.WriteError(h.stderr, msg.ID, "Response timeout: message too late for response")
	}
}

func (h *HTTPServer) uncacheResponseWriter(id string) (w http.ResponseWriter, ok bool) {
	wInfh, ok := h.cache.Get(id)
	if ok {
		w = wInfh.(http.ResponseWriter)
	}
	return
}

func statusFromResponsePayload(payload map[string]interface{}) (int, error) {
	ifStatus, ok := payload["status"]
	if !ok {
		return -1, fmt.Errorf("Status code missing in response message")
	}
	status, ok := ifStatus.(int)
	if !ok {
		return -1, fmt.Errorf("Status code has wrong type in response message")
	}
	return status, nil
}

func contentTypeFromResponsePayload(payload map[string]interface{}) (string, error) {
	ifContType, ok := payload["content-type"]
	if !ok {
		return "", fmt.Errorf("Content type missing in response message")
	}
	contType, ok := ifContType.(string)
	if !ok {
		return "", fmt.Errorf("Content type has wrong type in response message")
	}
	return contType, nil
}

func bodyFromResponsePayload(payload map[string]interface{}) ([]byte, error) {
	ifBodyStr, ok := payload["body"]
	if !ok {
		return nil, fmt.Errorf("Body missing in response message")
	}
	bodyStr, ok := ifBodyStr.(string)
	if !ok {
		return nil, fmt.Errorf("Body has wrong type in response message")
	}
	bytes, err := base64.StdEncoding.DecodeString(bodyStr)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
