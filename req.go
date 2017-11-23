package gor

import "net/http"

type Req struct {
	*http.Request
}

func httpRequestToReq(httpRequest *http.Request) *Req {
	return &Req{
		httpRequest,
	}
}
