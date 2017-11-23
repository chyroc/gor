package gor

import (
	"fmt"
	"net/http"
)

// Router router of api
type Router interface {
	Get(pattern string, h HandlerFunc)
	Listen(port int)
}

type HandlerFunc func(*Req, Res)

type Gor struct {
	handlers map[string]HandlerFunc
}

func NewGor() *Gor {
	return &Gor{
		handlers: make(map[string]HandlerFunc),
	}
}

func (g *Gor) Get(pattern string, h HandlerFunc) {
	g.handlers[http.MethodGet+pattern] = h
}

func (g *Gor) Post(pattern string, h HandlerFunc) {
	g.handlers[http.MethodPost+pattern] = h
}

func (g *Gor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httpMethod := r.Method
	path := r.URL.Path

	if handle, ok := g.handlers[httpMethod+path]; ok {
		handle(httpRequestToReq(r), httpResponseWriterToRes(w))
		return
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, ErrNotFound)
	return
}

func (g *Gor) Listen(addr string) error {
	return http.ListenAndServe(addr, g)
}
