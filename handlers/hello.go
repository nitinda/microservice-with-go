package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Hello is a simple handler
type Hello struct {
	l *log.Logger
}

// NewHello Method creats new Hello handler with logger
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

// ServeHTTP implements the go http.handler interface
func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello Print")

	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Ohh Nooooo", http.StatusBadRequest)
	}

	fmt.Fprintf(rw, "Hello %s\n", d)
}
