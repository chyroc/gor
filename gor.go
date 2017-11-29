package gor

import (
	"log"
	"strings"
	"fmt"
	"regexp"
)

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

var debug = true

func debugPrintf(format string, a ...interface{}) {
	if debug {
		log.Printf(format+"\n", a...)
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
				regS = append(regS, `(?P<`+strings.ToLower(v[1:])+`>.*[^/])`)
			} else {
				regS = append(regS, v)
			}
		}
		return regexp.MustCompile(strings.Join(regS, "/"))
	}

	return nil
}

func matchPath(routePath, requestPath string, matchtype matchType) (map[string]string, bool) {
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

	switch matchtype {
	case fullMatch:
		if containsColon {
			if len(strings.Split(routePath, "/")) != len(strings.Split(requestPath, "/")) {
				return nil, false
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
			return nil, false
		} else {
			return nil, routePath == requestPath
		}
	case preMatch:
		if containsColon {
			if len(strings.Split(routePath, "/")) > len(strings.Split(requestPath, "/")) {
				return nil, false
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
			return nil, false
		} else {
			return nil, strings.HasPrefix(requestPath, routePath) && (len(requestPath) == len(routePath) || strings.HasPrefix(requestPath[len(routePath):], "/"))
		}
	default:
		fmt.Errorf("matchtype must be one of fullMatch or preMatch")
	}
	return nil, false
}
