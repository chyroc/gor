# gor
Fast, minimalist web framework for [Golang](https://golang.org/)

[![CircleCI](https://circleci.com/gh/Chyroc/gor/tree/master.svg?style=svg&circle-token=5cf109814e08b0d6eee1b4ba4a6e8b2a5c792c84)](https://circleci.com/gh/Chyroc/gor/tree/master)

```go
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

	fmt.Printf("err: %+v\n", app.Listen(":3000"))
}
```
