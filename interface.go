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

type methodWithNext interface {
	GetWithNext(pattern string, hs ...HandlerFuncWithNext)
	HeadWithNext(pattern string, h HandlerFuncWithNext)
	PostWithNext(pattern string, h HandlerFuncWithNext)
	PutWithNext(pattern string, h HandlerFuncWithNext)
	PatchWithNext(pattern string, h HandlerFuncWithNext)
	DeleteWithNext(pattern string, h HandlerFuncWithNext)
	ConnectWithNext(pattern string, h HandlerFuncWithNext)
	OptionsWithNext(pattern string, h HandlerFuncWithNext)
	TraceWithNext(pattern string, h HandlerFuncWithNext)
}

var _ gorInterface = (*Gor)(nil)
var _ normalMethod = (*Gor)(nil)
var _ methodWithNext = (*Gor)(nil)
