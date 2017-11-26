package gor

// Router router
type Router struct {
	*Route
}

// NewRouter return *Router
func NewRouter() *Router {
	return &Router{
		NewRoute(),
	}
}
