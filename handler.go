package gor

import (
	"net/http"
	"strings"
	"fmt"
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
		if route.prepath == paths[0] || route.prepath == "*" {
			fmt.Printf("route.routeParams paths", route.routeParams, paths)
			if route.method != "ALL" && route.method != method {
				continue
			}
			matchRoutes := paths[1:]
			if len(matchRoutes) != len(route.routeParams) && route.prepath != "*"{
				continue
			}
			fmt.Printf("--1\n")
			fmt.Printf("route.routeParams matchRoutes", route.routeParams, matchRoutes)
			if len(route.routeParams) == 0 && len(matchRoutes) == 0 {
				fmt.Printf("--2\n")
				matchedRoutes = append(matchedRoutes, &matchedRoute{index: matchIndex})
				continue
			}
			fmt.Printf("--3\n")

			match := false
			matchParams := make(map[string]string)
			for i, j := 0, len(matchRoutes); i < j; i++ {
				fmt.Printf("--4\n")
				if route.routeParams[i].isParam {
					fmt.Printf("--5\n")
					matchParams[route.routeParams[i].name] = matchRoutes[i]
				} else if route.routeParams[i].name != matchRoutes[i] {
					fmt.Printf("--6\n")
					match = false
					break
				}
				fmt.Printf("--7\n")
				match = true
			}
			fmt.Printf("--8\n")
			if match {
				fmt.Printf("--9\n")
				matchedRoutes = append(matchedRoutes, &matchedRoute{index: matchIndex, params: matchParams})
				continue
			}
			fmt.Printf("--10\n")
		}
	}

	return matchedRoutes
}

func (g *Gor) doHandler(req *Req, res *Res, routeIndex int, matchedRoutes []*matchedRoute) bool {
	matchedRoute := matchedRoutes[routeIndex]
	//canNext = false
	//debugPrintf("res exit %s", res.exit)
	if res.exit {
		return true
	}
	//if !isFirst || canNext {
	//	return
	//}
	//debugPrintf("matchedRoute[%d] %+v", routeIndex, matchedRoute)
	route := g.routes[matchedRoute.index]
	matchParams := matchedRoute.params
	//debugPrintf("route %+v \nmatchParams %+v", route, matchParams)
	for k, v := range matchParams {
		req.Params[k] = v
	}
	if route.handlerFunc != nil {
		route.handlerFunc(req, res)
		return true
	} else if route.handlerFuncNext != nil {
		//canNext := false
		fmt.Printf("1-----\n")
		route.handlerFuncNext(req, res, func() {
			fmt.Printf("2-----\n")
			for i, _ := range matchedRoutes[routeIndex+1:] {
				fmt.Printf("3-----\n")
				g.doHandler(req, res, i, matchedRoutes)
				fmt.Printf("4-----\n")
			}
			fmt.Printf("5-----\n")
			//canNext = true
		})
		fmt.Printf("6-----\n")
		return true
	} else {
		panic("get uhdkshfhksdh")
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

	routes := strings.Split(strings.Split(r.URL.Path[1:], "?")[0], "/")
	matchedRoutes := matchRouter(r.Method, routes, g.routes)
	debugPrintf("g.routes len %d %s", len(g.routes), g.routes)
	//debugPrintf("matchedRoutes len %d %s", len(matchedRoutes), matchedRoutes)
	//canNext := true
	//isFirst := true
	for i, _ := range matchedRoutes {
		g.doHandler(req, res, i, matchedRoutes)
	}

	res.SendStatus(http.StatusNotFound)
}

// Listen bind port and start server
func (g *Gor) Listen(addr string) error {
	return http.ListenAndServe(addr, g)
}
