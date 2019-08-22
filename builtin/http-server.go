package builtin

import (
	"fmt"
	"io"
	"net/http"

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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to my website!")
	})
}

// MakeHTTPServer ...
func MakeHTTPServer(def *definition.Definition) *HTTPServer {
	return &HTTPServer{
		def: *def,
	}
}
