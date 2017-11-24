package gor

import (
	"context"
	"net/http"
)

// Req is http Request struct
type Req struct {
	httpr   *http.Request
	context context.Context
}

func httpRequestToReq(httpRequest *http.Request) *Req {
	return &Req{
		httpRequest,
		httpRequest.Context(),
	}
}

// AddContext add value to gor context
func (req *Req) AddContext(key, val interface{}) {
	req.context = context.WithValue(req.context, key, val)
}

// GetContext get context from gor by key
func (req *Req) GetContext(key interface{}) interface{} {
	return req.context.Value(key)
}
