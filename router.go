package gor

// Router router
type Router struct {
	middlewares []HandlerFunc
}

// NewRouter return *Router
func NewRouter() *Router {
	return &Router{}
}

// All add a handler for all HTTP verbs to this route
func (r *Router) All() {
}

// Method method
func (r *Router) Method() {
}

// Param map the given param placeholder `name`(s) to the given callback.
func (r *Router) Param() {
}

// Route Create a new Route for the given path.
func (r *Router) Route() {
}

// Use use the given middleware function, with optional path, defaulting to "/"
func (r *Router) Use(middlewares ...HandlerFunc) {
	r.middlewares = append(r.middlewares, middlewares...)
}
