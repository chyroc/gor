package gor

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
)

// HandlerFunc gor handler func like http.HandlerFunc func(ResponseWriter, *Request)
type HandlerFunc func(*Req, *Res)

// Next exec next handler or mid
type Next func()

// HandlerFuncNext gor handler func like http.HandlerFunc func(ResponseWriter, *Request),
// but return HandlerFunc to do somrthing at defer time
type HandlerFuncNext func(*Req, *Res, Next)

type matchType int

const (
	preMatch matchType = iota
	fullMatch
)

type route struct {
	method    string
	routePath string
	matchType matchType

	routePathReg *regexp.Regexp

	handlerFunc     HandlerFunc
	handlerFuncNext HandlerFuncNext
	middleware      Middleware

	children []*route
}

func (r *route) copy() *route {
	return &route{
		method:    r.method,
		routePath: r.routePath,
		matchType: r.matchType,

		routePathReg: r.routePathReg,

		handlerFunc:     r.handlerFunc,
		handlerFuncNext: r.handlerFuncNext,
		middleware:      r.middleware,

		children: copyRouteSlice(r.children),
	}
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

// NewRoute return *Router
func NewRoute() *Route {
	return &Route{}
}

func (r *Route) addHandlerFuncAndNextRoute(method string, pattern string, matchType matchType, h HandlerFunc, hn HandlerFuncNext) {
	if !strings.HasPrefix(pattern, "/") {
		panic("must start with /")
	}
	if strings.HasSuffix(pattern, "/") && pattern != "/" {
		pattern = pattern[:len(pattern)-1]
	}

	URL, err := url.Parse(pattern)
	if err != nil {
		panic(fmt.Sprintf("pattern invalid: %s", pattern))
	}

	routePath := URL.Path

	var routeH = &route{
		method:    method,
		routePath: routePath,
		matchType: matchType,

		routePathReg: genMatchPathReg(routePath),
	}
	if h != nil {
		routeH.handlerFunc = h
	} else if hn != nil {
		routeH.handlerFuncNext = hn
	} else {
		panic("handlerFunc or handlerFuncNext cannot be both nil")
	}

	r.routes = append(r.routes, routeH)
}

// Get http get method
func (r *Route) Get(pattern string, h HandlerFunc) {
	r.addHandlerFuncAndNextRoute(http.MethodGet, pattern, fullMatch, h, nil)
}

// Head http head method
func (r *Route) Head(pattern string, h HandlerFunc) {
	r.addHandlerFuncAndNextRoute(http.MethodHead, pattern, fullMatch, h, nil)
}

// Post http post method
func (r *Route) Post(pattern string, h HandlerFunc) {
	r.addHandlerFuncAndNextRoute(http.MethodPost, pattern, fullMatch, h, nil)
}

// Put http put method
func (r *Route) Put(pattern string, h HandlerFunc) {
	r.addHandlerFuncAndNextRoute(http.MethodPut, pattern, fullMatch, h, nil)
}

// Patch http patch method
func (r *Route) Patch(pattern string, h HandlerFunc) {
	r.addHandlerFuncAndNextRoute(http.MethodPatch, pattern, fullMatch, h, nil)
}

// Delete http delete method
func (r *Route) Delete(pattern string, h HandlerFunc) {
	r.addHandlerFuncAndNextRoute(http.MethodDelete, pattern, fullMatch, h, nil)
}

// Connect http connect method
func (r *Route) Connect(pattern string, h HandlerFunc) {
	r.addHandlerFuncAndNextRoute(http.MethodConnect, pattern, fullMatch, h, nil)
}

// Options http options method
func (r *Route) Options(pattern string, h HandlerFunc) {
	r.addHandlerFuncAndNextRoute(http.MethodOptions, pattern, fullMatch, h, nil)
}

// Trace http trace method
func (r *Route) Trace(pattern string, h HandlerFunc) {
	r.addHandlerFuncAndNextRoute(http.MethodTrace, pattern, fullMatch, h, nil)
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
		r.useWithOne("/", hs[0])
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
			r.useWithOne("/", h)
		}
	}
}

func (r *Route) useWithOne(pattern string, h interface{}) {
	// todo use 应该处理签名的params
	var err error
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
				r.useWithHandlerFunc("ALL", pattern, preMatch, HandlerFunc(f))
			} else {
				err = fmt.Errorf("cannot convert to gor.HandlerFunc")
			}
		case func(req *Req, res *Res, next Next):
			if f, ok := h.(func(req *Req, res *Res, next Next)); ok {
				r.useWithHandlerFuncNext("ALL", pattern, preMatch, HandlerFuncNext(f)) // todo parentrouteParams
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
			r.useWithMiddleware("ALL", pattern, preMatch, f)
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

func (r *Route) useWithHandlerFunc(method, pattern string, matchType matchType, h HandlerFunc) {
	//fmt.Printf("get pattern: %s, HandlerFunc: %+v\n", pattern, h)
	r.addHandlerFuncAndNextRoute(method, pattern, matchType, h, nil)
}

func (r *Route) useWithHandlerFuncNext(method, pattern string, matchType matchType, h HandlerFuncNext) {
	//fmt.Printf("get pattern: %s, HandlerFuncNext: %+v\n", pattern, h)
	r.addHandlerFuncAndNextRoute(method, pattern, matchType, nil, h)
}

func (r *Route) useWithMiddleware(method, pattern string, matchType matchType, mid Middleware) {
	if method != "ALL" {
		panic("middleware method must be ALL")
	}
	subRoutes := mid.handler(pattern)

	// 重新计算subRoutes的正则
	for _, v := range subRoutes {
		v.routePathReg = genMatchPathReg(pattern + v.routePath)
	}

	parent := &route{
		method:    method,
		matchType: matchType,
		routePath: pattern,

		routePathReg: genMatchPathReg(pattern),

		children: subRoutes,
	}
	r.routes = append(r.routes, parent)
}
