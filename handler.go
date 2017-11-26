package gor

import (
	"fmt"
	"net/http"
	"strings"
)

func matchRouter(method string, paths []string, routes []*route) (map[string]string, int) {
	for k, v := range routes {
		debugPrintf("route[%d] %+v", k, v)
		for k2, v2 := range v.routeParams {
			debugPrintf("routeParams[%d] %+v", k2, v2)
		}
		debugPrintf("=====")
	}

	for _, v := range paths {
		if strings.Contains(v, "/") {
			panic("paths cannot contain /")
		}
	}

	matchIndex := -1
	for _, route := range routes {
		matchIndex++
		if route.prepath == paths[0] {
			if method != "ALL" && route.method != "ALL" && route.method != method {
				continue
			}
			matchRoutes := paths[1:]
			if len(matchRoutes) != len(route.routeParams) {
				continue
			}

			if len(route.routeParams) == 0 && len(matchRoutes) == 0 {
				return nil, matchIndex
			}

			match := false
			matchParams := make(map[string]string)
			for i, j := 0, len(matchRoutes); i < j; i++ {
				if route.routeParams[i].isParam {
					matchParams[route.routeParams[i].name] = matchRoutes[i]
				} else if route.routeParams[i].name != matchRoutes[i] {
					match = false
					break
				}
				match = true
			}
			if match {
				return matchParams, matchIndex
			}
		}
	}

	return nil, -1
}

// ServeHTTP use to start server
func (g *Gor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := httpResponseWriterToRes(w)
	req, err := httpRequestToReq(r)
	if err != nil {
		res.Error(err.Error())
		return
	}

	fmt.Printf("all mids: %+v\n", g.mids)
	for _, mid := range g.mids {
		if deferFunc := mid(req, res); deferFunc != nil {
			fmt.Printf("deferFunc %v\n", deferFunc)
			defer deferFunc(req, res)
		}
	}

	routes := strings.Split(strings.Split(r.URL.Path[1:], "?")[0], "/")
	matchParams, matchIndex := matchRouter(r.Method, routes, g.routes)
	if matchIndex != -1 {
		for k, v := range matchParams {
			req.Params[k] = v
		}
		g.routes[matchIndex].handler(req, res)
		return
	}

	res.SendStatus(http.StatusNotFound)
}

// Use add middlewares
func (g *Gor) Use(middlewares ...func(g *Gor) http.Handler) {
	//g.middlewares = append(g.middlewares, middlewares...)
}

// Listen bind port and start server
func (g *Gor) Listen(addr string) error {
	return http.ListenAndServe(addr, g)
}
