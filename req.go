package gor

import "net/http"

// Req is http Request struct
type Req struct {
	*http.Request
}

func httpRequestToReq(httpRequest *http.Request) *Req {
	return &Req{
		httpRequest,
	}
}
