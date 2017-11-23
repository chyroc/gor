package gor

import (
	"fmt"
	"net/http"
)

// Router router of api
type Router interface {
	Get(pattern string, h HandlerFunc)
	Listen(port int)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
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

func (g *Gor) Head(pattern string, h HandlerFunc) {
	g.handlers[http.MethodHead+pattern] = h
}

func (g *Gor) Post(pattern string, h HandlerFunc) {
	g.handlers[http.MethodPost+pattern] = h
}

func (g *Gor) Put(pattern string, h HandlerFunc) {
	g.handlers[http.MethodPut+pattern] = h
}

func (g *Gor) Patch(pattern string, h HandlerFunc) {
	g.handlers[http.MethodPatch+pattern] = h
}

func (g *Gor) Delete(pattern string, h HandlerFunc) {
	g.handlers[http.MethodDelete+pattern] = h
}

func (g *Gor) Connect(pattern string, h HandlerFunc) {
	g.handlers[http.MethodConnect+pattern] = h
}

func (g *Gor) Options(pattern string, h HandlerFunc) {
	g.handlers[http.MethodOptions+pattern] = h
}

func (g *Gor) Trace(pattern string, h HandlerFunc) {
	g.handlers[http.MethodTrace+pattern] = h
}

func (g *Gor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httpMethod := r.Method
	path := r.URL.Path

	if handle, ok := g.handlers[httpMethod+path]; ok {
		handle(httpRequestToReq(r), httpResponseWriterToRes(w))
		return
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, http.StatusText(http.StatusNotFound))
	return
}

func (g *Gor) Listen(addr string) error {
	return http.ListenAndServe(addr, g)
}
