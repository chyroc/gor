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
}

// NewRouter return *Router
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
func (r *Route) Use(h HandlerFunc) {
	//r.handlerRoute(http.MethodTrace, pattern, h)
}
func splitPattern(pattern string) []string {
	//var rps []*routeParam
	if strings.HasPrefix(pattern, "/") {
		pattern = pattern[1:]
	}
	return strings.Split(pattern, "/")
	//for _, v := range paths {
	//	if strings.HasPrefix(v, ":") {
	//		rps = append(rps, &routeParam{name: v[1:], isParam: true})
	//	} else {
	//		rps = append(rps, &routeParam{name: v, isParam: false})
	//	}
	//}
	//return rps
}

// UseN http trace method
func (r *Route) UseN(pattern string, m Mid) {
	subroutes := m.Handler(pattern)
	//mainroute := r.matchMainRoute(pattern)
	patternPaths := splitPattern(pattern)

	matchParam, matchIndex := matchRouter2("ALL", patternPaths, r.routes)
	fmt.Printf("matchParam,matchIndex %s %s\n", matchParam, matchIndex)
	fmt.Printf("patternPaths %s %s\n", patternPaths, len(patternPaths))

	//patternPaths : `['']`, `['a']`
	var routeParams []*routeParam
	if matchIndex == -1 {
		for _, v := range patternPaths[1:] {
			if strings.HasPrefix(v, ":") {
				routeParams = append(routeParams, &routeParam{name: v, isParam: true})
			} else {
				routeParams = append(routeParams, &routeParam{name: v, isParam: false})
			}
		}
	}

	for _, v := range subroutes {
		var subRouteParams []*routeParam
		subRouteParams = append(routeParams, v.routeParams...)
		r.routes = append(r.routes, &route{
			prepath:     patternPaths[0],
			routeParams: subRouteParams,
		})
	}

	//r.handlerRoute(http.MethodTrace, pattern, h)
	//r.routes = append(r.routes, m.Handler(pattern)...)
}

// Handler http trace method
func (r *Route) Handler(pattern string) []*route {
	subroutes := copyRouteSlice(r.routes)

	//for _, subroute := range routes {
	//	fmt.Printf("subroute.prepath %s \n", subroute.prepath)
	//	fmt.Printf("subroute.routeParams %s \n", subroute.routeParams)
	//	fmt.Printf("pattern %s\n",pattern)
	//	if strings.HasPrefix(pattern, "/:") {
	//		//fmt.Printf("route.routeParams %s\n", route.routeParams)
	//		//subroute.routeParams = append([]*routeParam{{name: pattern[2:], isParam: true}, {name: subroute.prepath, isParam: false}}, subroute.routeParams...)
	//		subroute.routeParams = append([]*routeParam{{name: subroute.prepath, isParam: true}}, subroute.routeParams...)
	//		//fmt.Printf("route.routeParams %s\n", route.routeParams)
	//		subroute.prepath = ""
	//	} else {
	//		subroute.routeParams = append([]*routeParam{{name: subroute.prepath, isParam: false}}, subroute.routeParams...)
	//		subroute.prepath = pattern[1:]
	//	}
	//}
	return subroutes
}

type Mid interface {
	Handler(pattern string) []*route
}
