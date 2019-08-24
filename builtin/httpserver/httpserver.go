package httpserver

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/ReneKroon/ttlcache"
	"github.com/jaqmol/approx/definition"
)

// uuid.NewRandom()

// HTTPServer ...
type HTTPServer struct {
	def     definition.Definition
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
	running bool
	cache   *ttlcache.Cache
}

// SetStdin ...
func (h *HTTPServer) SetStdin(r io.Reader) {
	h.stdin = r
}

// SetStdout ...
func (h *HTTPServer) SetStdout(w io.Writer) {
	h.stdout = w
}

// SetStderr ...
func (h *HTTPServer) SetStderr(w io.Writer) {
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
	log.Printf("Server running at %v\n", port)
}

// MakeHTTPServer ...
func MakeHTTPServer(def *definition.Definition, upstreamTimeout time.Duration) *HTTPServer {
	if upstreamTimeout <= 0 {
		upstreamTimeout = 10000
	}
	h := &HTTPServer{
		def:   *def,
		cache: ttlcache.NewCache(),
	}
	h.cache.SetTTL(upstreamTimeout)
	return h
}

func catch(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}
