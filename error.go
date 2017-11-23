package gor

import "errors"

var (
	// ErrNotFound is not found error.
	ErrNotFound = errors.New("not found page")
	// ErrResponseTypeUnsupported is response type unsupported error.
	ErrResponseTypeUnsupported = errors.New("response type unsupported")
	// ErrJsonMarshal is json marshal error.
	ErrJsonMarshal = errors.New("json marshal err")
)
