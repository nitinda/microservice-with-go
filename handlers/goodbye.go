package handlers

import (
	"log"
	"net/http"
)

// GoodBye is a simple handler
type GoodBye struct {
	l *log.Logger
}

// NewGoodBye Method creats new GoodBye handler with logger
func NewGoodBye(l *log.Logger) *GoodBye {
	return &GoodBye{l}
}

// ServeHTTP implements the go http.handler interface
func (g *GoodBye) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Byeee\n"))
}
