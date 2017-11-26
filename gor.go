package gor

// HandlerFunc gor handler func like http.HandlerFunc func(ResponseWriter, *Request)
type HandlerFunc func(*Req, *Res)

// Gor gor framework core struct
type Gor struct {
	*Route
}

// NewGor return Gor struct
func NewGor() *Gor {
	return &Gor{
		NewRoute(),
	}
}
