package gor

import (
	"net/http"
	"strings"
	"fmt"
)

type matchedRouteIndex struct {
	index    int
	params   map[string]string
	children []*matchedRouteIndex
}

func nextMatchedRouteIndex(index int, matchedRouteIndex []*matchedRouteIndex) (int) {
	var matchIndex []int
	var next = func() {
		for k, v := range matchedRouteIndex {
			if len(v.children) > 0 {
				nextMatchedRouteIndex(index,v.children)
			}
		}
	}

}

func matchRouter(method string, requestPath string, routes []*route) []*matchedRouteIndex {
	var matchedRoutes []*matchedRouteIndex
	//fmt.Printf("___1____\n")
	for matchIndex, route := range routes {
		//fmt.Printf("____2___\n")
		//fmt.Printf("_______method , route.method ,\n", route.method, method, route.method != "ALL" && route.method != method)
		if route.method != "ALL" && route.method != method {
			continue
		}
		//fmt.Printf("___3____\n")

		if !strings.HasPrefix(requestPath, "/") {
			requestPath = "/" + requestPath
		}
		//fmt.Printf("____ route.routePath, requestPath, route.matchType\n\n", route.routePath, requestPath, route.matchType)
		params, matched := matchPath(route.routePath, requestPath, route.matchType)

		if matched {
			//fmt.Printf("___4____\n")
			if params == nil {
				params = make(map[string]string)
			}
			//fmt.Printf("____5___\n")
			matchedRouteIndex := &matchedRouteIndex{index: matchIndex, params: params}
			if len(route.children) > 0 {
				//fmt.Printf("____6___\n")
				requestPaths := strings.Split(requestPath, "/")
				if len(requestPaths) > 1 {
					//fmt.Printf("___7____\n")
					subrequestPath := strings.Join(strings.Split(requestPath, "/")[2:], "/")
					//fmt.Printf("==1 route.method, subrequestPath, route.children \n\n", route.method, subrequestPath, route.children)
					//fmt.Printf("\nroute.children %s\n", route.children)
					matchedRouteIndex.children = matchRouter(method, subrequestPath, route.children)
					//fmt.Printf("matchedRouteIndex.children %s\n", matchedRouteIndex.children)
				} else {
					//fmt.Printf("___8____\n")
					matchedRouteIndex.children = matchRouter(route.method, "/", route.children)
				}
			}
			//fmt.Printf("___9____\n")
			matchedRoutes = append(matchedRoutes, matchedRouteIndex)
		}
		//fmt.Printf("____10___\n")
	}
	//fmt.Printf("____11___\n")

	return matchedRoutes
}

func doHandler(req *Req, res *Res, matchedRouteIndex *matchedRouteIndex, routes []*route) bool {
	//routeIndex := matchedRoutes[routeIndex]
	if res.exit {
		return true
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("r %s\n", r)
			fmt.Printf("routes[matchedRouteIndex.index] %s %s %s\n", routes, len(routes), matchedRouteIndex.index)
		}
	}()
	matchedRoute := routes[matchedRouteIndex.index]
	if matchedRoute.handlerFunc != nil {
		matchedRoute.handlerFunc(req, res)
		return true
	} else if matchedRoute.handlerFuncNext != nil {
		matchedRoute.handlerFuncNext(req, res, func() {
			for i, _ := range routes[matchedRouteIndex.index+1:] {
				doHandler(req, res, i, matchedRoutes, routes)
			}
		})
		return true
	} else {
		//fmt.Printf("=-=-=-=-=-=-=-=-==-\n", matchedRoute.children)
		return false
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

	//fmt.Printf("%s\n", g.routes)
	//fmt.Printf("%s\n", g.routes[1].children)
	matchedRoutes := matchRouter(r.Method, strings.Split(r.URL.Path, "?")[0], g.routes)
	//fmt.Printf("matchedRoutes len %d\n", len(matchedRoutes))
	//fmt.Printf("~~~~~~~1~~~~~~\n")
	for _, v := range matchedRoutes {
		//fmt.Printf("~~~~~~~2~~~~~~\n")
		//fmt.Printf("matchedRoutes v %s\n", v)
		req.Params = v.params
		//fmt.Printf("routes length %d \n", len(g.routes))

		doHandler(req, res, v, g.routes)
		//fmt.Printf("~~~~~~~g.doHandler(req, res, i, matchedRoutes)~~~~~~\n", g.doHandler, req, res, i, matchedRoutes)

		//fmt.Printf("v.children len %d\n", len(v.children))
		//fmt.Printf("~~~~~~~3~~~~~~\n")
		for _, v2 := range v.children {
			//fmt.Printf("~~~~~~~4~~~~~~\n")
			req.Params = v2.params
			//fmt.Printf("routes length %d \n", len(g.routes[v.index].children))
			doHandler(req, res, v2, g.routes[v.index].children)
			//fmt.Printf("~~~~~~~g.doHandler(req, res, i2, matchedRoutes)~~~~~~\n", g.doHandler, req, res, i2, v.children)
			//fmt.Printf("~~~~~~5~~~~~~~\n")
		}
		//fmt.Printf("~~~~~~~6~~~~~~\n")
	}
	//fmt.Printf("~~~~~~~7~~~~~~\n")

	res.SendStatus(http.StatusNotFound)
	//fmt.Printf("~~~~~~~8~~~~~~\n")
}

// Listen bind port and start server
func (g *Gor) Listen(addr string) error {
	return http.ListenAndServe(addr, g)
}
