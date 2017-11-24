package main

import (
	"fmt"
	"github.com/Chyroc/gor"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("logger")

		//next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func main() {
	app := gor.NewGor()

	app.Use(Logger)
	app.Get("/", func(req *gor.Req, res gor.Res) {
		res.Send("1")
	})
	app.Use(Logger)

	app.Listen(":3000")
}
