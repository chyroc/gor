package main

import (
	"fmt"
	"github.com/Chyroc/gor"
	"net/http"
)

func Logger(g *gor.Gor) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("logger")
	}

	return http.HandlerFunc(fn)
}

func main() {
	app := gor.NewGor()
	router := gor.NewRouter()

	router.Use(func(req *gor.Req, res gor.Res) {
		res.Send()
	})

	app.Use(Logger)
	app.Get("/", func(req *gor.Req, res gor.Res) {
		res.Send("1")
	})
	app.Use(Logger)

	app.Listen(":3000")
}
