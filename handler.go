package gor

import (
	"net/http"
	"strings"
)

type matchedRouteArray []*route

func recursionMatch(method, requestPath, prePath string, parentRoutes []*route, matchedRoutes *matchedRouteArray) {
	if !strings.HasPrefix(requestPath, "/") {
		requestPath = "/" + requestPath
	}
	if strings.HasSuffix(requestPath, "/") {
		requestPath = requestPath[:len(requestPath)-1]
	}
	for _, route := range parentRoutes {
		if route.method != "ALL" && route.method != method {
			continue
		}
		_, matched := matchPath(route.routePath, requestPath, route.matchType)

		if matched {
			if len(route.children) > 0 {
				subrequestPath := strings.Join(strings.Split(requestPath, "/")[2:], "/")
				recursionMatch(method, subrequestPath, prePath+route.routePath, route.children, matchedRoutes)
			} else {
				route2 := route.copy()
				route2.routePath = prePath + route2.routePath
				(*matchedRoutes) = append((*matchedRoutes), route2)
			}
		}
	}
}

func matchRouter(method string, requestPath string, routes []*route) []*route {
	if strings.ContainsRune(requestPath, '?') {
		requestPath = strings.Split(requestPath, "?")[0]
	}
	var matchedRoutes matchedRouteArray
	recursionMatch(method, requestPath, "", routes, &matchedRoutes)
	return matchedRoutes
}

func doHandler(req *Req, res *Res, index int, matchRoutes []*route, requestPath string) {
	for i, j := index, len(matchRoutes); i < j; i++ {
		if res.exit {
			return
		}

		route := matchRoutes[i]
		req.Params, _ = matchPath(route.routePath, requestPath, route.matchType)

		if route.handlerFunc != nil {
			route.handlerFunc(req, res)
		} else if route.handlerFuncNext != nil {
			noCallNext := true
			route.handlerFuncNext(req, res, func() {
				noCallNext = false
				doHandler(req, res, index+1, matchRoutes, requestPath)
			})
			if noCallNext {
				res.exit = true
				return
			}
		} else {
			panic("sdafafdasfasfasfas")
		}
	}
}

// ServeHTTP use to start server
func (g *Gor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := httpResponseWriterToRes(w)
	req, err := httpRequestToReq(r)
	if err != nil {
		res.Error(err.Error())
		return
	}

	requestPath := strings.Split(r.URL.Path, "?")[0]
	matchedRoutes := matchRouter(r.Method, requestPath, g.routes)

	doHandler(req, res, 0, matchedRoutes, requestPath)

	res.SendStatus(http.StatusNotFound)
}

// Listen bind port and start server
func (g *Gor) Listen(addr string) error {
	return http.ListenAndServe(addr, g)
}
