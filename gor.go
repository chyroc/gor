package gor

import (
	"net/http"
)

// HandlerFunc gor handler func like http.HandlerFunc func(ResponseWriter, *Request)
type HandlerFunc func(*Req, *Res)

type routeParam struct {
	name    string
	isParam bool
}

type route struct {
	method      string
	handler     HandlerFunc
	prepath     string
	routeParams []*routeParam
}

// Gor gor framework core struct
type Gor struct {
	middlewares []func(g *Gor) http.Handler
	midWithPath map[string]int

	routes []*route
}

// NewGor return Gor struct
func NewGor() *Gor {
	return &Gor{
		midWithPath: make(map[string]int),
	}
}
