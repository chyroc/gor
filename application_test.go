package gor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"
)

func newTestServer(t *testing.T) (*Gor, *httptest.Server, *httpexpect.Expect, *assert.Assertions) {
	handler := NewGor()
	ts := httptest.NewServer(handler)
	e := httpexpect.New(t, ts.URL)
	as := assert.New(t)

	return handler, ts, e, as
}

func TestMethod(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Get("/", func(req *Req, res *Res) { res.Send(req.Method) })
	app.Head("/", func(req *Req, res *Res) { res.Send(req.Method) })
	app.Post("/", func(req *Req, res *Res) { res.Send(req.Method) })
	app.Put("/", func(req *Req, res *Res) { res.Send(req.Method) })
	app.Patch("/", func(req *Req, res *Res) { res.Send(req.Method) })
	app.Delete("/", func(req *Req, res *Res) { res.Send(req.Method) })
	app.Connect("/", func(req *Req, res *Res) { res.Send(req.Method) })
	app.Options("/", func(req *Req, res *Res) { res.Send(req.Method) })
	app.Trace("/", func(req *Req, res *Res) { res.Send(req.Method) })

	e.GET("/").Expect().Status(http.StatusOK).Text().Equal("GET")
	e.HEAD("/").Expect().Status(http.StatusOK).Text().Equal("") // todo ?
	e.POST("/").Expect().Status(http.StatusOK).Text().Equal("POST")
	e.PUT("/").Expect().Status(http.StatusOK).Text().Equal("PUT")
	e.PATCH("/").Expect().Status(http.StatusOK).Text().Equal("PATCH")
	e.DELETE("/").Expect().Status(http.StatusOK).Text().Equal("DELETE")
	e.Request("CONNECT", "/").Expect().Status(http.StatusOK).Text().Equal("CONNECT")
	e.OPTIONS("/").Expect().Status(http.StatusOK).Text().Equal("OPTIONS")
	e.Request("TRACE", "/").Expect().Status(http.StatusOK).Text().Equal("TRACE")
}

func testMid(g *Gor) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%+v\n", g.midWithPath)
	})
}

func TestNext(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Use(testMid)
	app.Get("/", func(req *Req, res *Res) { res.Send("11") })
	app.Use(testMid)
	app.Post("/2", func(req *Req, res *Res) { res.Send("22") })

	e.GET("/").Expect().Status(http.StatusOK).Text().Equal("11")
	e.POST("/2").Expect().Status(http.StatusOK).Text().Equal("22")

	// todo
	//as.Len(app.middlewares, 2)
	//as.Equal(map[string]int{"GET/": 0, "POST/2": 1}, app.midWithPath)
}
