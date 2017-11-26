package gor

import (
	"net/http"
	"strings"
	"fmt"
)

// splitRoute split url to path array
// / -> [""]
// /a -> [a]
func splitRoute(r *http.Request) []string {
	return strings.Split(strings.Split(r.URL.Path[1:], "?")[0], "/")
	//path := strings.Split(r.URL.Path, "?")[0]
	//if path == "/" {
	//	return []string{"/"}
	//}
	//paths := strings.Split(path, "/")
	//paths[0] = "/"
	//return paths
}

func (g *Gor) matchRouter(w http.ResponseWriter, r *http.Request, req *Req, res *Res) *HandlerFunc {
	fmt.Printf("routes %s\n", g.routes)
	fmt.Printf("routeParams %s\n", g.routes[0].routeParams)
	fmt.Printf("prepath %s\n", g.routes[0].prepath)

	routes := splitRoute(r)
	fmt.Printf("routes %+v %s\n", routes, len(routes))

	for _, route := range g.routes {
		fmt.Printf("==1\n")
		fmt.Printf("", route.method, r.Method, route.prepath, routes[0])
		if route.method == r.Method && (route.prepath == "" || route.prepath == routes[0]) {
			fmt.Printf("==2\n")
			matchRoutes := []string{}
			if route.prepath == "" {
				matchRoutes = routes
			} else {
				matchRoutes = routes[1:]
			}
			//fmt.Printf("", matchRoutes, route.routeParams[0], route.routeParams[1])
			if len(matchRoutes) != len(route.routeParams) {
				continue
			}
			fmt.Printf("==3\n")
			if len(route.routeParams) == 0 && len(matchRoutes) == 0 {
				return &route.handler
			}
			fmt.Printf("==4\n")
			match := false
			for i, j := 0, len(matchRoutes); i < j; i++ {
				if route.routeParams[i].isParam {
					req.Params[route.routeParams[i].name] = matchRoutes[i]
				} else if route.routeParams[i].name != matchRoutes[i] {
					match = false
					break
				}
				match = true
			}
			if match {
				return &route.handler
			}
		}
	}

	return nil
}

// ServeHTTP use to start server
func (g *Gor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := httpResponseWriterToRes(w)
	req, err := httpRequestToReq(r)
	if err != nil {
		res.Error(err.Error())
		return
	}

	if handler := g.matchRouter(w, r, req, res); handler != nil {
		(*handler)(req, res)
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
