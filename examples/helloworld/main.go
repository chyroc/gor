package main

import (
	"log"

	"github.com/Chyroc/gor"
	"github.com/Chyroc/gor/middlerware"
)

func main() {
	app := gor.NewGor()
	app.Use(middlerware.Logger)
	app.SetRenderDir("testdata/examples/helloword")
	app.Static("./vendor")

	app.Get("/", func(req *gor.Req, res *gor.Res) {
		res.Send("Hello World")
	})

	app.Get("/json", func(req *gor.Req, res *gor.Res) {
		res.JSON([]string{"a", "b", "c"})
	})

	app.Get("/redirect", func(req *gor.Req, res *gor.Res) {
		res.Redirect("/404")
	})

	app.Get("/500", func(req *gor.Req, res *gor.Res) {
		res.Status(500).End()
	})

	router := gor.NewRouter()
	router.Get("/1", func(req *gor.Req, res *gor.Res) { res.HTML("1", "") })
	router.Get("/2", func(req *gor.Req, res *gor.Res) { res.HTML("2", map[string]string{"name": "Chyroc"}) })
	// for more, to see: https://golang.org/pkg/text/template
	app.All("/html", router)

	log.Fatal(app.Listen(":3000"))
}
