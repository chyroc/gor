package gor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

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

func (res *Res) errResponseTypeUnsupported(vtype string, v interface{}) {
	http.Error(res, fmt.Sprintf("[%s] [%s] %+v", ErrResponseTypeUnsupported, vtype, v), http.StatusInternalServerError)
}

func (res *Res) Json(v interface{}) {
	if v == nil {
		res.errResponseTypeUnsupported("nil", v)
		return
	}

	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Map, reflect.Struct, reflect.Slice, reflect.Array:
		break
	default:
		res.errResponseTypeUnsupported(t.Kind().String(), v)
		return
	}

	b, err := json.Marshal(v)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(res, ErrJsonMarshal.Error())
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(b)
}
