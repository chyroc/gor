package gor

import "strings"

// Gor gor framework core struct
type Gor struct {
	*Route
	renderDir      string
	staticFilePath string
	staticFielDir  string
}

// NewGor return Gor struct
func NewGor() *Gor {
	return &Gor{
		Route: NewRoute(),
	}
}

// SetRenderDir set rendir tmpl dir
func (g *Gor) SetRenderDir(dir string) {
	g.renderDir = dir
}

// SetStaticPath set static url path start string
func (g *Gor) SetStaticPath(path string) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	g.staticFilePath = path
}

// Static start static file server
func (g *Gor) Static(dir string) {
	g.staticFielDir = dir
}
