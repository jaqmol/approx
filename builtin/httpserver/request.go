package httpserver

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/jaqmol/approx/message"
)

func (h *HTTPServer) startReceiving(port int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		payloadRawMsg := payload(r)
		msg := message.Message{
			ID:      createID(),
			Role:    "request",
			Payload: payloadRawMsg,
		}
		msgBytes, err := json.Marshal(msg)
		catch(err)

		h.cacheResponseWriter(msg.ID, w)
		h.dispatchLine(msgBytes)
	})
	addr := fmt.Sprintf(":%v", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (h *HTTPServer) dispatchLine(lineBytes []byte) {
	fmt.Fprintf(h.stdout, "%v\n", string(lineBytes))
}

func (h *HTTPServer) cacheResponseWriter(id string, w http.ResponseWriter) {
	h.cache.Set(id, w)
}

func createID() string {
	u, err := uuid.NewRandom()
	catch(err)
	return u.String()
}

func payload(r *http.Request) *json.RawMessage {
	body := readBody(r)
	payload := map[string]interface{}{
		"method":  r.Method,
		"url":     urlMap(r),
		"headers": headerToMap(r.Header),
		"body":    *body,
	}
	bytes, err := json.Marshal(payload)
	catch(err)
	rawMsg := json.RawMessage(bytes)
	return &rawMsg
}

func strToInt(str string) int {
	if len(str) == 0 {
		return -1
	}
	i, err := strconv.Atoi(str)
	catch(err)
	return i
}

func urlMap(r *http.Request) map[string]interface{} {
	return map[string]interface{}{
		"host":  r.Host,
		"path":  r.URL.Path,
		"query": urlValuesToMap(r.URL.Query()),
	}
}

func urlValuesToMap(values url.Values) map[string][]string {
	acc := make(map[string][]string)
	for k, v := range values {
		acc[k] = v
	}
	return acc
}

func headerToMap(header http.Header) map[string][]string {
	acc := make(map[string][]string)
	for k, v := range header {
		acc[k] = v
	}
	return acc
}

func readBody(r *http.Request) *string {
	bytes, err := ioutil.ReadAll(r.Body)
	catch(err)
	str := base64.StdEncoding.EncodeToString(bytes)
	return &str
}
