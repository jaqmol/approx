package run

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/jaqmol/approx/utils"

	"github.com/jaqmol/approx/conf"
	"github.com/jaqmol/approx/flow"
)

// NewHTTPServerProc ...
func NewHTTPServerProc(conf *conf.HTTPServerConf) (*HTTPServerProc, error) {
	return &HTTPServerProc{
		conf: conf,
		ins:  make(map[string]proc.Conn),
		outs: make(map[string]proc.Conn),
	}, nil
}

// HTTPServerProc ...
type HTTPServerProc struct {
	conf *conf.HTTPServerConf
	ins  map[string]proc.Conn
	outs map[string]proc.Conn
}

// Type ...
func (hp *HTTPServerProc) Type() conf.Type {
	return hp.conf.Type()
}

// Conf ...
func (hp *HTTPServerProc) Conf() conf.Conf {
	return hp.conf
}

// In ...
func (hp *HTTPServerProc) In(name string) proc.Conn {
	return hp.ins[name]
}

// Out ...
func (hp *HTTPServerProc) Out(name string) proc.Conn {
	return hp.outs[name]
}

// AddIn ...
func (hp *HTTPServerProc) AddIn(name string, c proc.Conn) {
	hp.ins[name] = c
}

// AddOut ...
func (hp *HTTPServerProc) AddOut(name string, c proc.Conn) {
	hp.outs[name] = c
}

// Start ...
func (hp *HTTPServerProc) Start(errChan chan<- error) {
	go func() {
		idCounter := 0
		responseWriterForID := make(map[int]http.ResponseWriter)

		log.Printf("Starting HTTP server for \"%v\" @%v", hp.conf.Endpoint(), hp.conf.Port())

		http.HandleFunc(hp.conf.Endpoint(), func(w http.ResponseWriter, r *http.Request) {
			idCounter++
			id := idCounter
			responseWriterForID[id] = w

			msg, err := utils.MessageFromRequest(id, r)
			if err != nil {
				errChan <- err
			}

			io.Copy(w, msg)
			// fmt.Fprintf(w, "HTTP server request: %v\n", r)
		})

		// Static file server
		// fs := http.FileServer(http.Dir("static/"))
		// http.Handle("/static/", http.StripPrefix("/static/", fs))

		p := fmt.Sprintf(":%v", hp.conf.Port())
		err := http.ListenAndServe(p, nil)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}()
}
