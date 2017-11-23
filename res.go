package gor

import "net/http"

type Res struct {
}

//ResponseWriter, *Request

func httpResponseWriterToRes(r http.ResponseWriter) Res {
	return Res{}
}

func (req *Req) Header() http.Header {
	return make(http.Header)
}

func (req *Req) Write([]byte) (int, error) {
	return 0, nil
}

func (req *Req) WriteHeader(int) {
}

func (res *Res) Send(v interface{}) {

}
