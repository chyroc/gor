package gor

import (
	"net/http"
)

// HandlerFunc gor handler func like http.HandlerFunc func(ResponseWriter, *Request)
type HandlerFunc func(*Req, Res)

func (h *HandlerFunc) ServeHTTP(w http.ResponseWriter, res *http.Request) {
	// todo
}

func (h *HandlerFunc) Next() {

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

// Trace http trace method
func (g *Gor) Trace(pattern string, h HandlerFunc) {
	g.handlers[http.MethodTrace+pattern] = h
}
