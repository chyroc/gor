package gor

import "net/http"

//http.Request

type Req struct {
}

func httpRequestToReq(r *http.Request) *Req {
	return &Req{

	}
}
