package main

import (
	"log"
	"time"

	"github.com/Chyroc/gor"
)

func Logger(req *gor.Req, res *gor.Res) gor.HandlerFunc {
	req.AddContext("time", time.Now())
	return func(req *gor.Req, res *gor.Res) {
		startTime := req.GetContext("time").(time.Time)
		log.Printf("startTime %+v , endTime %+v\n", startTime.UTC(), time.Now().UTC())
	}
}

func main() {
	app := gor.NewGor()
	router := gor.NewRouter()

	router.Use(Logger)
	router.Get("/sub/:uu", func(req *gor.Req, res *gor.Res) {
		time.Sleep(time.Microsecond * 100)
		res.JSON(map[string]interface{}{
			"query":  req.Query,
			"params": req.Params,
		})
	})

	app.UseN("/user", router)

	log.Fatal(app.Listen(":3000"))
}
