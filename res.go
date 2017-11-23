package gor

import "net/http"

type Res struct {
	http.ResponseWriter
}

func httpResponseWriterToRes(httpResponseWriter http.ResponseWriter) Res {
	return Res{
		httpResponseWriter,
	}
}

func (res *Res) Send(v string) {
	res.Write([]byte(v))
}
