package gor

import (
	"net/http"
)

type gorInterface interface {
	normalMethod
	http.Handler
}

type normalMethod interface {
	Get(pattern string, h HandlerFunc)
	Head(pattern string, h HandlerFunc)
	Post(pattern string, h HandlerFunc)
	Put(pattern string, h HandlerFunc)
	Patch(pattern string, h HandlerFunc)
	Delete(pattern string, h HandlerFunc)
	Connect(pattern string, h HandlerFunc)
	Options(pattern string, h HandlerFunc)
	Trace(pattern string, h HandlerFunc)
}

type routerInterface interface {
	All()
	Method()
	Param()
	Route()
	Use(middlewares ...HandlerFunc)
}

type RouteInterface interface {
	Use(h HandlerFunc)
	UseN(pattern string, m Mid)

	normalMethod
}

var _ gorInterface = (*Gor)(nil)
var _ normalMethod = (*Gor)(nil)
var _ RouteInterface = (*Route)(nil)
