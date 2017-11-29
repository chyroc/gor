package gor

import (
	"github.com/stretchr/testify/assert"
)

type handlerType int

const (
	unkonwFunc handlerType = iota
	handlerFunc
	handlerFuncNext
	midFunc
)

func renderParamQuery(req *Req, res *Res) {
	res.JSON(map[string]interface{}{
		"params": req.Params,
		"query":  req.Query,
	})
}

func assertBetweenRoute(as *assert.Assertions, expected, actual *route) {
	assertRoute(as, expected.method, expected.routePath, expected.matchType, unkonwFunc, actual)
	if expected.handlerFunc == nil {
		as.Nil(actual.handlerFunc)
	} else {
		as.NotNil(actual.handlerFunc)
	}

	if expected.handlerFuncNext == nil {
		as.Nil(actual.handlerFuncNext)
	} else {
		as.NotNil(actual.handlerFuncNext)
	}

	if expected.middleware == nil {
		as.Nil(actual.middleware)
	} else {
		as.NotNil(actual.middleware)
	}
}

func assertOneRoute(as *assert.Assertions, method, routePath string, matchType matchType, handlerType handlerType, actuals []*route) {
	as.Len(actuals, 1)
	assertRoute(as, method, routePath, matchType, handlerType, actuals[0])
}

func assertRoute(as *assert.Assertions, method, routePath string, matchType matchType, handlerType handlerType, actual *route) {
	if method != "ALL" {
		as.Equal(method, actual.method)
	}
	as.Equal(routePath, actual.routePath)
	as.Equal(matchType, actual.matchType)

	switch handlerType {
	case handlerFunc:
		as.NotNil(actual.handlerFunc)
		as.Nil(actual.handlerFuncNext)
		as.Nil(actual.middleware)
	case handlerFuncNext:
		as.Nil(actual.handlerFunc)
		as.NotNil(actual.handlerFuncNext)
		as.Nil(actual.middleware)
	case midFunc:
		as.Nil(actual.handlerFunc)
		as.Nil(actual.handlerFuncNext)
		as.NotNil(actual.middleware)
	}
}
