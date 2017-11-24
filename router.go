package gor

// Gor gor framework core struct
type Gor struct {
	handlers         map[string]HandlerFunc
	handlersWithNext map[string]HandlerFuncWithNext
	ttt              map[string][]HandlerFuncWithNext
}

// NewGor return Gor struct
func NewGor() *Gor {
	return &Gor{
		handlers:         make(map[string]HandlerFunc),
		handlersWithNext: make(map[string]HandlerFuncWithNext),
		ttt:              make(map[string][]HandlerFuncWithNext),
	}
}
