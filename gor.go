package gor

import (
	"net/http"
	"regexp"
	"strings"
)

type matchedRouteArray []*route

// Gor gor framework core struct
type Gor struct {
	*Route
}

// NewGor return Gor struct
func NewGor() *Gor {
	return &Gor{
		NewRoute(),
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
			route.handlerFuncNext(req, res, func(errs ...string) {
				if len(errs) > 0 {
					res.Error(strings.Join(errs, ", "))
					return
				}
				noCallNext = false
				doHandler(req, res, index+1, matchRoutes, requestPath)
			})
			if noCallNext {
				res.exit = true
				return
			}
		} else {
			panic("This can not exist when handler the request, this is a bug, please report : https://github.com/Chyroc/gor/issues")
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

func recursionMatch(method, requestPath, prePath string, parentRoutes []*route, matchedRoutes *matchedRouteArray) {
	//fmt.Printf("\n\nrecursionMatch %s %s %s %s %s\n", method, requestPath, prePath, parentRoutes, matchedRoutes)
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
		mType := route.matchType
		if mType == onlyLastFull {
			if len(route.children) > 0 {
				mType = preMatch
			} else {
				mType = fullMatch
			}
		}

		_, matched := matchPath(route.routePath, requestPath, mType)
		//fmt.Printf("matched %s\n", matched, route.routePath, requestPath, route.matchType)

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

func genMatchPathReg(routePath string) *regexp.Regexp {
	if strings.HasSuffix(routePath, "/") {
		routePath = routePath[:len(routePath)-1]
	}

	if strings.ContainsRune(routePath, ':') {
		routePaths := strings.Split(routePath, "/")
		var regS []string
		for _, v := range routePaths {
			if strings.HasPrefix(v, ":") {
				regS = append(regS, `(?P<`+strings.ToLower(v[1:])+`>.*)`)
			} else {
				regS = append(regS, v)
			}
		}
		return regexp.MustCompile(strings.Join(regS, "/"))
	}

	return nil
}

// matchtype pre full onlyLastFull
func matchPath(routePath, requestPath string, matchtype matchType) (params map[string]string, matched bool) {
	if !strings.HasPrefix(routePath, "/") {
		routePath = "/" + routePath
	}
	if strings.HasSuffix(routePath, "/") {
		routePath = routePath[:len(routePath)-1]
	}

	if strings.HasSuffix(requestPath, "/") {
		requestPath = requestPath[:len(requestPath)-1]
	}

	containsColon := strings.ContainsRune(routePath, ':')
	var reg *regexp.Regexp
	if containsColon {
		reg = genMatchPathReg(routePath)
	}

	params = make(map[string]string)
	matched = false

	switch matchtype {
	case fullMatch:
		if containsColon {
			if len(strings.Split(routePath, "/")) != len(strings.Split(requestPath, "/")) {
				return
			}
			paramsMap := make(map[string]string)
			match := reg.FindStringSubmatch(requestPath)
			for i, name := range reg.SubexpNames() {
				if i > 0 && i <= len(match) {
					paramsMap[name] = match[i]
				}
			}
			if len(match) > 0 {
				return paramsMap, true
			}
			return
		}
		return make(map[string]string), routePath == requestPath
	case preMatch:
		if containsColon {
			if len(strings.Split(routePath, "/")) > len(strings.Split(requestPath, "/")) {
				return
			}

			paramsMap := make(map[string]string)
			match := reg.FindStringSubmatch(strings.Join(strings.Split(requestPath, "/")[:len(strings.Split(routePath, "/"))], "/"))
			for i, name := range reg.SubexpNames() {
				if i > 0 && i <= len(match) {
					paramsMap[name] = match[i]
				}
			}

			if len(match) > 0 {
				return paramsMap, true
			}
			return
		}
		return make(map[string]string), strings.HasPrefix(requestPath, routePath) && (len(requestPath) == len(routePath) || strings.HasPrefix(requestPath[len(routePath):], "/"))
	}
	return
}
