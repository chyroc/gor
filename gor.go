package gor

import (
	"log"
)

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

var debug = true

func debugPrintf(format string, a ...interface{}) {
	if debug {
		log.Printf(format+"\n", a...)
	}
}
