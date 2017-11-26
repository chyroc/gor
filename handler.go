package gor

import (
	"net/http"
	"strings"
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

func matchRouter2(method string, paths []string, routes []*route) (map[string]string, int) {
	for _, v := range paths {
		if strings.Contains(v, "/") {
			panic("paths cannot contain /")
		}
	}
	matchIndex := -1
	for _, route := range routes {
		matchIndex++
		//fmt.Printf("route.prepath %s == paths[0] %s %s\n", route.prepath, paths[0], route.prepath == paths[0])
		if route.prepath == paths[0] {
			if method != "ALL" && route.method != method {
				continue
			}
			//fmt.Printf("11\n")

			matchRoutes := paths[1:]
			//matchRoutes := []string{}
			//if route.prepath == "" {
			//	matchRoutes = paths[1:]
			//} else {
			//	matchRoutes = paths[1:]
			//}
			//fmt.Printf(" %s %s %s %s %s\n",route.prepath, matchRoutes, route.routeParams, len(matchRoutes), len(route.routeParams))
			if len(matchRoutes) != len(route.routeParams) {
				continue
			}
			//fmt.Printf("22\n")

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

func (g *Gor) matchRouter(w http.ResponseWriter, r *http.Request, req *Req, res *Res) *HandlerFunc {
	//fmt.Printf("routes %s\n", g.routes)
	//fmt.Printf("routeParams %s\n", g.routes[0].routeParams)
	//fmt.Printf("prepath %s\n", g.routes[0].prepath)

	routes := splitRoute(r)
	//fmt.Printf("routes %+v %s\n", routes, len(routes))

	//fmt.Printf("r.Method %s, routes %s %s, g.routes %s\n", r.Method, routes, len(routes), g.routes)
	matchParams, matchIndex := matchRouter2(r.Method, routes, g.routes)
	//fmt.Printf("matchParams %s, matchIndex %s\n", matchParams, matchIndex)
	if matchIndex != -1 {
		if matchParams != nil {
			for k, v := range matchParams {
				req.Params[k] = v
			}
		}
		return &g.routes[matchIndex].handler
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
