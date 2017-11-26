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

	router.Get("/sub", func(req *gor.Req, res *gor.Res) {
		res.JSON(req.Params)
	})

	app.UseN("/user", router)

	app.Listen(":3000")
}
