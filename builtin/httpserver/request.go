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

type requestPayload struct {
	Method  string              `json:"method"`
	URL     requestURL          `json:"url"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

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

		rc := make(chan *message.Message)
		h.cacheResponseChannel(msg.ID, rc)
		h.dispatchLine(msgBytes)

		response := <-rc
		h.respond(w, response)
	})
	addr := fmt.Sprintf(":%v", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (h *HTTPServer) dispatchLine(bytes []byte) {
	bytes = append(bytes, []byte("\n")...)
	_, err := h.stdout.Write(bytes)
	catch(err)
}

func (h *HTTPServer) cacheResponseChannel(id string, rc chan<- *message.Message) {
	h.cache.Set(id, rc)
}

func createID() string {
	u, err := uuid.NewRandom()
	catch(err)
	return u.String()
}

func payload(r *http.Request) *json.RawMessage {
	body := readBody(r)
	url := makeRequestURL(r)
	payload := requestPayload{
		Method:  r.Method,
		URL:     *url,
		Headers: headerToMap(r.Header),
		Body:    *body,
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

type requestURL struct {
	Host  string              `json:"host"`
	Path  string              `json:"path"`
	Query map[string][]string `json:"query"`
}

func makeRequestURL(r *http.Request) *requestURL {
	return &requestURL{
		Host:  r.Host,
		Path:  r.URL.Path,
		Query: urlValuesToMap(r.URL.Query()),
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
