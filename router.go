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

// HandlerFunc gor handler func like http.HandlerFunc func(ResponseWriter, *Request)
type HandlerFunc func(*Req, Res)

// Gor gor framework core struct
type Gor struct {
	handlers map[string]HandlerFunc
}

// NewGor return Gor struct
func NewGor() *Gor {
	return &Gor{
		handlers: make(map[string]HandlerFunc),
	}
}

// Get http get method
func (g *Gor) Get(pattern string, h HandlerFunc) {
	g.handlers[http.MethodGet+pattern] = h
}

// Head http head method
func (g *Gor) Head(pattern string, h HandlerFunc) {
	g.handlers[http.MethodHead+pattern] = h
}

// Post http post method
func (g *Gor) Post(pattern string, h HandlerFunc) {
	g.handlers[http.MethodPost+pattern] = h
}

// Put http put method
func (g *Gor) Put(pattern string, h HandlerFunc) {
	g.handlers[http.MethodPut+pattern] = h
}

// Patch http patch method
func (g *Gor) Patch(pattern string, h HandlerFunc) {
	g.handlers[http.MethodPatch+pattern] = h
}

// Delete http delete method
func (g *Gor) Delete(pattern string, h HandlerFunc) {
	g.handlers[http.MethodDelete+pattern] = h
}

// Connect http connect method
func (g *Gor) Connect(pattern string, h HandlerFunc) {
	g.handlers[http.MethodConnect+pattern] = h
}

// Options http options method
func (g *Gor) Options(pattern string, h HandlerFunc) {
	g.handlers[http.MethodOptions+pattern] = h
}

// 1Trace http trace method
func (g *Gor) Trace(pattern string, h HandlerFunc) {
	g.handlers[http.MethodTrace+pattern] = h
}

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
	return
}

// Listen bind port and start server
func (g *Gor) Listen(addr string) error {
	return http.ListenAndServe(addr, g)
}
