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

func (g *Gor) handlerRoute(method string, pattern string, h HandlerFunc) {
	route := method + pattern
	g.handlers[route] = h
	g.midWithPath[route] = len(g.middlewares) - 1
}

// Get http get method
func (g *Gor) Get(pattern string, h HandlerFunc) {
	g.handlerRoute(http.MethodGet, pattern, h)
}

// Head http head method
func (g *Gor) Head(pattern string, h HandlerFunc) {
	g.handlerRoute(http.MethodHead, pattern, h)
}

// Post http post method
func (g *Gor) Post(pattern string, h HandlerFunc) {
	g.handlerRoute(http.MethodPost, pattern, h)
}

// Put http put method
func (g *Gor) Put(pattern string, h HandlerFunc) {
	g.handlerRoute(http.MethodPut, pattern, h)
}

// Patch http patch method
func (g *Gor) Patch(pattern string, h HandlerFunc) {
	g.handlerRoute(http.MethodPatch, pattern, h)
}

// Delete http delete method
func (g *Gor) Delete(pattern string, h HandlerFunc) {
	g.handlerRoute(http.MethodDelete, pattern, h)
}

// Connect http connect method
func (g *Gor) Connect(pattern string, h HandlerFunc) {
	g.handlerRoute(http.MethodConnect, pattern, h)
}

// Options http options method
func (g *Gor) Options(pattern string, h HandlerFunc) {
	g.handlerRoute(http.MethodOptions, pattern, h)
}

// Trace http trace method
func (g *Gor) Trace(pattern string, h HandlerFunc) {
	g.handlerRoute(http.MethodTrace, pattern, h)
}
