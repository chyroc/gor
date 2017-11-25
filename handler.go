package gor

import (
	"fmt"
	"net/http"
)

// ServeHTTP use to start server
func (g *Gor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httpMethod := r.Method
	path := r.URL.Path
	route := httpMethod + path

	res := httpResponseWriterToRes(w)
	req, err := httpRequestToReq(r)
	if err != nil {
		res.Error(err.Error())
		return
	}

	// normal method
	if handle, ok := g.handlers[route]; ok {
		for i := 0; i <= g.midWithPath[route]; i++ {
			g.middlewares[i](g).ServeHTTP(w, r)
		}
		handle(req, res)
		return
	}

	// todo method with next
	if handles, ok := g.ttt[route]; ok {
		for _, handle := range handles {
			handle(req, res)
			fmt.Printf("req3 %s \n", req.context)
		}
		return
	}

	res.SendStatus(http.StatusNotFound)
}

// Use add middlewares
func (g *Gor) Use(middlewares ...func(g *Gor) http.Handler) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// Listen bind port and start server
func (g *Gor) Listen(addr string) error {
	return http.ListenAndServe(addr, g)
}
