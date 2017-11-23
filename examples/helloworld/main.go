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

	app.Get("/json", func(req *gor.Req, res gor.Res) {
		res.JSON(1)
	})

	fmt.Printf("err: %+v\n", app.Listen(":3000"))
}
