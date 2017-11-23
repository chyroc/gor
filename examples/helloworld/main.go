package main

import (
	"fmt"

	"github.com/Chyroc/gor"
)

func main() {
	app := gor.NewGor()

	app.Get("/", func(req *gor.Req, res gor.Res) {
		res.Send("Hello World")
	})

	app.Post("/", func(req *gor.Req, res gor.Res) {
		res.Send("Hello World : POST")
	})

	fmt.Printf("err: %+v\n", app.Listen(":3000"))
}
