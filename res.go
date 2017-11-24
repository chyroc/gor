package gor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// Res is http ResponseWriter and some gor Response method
type Res struct {
	http.ResponseWriter
}

func httpResponseWriterToRes(httpResponseWriter http.ResponseWriter) Res {
	return Res{
		httpResponseWriter,
	}
}

func (res *Res) errResponseTypeUnsupported(vtype string, v interface{}) {
	http.Error(res, fmt.Sprintf("[%s] [%s] %+v", ErrResponseTypeUnsupported, vtype, v), http.StatusInternalServerError)
}

// Send Send a response
func (res *Res) Status(code int) *Res {
	res.WriteHeader(code)
	return res
}

// Send Send a response
func (res *Res) Send(v string) {
	res.Write([]byte(v))
}

// JSON Send json response
func (res *Res) JSON(v interface{}) {
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
		fmt.Fprintf(res, ErrJSONMarshal.Error())
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(b)
}

// Redirect Redirect to another url
func (res *Res) Redirect(path string) {
	res.Header().Set("Location", path)
	res.WriteHeader(http.StatusFound)
	res.Write([]byte(fmt.Sprintf(`%s. Redirecting to %s`, http.StatusText(http.StatusFound), path)))
}

// End end the request
func (res *Res) End() {
}
