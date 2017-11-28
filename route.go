package gor

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

// HandlerFunc gor handler func like http.HandlerFunc func(ResponseWriter, *Request)
type HandlerFunc func(*Req, *Res)

// Next exec next handler or mid
type Next func()

// HandlerFuncNext gor handler func like http.HandlerFunc func(ResponseWriter, *Request),
// but return HandlerFunc to do somrthing at defer time
type HandlerFuncNext func(*Req, *Res, Next)

type routeParam struct {
	name    string
	isParam bool
}

type route struct {
	handlerFunc     HandlerFunc
	handlerFuncNext HandlerFuncNext
	middleware      Middleware

	method      string
	prepath     string
	routeParams []*routeParam

	//parentIndex string
}
type matchedRoute struct {
	index  int
	params map[string]string
}

func (r *route) copy() *route {
	var t = &route{
		handlerFunc:     r.handlerFunc,
		handlerFuncNext: r.handlerFuncNext,
		middleware:      r.middleware,

		method:  r.method,
		prepath: r.prepath,

		//parentIndex: r.parentIndex,
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

func copyRouteParamSlice(routeParams []*routeParam) []*routeParam {
	var rs []*routeParam
	for _, v := range routeParams {
		rs = append(rs, &routeParam{
			name:    v.name,
			isParam: v.isParam,
		})
	}
	return rs
}

func mergeRouteParamSlice(parentRouteParams []*routeParam, subRouteParams ...*routeParam) []*routeParam {
	var rs []*routeParam
	for _, v := range parentRouteParams {
		rs = append(rs, &routeParam{
			name:    v.name,
			isParam: v.isParam,
		})
	}
	for _, v := range subRouteParams {
		rs = append(rs, &routeParam{
			name:    v.name,
			isParam: v.isParam,
		})
	}
	return rs
}

//[]*routeParam
// Route route
type Route struct {
	routes []*route
	mids   []HandlerFuncNext
}

// NewRoute return *Router
func NewRoute() *Route {
	return &Route{}
}

func (r *Route) addHandlerFuncRoute(method string, pattern string, h HandlerFunc, parentrouteParams []*routeParam) {
	fmt.Printf("method %s, pattern %s, h %s, parentrouteParams%s\n ", method, pattern, h, parentrouteParams)
	if pattern == "*" {
		pattern = "/" + pattern
	}
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

	//var rps []*routeParam
	for _, i := range paths {
		if strings.HasPrefix(i, ":") {
			//parentrouteParams = append(parentrouteParams, &routeParam{name: i[1:], isParam: true})
			parentrouteParams = mergeRouteParamSlice(parentrouteParams, &routeParam{name: i[1:], isParam: true})
		} else {
			//parentrouteParams = append(parentrouteParams, &routeParam{name: i, isParam: false})
			parentrouteParams = mergeRouteParamSlice(parentrouteParams, &routeParam{name: i, isParam: false})
		}
	}
	//fmt.Printf("handler func params %+v\n", parentrouteParams)
	r.routes = append(r.routes, &route{
		handlerFunc: h,

		method:      method,
		prepath:     prepath,
		routeParams: parentrouteParams,

		//parentIndex: pattern,
	})
}

func (r *Route) addHandlerFuncNextRoute(method string, pattern string, h HandlerFuncNext, parentrouteParams []*routeParam) {
	//fmt.Printf("pattern", pattern)
	if pattern == "*" {
		pattern = "/" + pattern
	}
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

	//var rps []*routeParam
	for _, i := range paths {
		if strings.HasPrefix(i, ":") {
			//parentrouteParams = append(parentrouteParams, &routeParam{name: i[1:], isParam: true})
			parentrouteParams = mergeRouteParamSlice(parentrouteParams, &routeParam{name: i[1:], isParam: true})
		} else {
			//parentrouteParams = append(parentrouteParams, &routeParam{name: i, isParam: false})
			parentrouteParams = mergeRouteParamSlice(parentrouteParams, &routeParam{name: i, isParam: false})
		}
	}

	r.routes = append(r.routes, &route{
		handlerFuncNext: h,

		method:      method,
		prepath:     prepath,
		routeParams: parentrouteParams,

		//parentIndex: pattern,
	})
}

func (r *Route) addMiddlerwareRoute(method string, pattern string, mid Middleware) {
	if pattern == "*" {
		pattern = "/" + pattern
	}
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

	var rps []*routeParam
	for _, i := range paths {
		if strings.HasPrefix(i, ":") {
			rps = append(rps, &routeParam{name: i[1:], isParam: true})
		} else {
			rps = append(rps, &routeParam{name: i, isParam: false})
		}
	}

	r.routes = append(r.routes, &route{
		middleware: mid,

		method:      method,
		prepath:     prepath,
		routeParams: rps,

		//parentIndex: pattern,
	})
}

// Get http get method
func (r *Route) Get(pattern string, h HandlerFunc) {
	r.addHandlerFuncRoute(http.MethodGet, pattern, h, []*routeParam{})
}

// Head http head method
func (r *Route) Head(pattern string, h HandlerFunc) {
	r.addHandlerFuncRoute(http.MethodHead, pattern, h, []*routeParam{})
}

// Post http post method
func (r *Route) Post(pattern string, h HandlerFunc) {
	r.addHandlerFuncRoute(http.MethodPost, pattern, h, []*routeParam{})
}

// Put http put method
func (r *Route) Put(pattern string, h HandlerFunc) {
	r.addHandlerFuncRoute(http.MethodPut, pattern, h, []*routeParam{})
}

// Patch http patch method
func (r *Route) Patch(pattern string, h HandlerFunc) {
	r.addHandlerFuncRoute(http.MethodPatch, pattern, h, []*routeParam{})
}

// Delete http delete method
func (r *Route) Delete(pattern string, h HandlerFunc) {
	r.addHandlerFuncRoute(http.MethodDelete, pattern, h, []*routeParam{})
}

// Connect http connect method
func (r *Route) Connect(pattern string, h HandlerFunc) {
	r.addHandlerFuncRoute(http.MethodConnect, pattern, h, []*routeParam{})
}

// Options http options method
func (r *Route) Options(pattern string, h HandlerFunc) {
	r.addHandlerFuncRoute(http.MethodOptions, pattern, h, []*routeParam{})
}

// Trace http trace method
func (r *Route) Trace(pattern string, h HandlerFunc) {
	r.addHandlerFuncRoute(http.MethodTrace, pattern, h, []*routeParam{})
}

// Use http trace method
//
// must belong one of below type
// string
// type HandlerFunc func(*Req, *Res)
// type HandlerFuncNext func(*Req, *Res, Next)
// type Middleware interface
func (r *Route) Use(hs ...interface{}) {
	if len(hs) == 1 {
		r.useWithOne("*", hs[0])
		return
	}

	first := hs[0]
	firstType := reflect.TypeOf(first)
	if firstType.Kind() == reflect.String {
		firstValue := reflect.ValueOf(first)
		pattern := firstValue.String()
		for _, h := range hs[1:] {
			r.useWithOne(pattern, h)
		}
	} else {
		for _, h := range hs {
			r.useWithOne("*", h)
		}
	}
}

func (r *Route) useWithOne(pattern string, h interface{}) {
	// todo use 应该处理签名的params
	var err error = nil
	defer func() {
		if err != nil {
			panic(err)
		}
	}()
	hType := reflect.TypeOf(h)
	switch hType.Kind() {
	case reflect.Func:
		switch h.(type) {
		case func(req *Req, res *Res):
			if f, ok := h.(func(req *Req, res *Res)); ok {
				r.useWithHandlerFunc("ALL", pattern, HandlerFunc(f), []*routeParam{})
			} else {
				err = fmt.Errorf("cannot convert to gor.HandlerFunc")
			}
		case func(req *Req, res *Res, next Next):
			if f, ok := h.(func(req *Req, res *Res, next Next)); ok {
				r.useWithHandlerFuncNext("ALL", pattern, HandlerFuncNext(f), []*routeParam{}) // todo parentrouteParams
			} else {
				err = fmt.Errorf("cannot convert to gor.HandlerFuncNext")
			}
		default:
			err = fmt.Errorf("maybe you are transmiting gor.HandlerFunc / gor.HandlerFuncNext, but the function signature is wrong")
		}
	case reflect.Struct:
		err = fmt.Errorf("maybe you are transmiting gor.Middleware, but please use Pointer, not Struct")
	case reflect.Ptr:
		if f, ok := h.(Middleware); ok {
			r.useWithMiddleware("ALL", pattern, Middleware(f), []*routeParam{}) // todo parentrouteParams
		} else {
			err = fmt.Errorf("cannot convert to gor.Middleware")
		}
	default:
		err = fmt.Errorf("when middleware length is one, that type must belong gor.HandlerFunc / gor.HandlerFuncNext / gor.Route, but get %s", hType.Kind())
	}
}

func (r *Route) handler(pattern string) []*route {
	return copyRouteSlice(r.routes)
}

// UseN http trace method
func (r *Route) useWithHandlerFunc(method, pattern string, h HandlerFunc, parentrouteParams []*routeParam) {
	//fmt.Printf("get pattern: %s, HandlerFunc: %+v\n", pattern, h)
	r.addHandlerFuncRoute(method, pattern, h, parentrouteParams)
}

// UseN http trace method
func (r *Route) useWithHandlerFuncNext(method, pattern string, h HandlerFuncNext, parentrouteParams []*routeParam) {
	//fmt.Printf("get pattern: %s, HandlerFuncNext: %+v\n", pattern, h)
	r.addHandlerFuncNextRoute(method, pattern, h, parentrouteParams)
}

// UseN http trace method
func (r *Route) useWithMiddleware(method, pattern string, mid Middleware, parentrouteParams []*routeParam) {
	//parentRoutes := copyRouteSlice(r.routes)
	subRoutes := mid.handler(pattern)
	//fmt.Printf("get pattern: %s, Middleware: %+v\n", pattern, mid)
	//fmt.Printf("parentRoutes %+v\n", parentRoutes)
	for _, subRoute := range subRoutes {
		fmt.Printf("subRoute %s\n", subRoute)
		var newParentrouteParams []*routeParam
		if subRoute.prepath != "" {
			fmt.Printf("1\n")
			if strings.HasPrefix(subRoute.prepath, ":") {
				fmt.Printf("2\n")
				newParentrouteParams = mergeRouteParamSlice(parentrouteParams, &routeParam{name: subRoute.prepath[1:], isParam: true})
			} else {
				fmt.Printf("3\n")
				newParentrouteParams = mergeRouteParamSlice(parentrouteParams, &routeParam{name: subRoute.prepath, isParam: false})
			}
			fmt.Printf("4\n")
		}
		fmt.Printf("5\n")
		newParentrouteParams = mergeRouteParamSlice(newParentrouteParams, subRoute.routeParams...)
		//fmt.Printf("newParentrouteParams %s\n", newParentrouteParams)
		if subRoute.handlerFunc != nil {
			fmt.Printf("6\n")
			r.routes = append(r.routes, &route{
				handlerFunc: subRoute.handlerFunc,
				method:      subRoute.method,
				prepath:     pattern[1:],
				routeParams: newParentrouteParams,
				//parentIndex: pattern,
			})
			//r.useWithHandlerFunc(subRoute.method, pattern+"/"+subRoute.prepath, subRoute.handlerFunc, newParentrouteParams)
		} else if subRoute.handlerFuncNext != nil {
			fmt.Printf("7\n")
			fmt.Printf("6\n")
			r.routes = append(r.routes, &route{
				handlerFuncNext: subRoute.handlerFuncNext,
				method:          subRoute.method,
				prepath:         pattern[1:],
				routeParams:     newParentrouteParams,
				//parentIndex: pattern,
			})
			//r.useWithHandlerFuncNext(subRoute.method, pattern+"/"+subRoute.prepath, subRoute.handlerFuncNext, newParentrouteParams)
		} else if subRoute.middleware != nil {
			fmt.Printf("8\n")
			r.useWithMiddleware(subRoute.method, pattern+"/"+subRoute.prepath, subRoute.middleware, newParentrouteParams)
		} else {
			fmt.Printf("9\n")
			panic("notklsadjlfajs")
		}
	}
}

//// UseN http trace method
//func (r *Route) UseN(pattern string, m Middleware) {
//	midRouter := m.handler(pattern)
//	patternPaths := strings.Split(strings.TrimPrefix(pattern, "/"), "/")
//	_, matchIndex := matchRouter("ALL", patternPaths, r.routes)
//
//	var routeParams []*routeParam
//	if matchIndex == -1 {
//		for _, v := range patternPaths[1:] {
//			routeParams = append(routeParams, &routeParam{name: v, isParam: strings.HasPrefix(v, ":")})
//		}
//	} else {
//		routeParams = append(routeParams, r.routes[matchIndex].routeParams...)
//	}
//
//	for _, v := range midRouter.routes {
//		var subRouteParams []*routeParam
//		if v.prepath != "" {
//			subRouteParams = append(routeParams, &routeParam{name: v.prepath, isParam: false})
//		}
//		subRouteParams = append(subRouteParams, v.routeParams...)
//		r.routes = append(r.routes, &route{
//			method:      "ALL",
//			handlerFunc: v.handlerFunc,
//			prepath:     patternPaths[0],
//			routeParams: subRouteParams,
//		})
//	}
//	r.mids = append(r.mids, midRouter.mids...)
//}
