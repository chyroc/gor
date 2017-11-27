package gor

import (
	"net/http"
	"strings"
)

func matchRouter(method string, paths []string, routes []*route) []*matchedRoute {
	debugPrintf("routes len %s", len(routes))
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

	var matchedRoutes []*matchedRoute
	for matchIndex, route := range routes {
		if route.prepath == paths[0] {
			if route.method != "ALL" && route.method != method {
				continue
			}
			matchRoutes := paths[1:]
			if len(matchRoutes) != len(route.routeParams) {
				continue
			}

			if len(route.routeParams) == 0 && len(matchRoutes) == 0 {
				matchedRoutes = append(matchedRoutes, &matchedRoute{index: matchIndex})
				continue
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
				matchedRoutes = append(matchedRoutes, &matchedRoute{index: matchIndex, params: matchParams})
				continue
			}
		}
	}

	return matchedRoutes
}

// ServeHTTP use to start server
func (g *Gor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := httpResponseWriterToRes(w)
	req, err := httpRequestToReq(r)
	if err != nil {
		res.Error(err.Error())
		return
	}

	routes := strings.Split(strings.Split(r.URL.Path[1:], "?")[0], "/")
	matchedRoutes := matchRouter(r.Method, routes, g.routes)
	debugPrintf("matchedRoutes len %d", len(matchedRoutes))
	canNext := true
	isFirst := true
	for i, matchedRoute_ := range matchedRoutes {
		canNext = false
		debugPrintf("res exit %s", res.exit)
		if res.exit {
			return
		}
		if !isFirst || canNext {
			return
		}
		debugPrintf("matchedRoute[%d] %+v", i, matchedRoute_)
		route := g.routes[matchedRoute_.index]
		matchParams := matchedRoute_.params
		debugPrintf("route %+v \nmatchParams %+v", route, matchParams)
		for k, v := range matchParams {
			req.Params[k] = v
		}
		if route.handlerFunc != nil {
			route.handlerFunc(req, res)
			return
			//continue
		} else if route.handlerFuncNext != nil {
			isFirst = false
			route.handlerFuncNext(req, res, func() {
				canNext = true
			})
			continue
		} else if route.middleware != nil {
			route.middleware.handler("")
			continue
		}

		panic("")
	}

	res.SendStatus(http.StatusNotFound)
}

// Listen bind port and start server
func (g *Gor) Listen(addr string) error {
	return http.ListenAndServe(addr, g)
}
