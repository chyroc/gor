package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Chyroc/gor"
	"time"
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

	router.Use(func(req *gor.Req, res *gor.Res) gor.HandlerFunc {
		fmt.Printf("add mid 1(should exec in real handler)\n")
		req.AddContext("time", time.Now())
		return func(req *gor.Req, res *gor.Res) {
			startTime := req.GetContext("time").(time.Time)
			fmt.Printf("startTime %+v , endTime %+v\n", startTime.UTC(), time.Now().UTC())
			fmt.Printf("response is %+v\n",res.Response)
		}
	})

	router.Get("/sub/:uu", func(req *gor.Req, res *gor.Res) {
		fmt.Printf("add handler\n")
		time.Sleep(time.Microsecond * 100)
		res.JSON(map[string]interface{}{
			"query":  req.Query,
			"params": req.Params,
		})
	})

	app.UseN("/user", router)

	log.Fatal(app.Listen(":3000"))
}
