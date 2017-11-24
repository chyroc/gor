package gor

import (
	"fmt"
	"net/http"
)

// ServeHTTP use to start server
func (g *Gor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httpMethod := r.Method
	path := r.URL.Path
	req, res := httpRequestToReq(r), httpResponseWriterToRes(w)

	if handle, ok := g.handlers[httpMethod+path]; ok {
		for i := 0; i <= g.midWithPath[httpMethod+path]; i++ {
			g.middlewares[i](g).ServeHTTP(w, r)
		}
		handle(req, res)
		return
	}

	if handles, ok := g.ttt[httpMethod+path]; ok {
		for _, handle := range handles {
			handle(req, res)
			fmt.Printf("req3 %s \n", req.context)
		}
		return
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, http.StatusText(http.StatusNotFound))
}

// Use add middlewares
func (g *Gor) Use(middlewares ...func(g *Gor) http.Handler) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// Listen bind port and start server
func (g *Gor) Listen(addr string) error {
	return http.ListenAndServe(addr, g)
}
