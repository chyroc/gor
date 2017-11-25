package gor

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (g *Gor) handlerRoute(method string, pattern string, h HandlerFunc) {
	if !strings.HasPrefix(pattern, "/") {
		panic("not start with /")
	}

	URL, err := url.Parse(pattern)
	if err != nil {
		panic(fmt.Sprintf("pattern invalid: %s", pattern))
	}

	paths := strings.Split(URL.Path[1:], "/")
	var prepath string
	if strings.HasPrefix(paths[0], ":") {
		prepath = "/"
	} else {
		prepath = paths[0]
		paths = paths[1:]
	}

	//g.midWithPath[route] = len(g.middlewares) - 1 todo
	var rps []*routeParam
	for _, i := range paths {
		if strings.HasPrefix(i, ":") {
			rps = append(rps, &routeParam{name: i[1:], isParam: true})
		} else {
			rps = append(rps, &routeParam{name: i, isParam: false})
		}
	}

	g.routes = append(g.routes, &route{
		method:      method,
		handler:     h,
		prepath:     prepath,
		routeParams: rps,
	})
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
