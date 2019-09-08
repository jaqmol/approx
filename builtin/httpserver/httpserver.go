package httpserver

import (
	"log"
	"os"
	"time"

	"github.com/ReneKroon/ttlcache"
	"github.com/jaqmol/approx/definition"
	"github.com/jaqmol/approx/pipe"
)

// uuid.NewRandom()

// HTTPServer ...
type HTTPServer struct {
	def             definition.Definition
	stdin           *pipe.Reader
	stdout          *pipe.Writer
	stderr          *pipe.Writer
	running         bool
	timeout         time.Duration
	cache           *ttlcache.Cache
	dispatchChannel chan *dispatchData
}

// SetStdin ...
func (h *HTTPServer) SetStdin(r *pipe.Reader) {
	h.stdin = r
}

// SetStdout ...
func (h *HTTPServer) SetStdout(w *pipe.Writer) {
	h.stdout = w
}

// SetStderr ...
func (h *HTTPServer) SetStderr(w *pipe.Writer) {
	h.stderr = w
}

// Definition ...
func (h *HTTPServer) Definition() *definition.Definition {
	return &h.def
}

// Start ...
func (h *HTTPServer) Start() {
	portStr := *h.def.Env["PORT"]
	port := strToInt(portStr)
	go h.startReceiving(port)
	go h.startResponding()
	go h.startDispatching()
	log.Printf("Server running at %v\n", port)
}

func parseTimeout(def *definition.Definition) (timeout time.Duration) {
	timeoutStr, ok := def.Env["TIMEOUT"]
	if ok && len(*timeoutStr) > 0 {
		var err error
		timeout, err = time.ParseDuration(*timeoutStr)
		if err != nil {
			catch(err)
		}
	} else {
		timeout, _ = time.ParseDuration("10s")
	}
	return
}

// MakeHTTPServer ...
func MakeHTTPServer(def *definition.Definition) *HTTPServer {
	h := &HTTPServer{
		def:             *def,
		timeout:         parseTimeout(def),
		cache:           ttlcache.NewCache(),
		dispatchChannel: make(chan *dispatchData),
	}
	h.cache.SetTTL(h.timeout)
	return h
}

func catch(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}
