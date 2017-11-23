package gor

import (
	"testing"
)

func TestNewRouter(t *testing.T) {
	app := NewGor()

	app.Get("/", func(req Req, res Res) {
		res.Send("Hello World")
	})

	//app.Listen(":3000")
}
