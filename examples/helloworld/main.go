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

	app.Get("/json", func(req *gor.Req, res gor.Res) {
		res.JSON([]string{"a", "b", "c"})
	})

	fmt.Printf("err: %+v\n", app.Listen(":3000"))
}
