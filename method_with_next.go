package gor

import (
	"net/http"
)

// HandlerFuncWithNext gor handler func like http.HandlerFunc func(ResponseWriter, *Request) but add next func
type HandlerFuncWithNext func(*Req, Res, Next)

func (hn *HandlerFuncWithNext) ServeHTTP(w http.ResponseWriter, res *http.Request) {
	// todo
}

// GetWithNext http get method
func (g *Gor) GetWithNext(pattern string, h ...HandlerFuncWithNext) {
	//g.handlersWithNext[http.MethodGet+pattern] = h
	g.ttt[http.MethodGet+pattern] = h
}

// HeadWithNext http head method
func (g *Gor) HeadWithNext(pattern string, h HandlerFuncWithNext) {
	g.handlersWithNext[http.MethodHead+pattern] = h
}

// PostWithNext http post method
func (g *Gor) PostWithNext(pattern string, h HandlerFuncWithNext) {
	g.handlersWithNext[http.MethodPost+pattern] = h
}

// PutWithNext http put method
func (g *Gor) PutWithNext(pattern string, h HandlerFuncWithNext) {
	g.handlersWithNext[http.MethodPut+pattern] = h
}

// PatchWithNext http patch method
func (g *Gor) PatchWithNext(pattern string, h HandlerFuncWithNext) {
	g.handlersWithNext[http.MethodPatch+pattern] = h
}

// DeleteWithNext http delete method
func (g *Gor) DeleteWithNext(pattern string, h HandlerFuncWithNext) {
	g.handlersWithNext[http.MethodDelete+pattern] = h
}

// ConnectWithNext http connect method
func (g *Gor) ConnectWithNext(pattern string, h HandlerFuncWithNext) {
	g.handlersWithNext[http.MethodConnect+pattern] = h
}

// OptionsWithNext http options method
func (g *Gor) OptionsWithNext(pattern string, h HandlerFuncWithNext) {
	g.handlersWithNext[http.MethodOptions+pattern] = h
}

// TraceWithNext http trace method
func (g *Gor) TraceWithNext(pattern string, h HandlerFuncWithNext) {
	g.handlersWithNext[http.MethodTrace+pattern] = h
}
