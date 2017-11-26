package gor

import "fmt"

// HandlerFunc gor handler func like http.HandlerFunc func(ResponseWriter, *Request)
type HandlerFunc func(*Req, *Res)

// HandlerFuncDefer gor handler func like http.HandlerFunc func(ResponseWriter, *Request),
// but return HandlerFunc to do somrthing at defer time
type HandlerFuncDefer func(*Req, *Res) HandlerFunc

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

var debug = false

func debugPrintf(format string, a ...interface{}) {
	if debug {
		fmt.Printf(format+"\n", a...)
	}
}
