package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Chyroc/gor"
)

func Logger(req *gor.Req, res *gor.Res, next gor.Next) {
	startTime := time.Now()
	next()
	fmt.Printf("[LOG] method: %s, time: %s, response %s\n", req.Method, time.Now().Sub(startTime), res.Response)
}

func Print(req *gor.Req, res *gor.Res, next gor.Next) {
	fmt.Printf("before\n")
	next()
	fmt.Printf("after\n")
}

func main() {
	app := gor.NewGor()
	router := gor.NewRouter()

	app.Use(Logger)
	app.Use(Print)
	router.Get("/sub", func(req *gor.Req, res *gor.Res) {
		fmt.Printf("content... \n")
		res.Send("this is data.")
	})

	app.Use("/main", router)

	log.Fatal(app.Listen(":3000"))
}
