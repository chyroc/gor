package gor

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

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

func (r *route) copy() *route {
	var t = &route{
		method:  r.method,
		handler: r.handler, // not deep copy
		prepath: r.prepath,
	}
	var rs []*routeParam
	for _, v := range r.routeParams {
		rs = append(rs, &routeParam{
			name:    v.name,
			isParam: v.isParam,
		})
	}
	t.routeParams = rs
	return t
}

func copyRouteSlice(routes []*route) []*route {
	var rs []*route
	for _, v := range routes {
		rs = append(rs, v.copy())
	}
	return rs
}

// Route route
type Route struct {
	routes []*route
	mids   []HandlerFuncDefer
}

// NewRoute return *Router
func NewRoute() *Route {
	return &Route{}
}

func (r *Route) handlerRoute(method string, pattern string, h HandlerFunc) {
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
		prepath = ""
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

	r.routes = append(r.routes, &route{
		method:      method,
		handler:     h,
		prepath:     prepath,
		routeParams: rps,
	})
}

func (r *Route) handlerMid(hd HandlerFuncDefer) {
	fmt.Printf("add mid 4\n")
	r.mids = append(r.mids, hd)
}

// Get http get method
func (r *Route) Get(pattern string, h HandlerFunc) {
	r.handlerRoute(http.MethodGet, pattern, h)
}

// Head http head method
func (r *Route) Head(pattern string, h HandlerFunc) {
	r.handlerRoute(http.MethodHead, pattern, h)
}

// Post http post method
func (r *Route) Post(pattern string, h HandlerFunc) {
	r.handlerRoute(http.MethodPost, pattern, h)
}

// Put http put method
func (r *Route) Put(pattern string, h HandlerFunc) {
	r.handlerRoute(http.MethodPut, pattern, h)
}

// Patch http patch method
func (r *Route) Patch(pattern string, h HandlerFunc) {
	r.handlerRoute(http.MethodPatch, pattern, h)
}

// Delete http delete method
func (r *Route) Delete(pattern string, h HandlerFunc) {
	r.handlerRoute(http.MethodDelete, pattern, h)
}

// Connect http connect method
func (r *Route) Connect(pattern string, h HandlerFunc) {
	r.handlerRoute(http.MethodConnect, pattern, h)
}

// Options http options method
func (r *Route) Options(pattern string, h HandlerFunc) {
	r.handlerRoute(http.MethodOptions, pattern, h)
}

// Trace http trace method
func (r *Route) Trace(pattern string, h HandlerFunc) {
	r.handlerRoute(http.MethodTrace, pattern, h)
}

// Use http trace method
func (r *Route) Use(h HandlerFuncDefer) {
	fmt.Printf("add mid 2\n")
	r.handlerMid(h)
}

// UseN http trace method
func (r *Route) UseN(pattern string, m Mid) {
	midRouter := m.handler(pattern)
	patternPaths := strings.Split(strings.TrimPrefix(pattern, "/"), "/")
	_, matchIndex := matchRouter("ALL", patternPaths, r.routes)

	var routeParams []*routeParam
	if matchIndex == -1 {
		for _, v := range patternPaths[1:] {
			routeParams = append(routeParams, &routeParam{name: v, isParam: strings.HasPrefix(v, ":")})
		}
	} else {
		routeParams = append(routeParams, r.routes[matchIndex].routeParams...)
	}

	for _, v := range midRouter.routes {
		var subRouteParams []*routeParam
		if v.prepath != "" {
			subRouteParams = append(routeParams, &routeParam{name: v.prepath, isParam: false})
		}
		subRouteParams = append(subRouteParams, v.routeParams...)
		r.routes = append(r.routes, &route{
			method:      "ALL",
			handler:     v.handler,
			prepath:     patternPaths[0],
			routeParams: subRouteParams,
		})
	}
	r.mids = append(r.mids, midRouter.mids...)
}

func (r *Route) handler(pattern string) *Route {
	return &Route{
		routes: copyRouteSlice(r.routes),
		mids:   r.mids,
	}
}

// Mid mid
type Mid interface {
	handler(pattern string) *Route
}
