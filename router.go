package gor

type Router struct {
	*Route
}

func NewRouter() *Router {
	return &Router{
		NewRoute(),
	}
}
