package gor

// Gor gor framework core struct
type Gor struct {
	*Route
	renderDir string
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
