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
func MakeHTTPServer(def *definition.Definition) *HTTPServer {
	h := &HTTPServer{
		def:   *def,
		cache: ttlcache.NewCache(),
	}
	var timeout time.Duration
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
	h.cache.SetTTL(timeout)
	return h
}

func catch(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}
