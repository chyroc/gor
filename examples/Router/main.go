package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Chyroc/gor"
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

	router.Get("/sub/:uu", func(req *gor.Req, res *gor.Res) {
		res.JSON(map[string]interface{}{
			"query":  req.Query,
			"params": req.Params,
		})
	})

	app.UseN("/user", router)

	log.Fatal(app.Listen(":3000"))
}
