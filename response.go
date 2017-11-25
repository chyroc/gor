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
	exit bool
}

func httpResponseWriterToRes(httpResponseWriter http.ResponseWriter) *Res {
	return &Res{
		httpResponseWriter,
		false,
	}
}

func (res *Res) errResponseTypeUnsupported(vtype string, v interface{}) {
	http.Error(res, fmt.Sprintf("[%s] [%s] %+v", ErrResponseTypeUnsupported, vtype, v), http.StatusInternalServerError)
}

// Status set response http status code
func (res *Res) Status(code int) *Res {
	if http.StatusText(code) == "" {
		http.Error(res, ErrHTTPStatusCodeInvalid.Error(), http.StatusInternalServerError)
		res.exit = true
	} else {
		res.WriteHeader(code)
	}

	return res
}

// SendStatus set response http status code with its text
func (res *Res) SendStatus(code int) {
	if res.exit {
		return
	}

	res.Status(code).Send(http.StatusText(code))
}

// Send Send a response
func (res *Res) Send(v string) {
	if res.exit {
		return
	}

	res.Write([]byte(v))
	res.exit = true
}

// JSON Send json response
func (res *Res) JSON(v interface{}) {
	defer func() {
		res.exit = true
	}()

	if res.exit {
		return
	}

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

// AddHeader append (key, val) to headers
func (res *Res) AddHeader(key, val string) {
	res.Header().Add(key, val)
}

// SetCookie set cookie
func (res *Res) SetCookie(key, val string, option ...Cookie) {
	var cookie *http.Cookie
	if len(option) > 1 {
		res.Error("only support one cookie option")
	} else if len(option) == 1 {
		cookie = option[0].toHTTPCookie(key, val)
	} else {
		cookie = &http.Cookie{
			Name:  key,
			Value: val,
		}
	}

	http.SetCookie(res, cookie)
}

// Error send erroe response
func (res *Res) Error(v string) {
	http.Error(res, v, http.StatusInternalServerError)
	res.exit = true
}

// End end the request
func (res *Res) End() {
	res.exit = true
}
