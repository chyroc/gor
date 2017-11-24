package gor

import "errors"

var (
	// ErrNotFound is not found error.
	ErrNotFound = errors.New("not found page")
	// ErrResponseTypeUnsupported is response type unsupported error.
	ErrResponseTypeUnsupported = errors.New("response type unsupported")
	// ErrJSONMarshal is json marshal error.
	ErrJSONMarshal = errors.New("json marshal err")
	// ErrHTTPStatusCodeInvalid is given http status code is invalid error.
	ErrHTTPStatusCodeInvalid = errors.New("http status code is invalid")
)
