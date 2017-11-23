package gor

import (
	"fmt"
	"net/http"
)

// ServeHTTP use to start server
func (g *Gor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httpMethod := r.Method
	path := r.URL.Path

	if handle, ok := g.handlers[httpMethod+path]; ok {
		handle(httpRequestToReq(r), httpResponseWriterToRes(w))
		return
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, http.StatusText(http.StatusNotFound))
}

// Listen bind port and start server
func (g *Gor) Listen(addr string) error {
	return http.ListenAndServe(addr, g)
}
