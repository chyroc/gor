package gor

import (
	"net/http"
)

// Gor gor framework core struct
type Gor struct {
	handlers         map[string]HandlerFunc
	handlersWithNext map[string]HandlerFuncWithNext
	ttt              map[string][]HandlerFunc
	middlewares      []func(g *Gor) http.Handler
	midWithPath      map[string]int
}

// NewGor return Gor struct
func NewGor() *Gor {
	return &Gor{
		handlers:         make(map[string]HandlerFunc),
		handlersWithNext: make(map[string]HandlerFuncWithNext),
		ttt:              make(map[string][]HandlerFunc),
		midWithPath:      make(map[string]int),
	}
}
