package gor

import (
	"net/http"
	"strings"
)

func (g *Gor) matchRouter(w http.ResponseWriter, r *http.Request, req *Req, res *Res) *HandlerFunc {
	routes := strings.Split(strings.Split(r.URL.Path[1:], "?")[0], "/")

	for _, route := range g.routes {
		if route.method == r.Method && route.prepath == routes[0] {
			matchRoutes := routes[1:]
			if len(matchRoutes) != len(route.routeParams) {
				continue
			}
			if len(route.routeParams) == 0 && len(matchRoutes) == 0 {
				return &route.handler
			}
			match := false
			for i, j := 0, len(matchRoutes); i < j; i++ {
				if route.routeParams[i].isParam {
					req.Params[route.routeParams[i].name] = matchRoutes[i]
				} else if route.routeParams[i].name != matchRoutes[i] {
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
	g.middlewares = append(g.middlewares, middlewares...)
}

// Listen bind port and start server
func (g *Gor) Listen(addr string) error {
	return http.ListenAndServe(addr, g)
}
