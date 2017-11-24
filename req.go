package gor

import (
	"net/http"
	"context"
)

// Req is http Request struct
type Req struct {
	httpr *http.Request
}

func httpRequestToReq(httpRequest *http.Request) *Req {
	return &Req{
		httpRequest,
	}
}

func (req *Req) AddContext(key, val interface{}) *Req {
	newContext := context.WithValue(req.httpr.Context(), key, val)
	req = httpRequestToReq(req.httpr.WithContext(newContext))
	return req
}
