package gor

import (
	"context"
	"net/http"
	"net/url"
)

// Req is http Request struct
type Req struct {
	r       *http.Request
	context context.Context

	Method string
	Query  map[string][]string
}

func getQuery(r *http.Request) (map[string][]string, error) {
	URL, err := url.Parse(r.URL.Path)
	if err != nil {
		return nil, err
	}
	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}
	return query, nil
}

func httpRequestToReq(r *http.Request) (*Req, error) {
	query, err := getQuery(r)
	if err != nil {
		return nil, err
	}

	return &Req{
		r:       r,
		context: r.Context(),

		Method: r.Method,
		Query:  query,
	}, nil
}

// AddContext add value to gor context
func (req *Req) AddContext(key, val interface{}) {
	req.context = context.WithValue(req.context, key, val)
}

// GetContext get context from gor by key
func (req *Req) GetContext(key interface{}) interface{} {
	return req.context.Value(key)
}
