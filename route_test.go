package gor

import (
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

func TestMid(t *testing.T) {
	{
		// change req
		app, ts, e, _ := newTestServer(t)
		defer ts.Close()

		app.Use(func(req *Req, res *Res) HandlerFunc {
			req.AddContext("test", "test-value")
			return nil
		})
		app.Get("/", func(req *Req, res *Res) { res.Send(req.GetContext("test")) })
		e.GET("/").Expect().Status(http.StatusOK).Text().Equal("test-value")
	}
	{
		// get response
		var r interface{}
		app, ts, e, as := newTestServer(t)
		defer ts.Close()

		app.Use(func(req *Req, res *Res) HandlerFunc {
			return func(req *Req, res *Res) {
				r = res.Response
			}
		})
		app.Get("/", func(req *Req, res *Res) { res.Send("this is response") })
		e.GET("/").Expect().Status(http.StatusOK)
		as.Equal("this is response", r)
	}
	{
		// todo exec before response
		// need ?
	}
}
