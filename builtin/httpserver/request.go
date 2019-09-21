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
	"time"

	"github.com/google/uuid"
	"github.com/jaqmol/approx/message"
)

func (h *HTTPServer) startReceiving(port int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		respChan := make(chan *message.Message)
		dd := dispatchData{
			request: message.Message{
				// ID:        createID(),
				ID:        r.URL.Path,
				Role:      "request",
				IsEnd:     true,
				MediaType: "application/json",
				Encoding:  "utf8",
				Data:      makeRequestMessageData(r),
			},
			respChan: respChan,
		}
		h.dispatchChannel <- &dd
		h.waitForResponse(respChan, w, dd.request.ID)
		// expectingResponses := true
		// for expectingResponses {
		// 	select {
		// 	case response := <-respChan:
		// 		h.respond(w, response)
		// 		expectingResponses = !response.IsEnd
		// 	case <-time.After(h.timeout):
		// 		h.respondWithPipelineResponseTimeout(w, dd.request.ID)
		// 		expectingResponses = false
		// 	}
		// }
	})
	addr := fmt.Sprintf(":%v", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (h *HTTPServer) waitForResponse(respChan <-chan *message.Message, w http.ResponseWriter, id string) {
	waiting := true
	for waiting {
		select {
		case response := <-respChan:
			h.respond(w, response)
			waiting = !response.IsEnd
		case <-time.After(h.timeout):
			h.respondWithPipelineResponseTimeout(w, id)
			waiting = false
		}
	}
}

func (h *HTTPServer) startDispatching() {
	for dd := range h.dispatchChannel {
		h.cacheResponseChannel(dd.request.ID, dd.respChan)
		env := dd.request.Envelope()
		h.stdout.Write() <- env.Bytes
	}
}

type dispatchData struct {
	request  message.Message
	respChan chan<- *message.Message
}

func (h *HTTPServer) cacheResponseChannel(id string, rc chan<- *message.Message) {
	h.cache.Set(id, rc)
}

func createID() string {
	u, err := uuid.NewRandom()
	catch(err)
	return u.String()
}

type requestMessageData struct {
	Method   string              `json:"method"`
	URL      requestURL          `json:"url"`
	Headers  map[string][]string `json:"headers"`
	Encoding string              `json:"encoding"`
	Body     string              `json:"body"`
}

func makeRequestMessageData(r *http.Request) (bytes []byte) {
	bodyB64 := readBodyBase64(r)
	url := makeRequestURL(r)
	rp := requestMessageData{
		Method:   r.Method,
		URL:      *url,
		Headers:  headerToMap(r.Header),
		Encoding: "base64",
		Body:     *bodyB64,
	}
	bytes, err := json.Marshal(rp)
	catch(err)
	return
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

func readBodyBase64(r *http.Request) *string {
	bytes, err := ioutil.ReadAll(r.Body)
	catch(err)
	str := base64.StdEncoding.EncodeToString(bytes)
	return &str
}
