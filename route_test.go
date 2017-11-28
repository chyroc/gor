package gor

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"

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

func TestUse(t *testing.T) {
	{
		// only app use
		app, ts, _, as := newTestServer(t)
		defer ts.Close()

		func1 := func(req *Req, res *Res) {}
		app.Use("/1", func1)
		as.Len(app.routes, 1)
		as.Equal("ALL", app.routes[0].method)
		as.Equal("1", app.routes[0].prepath)
		//as.Equal("/1", app.routes[0].parentIndex)
		as.Len(app.routes[0].routeParams, 0)
		as.Nil(app.routes[0].handlerFuncNext)
		as.Nil(app.routes[0].middleware)

		app.Use("/2", func(req *Req, res *Res, next Next) {})
		as.Len(app.routes, 2)
		as.Equal("ALL", app.routes[1].method)
		as.Equal("2", app.routes[1].prepath)
		//as.Equal("/2", app.routes[1].parentIndex)
		as.Len(app.routes[1].routeParams, 0)
		as.Nil(app.routes[1].handlerFunc)
		as.Nil(app.routes[1].middleware)
	}
	{
		// app use and router
		app, ts, e, as := newTestServer(t)
		defer ts.Close()
		router := NewRouter()
		app.Get("/no-sub", func(req *Req, res *Res) {})
		router.Get("/1", func(req *Req, res *Res) {})
		router.Get("/2", func(req *Req, res *Res) {})
		app.Use("/main", router)

		as.Len(app.routes, 3)
		as.Equal("GET", app.routes[0].method)
		as.Equal("no-sub", app.routes[0].prepath)
		//as.Equal("/no-sub", app.routes[0].parentIndex)
		as.Len(app.routes[0].routeParams, 0)
		as.Nil(app.routes[0].handlerFuncNext)
		as.Nil(app.routes[0].middleware)

		as.Equal("GET", app.routes[1].method)
		as.Equal("main", app.routes[1].prepath)
		//as.Equal("/main/1", app.routes[1].parentIndex)
		as.Len(app.routes[1].routeParams, 1)
		as.Equal(&routeParam{name: "1", isParam: false}, app.routes[1].routeParams[0])
		as.Nil(app.routes[1].handlerFuncNext)
		as.Nil(app.routes[1].middleware)

		as.Equal("GET", app.routes[2].method)
		as.Equal("main", app.routes[2].prepath)
		//as.Equal("/main/2", app.routes[2].parentIndex)
		as.Len(app.routes[2].routeParams, 1)
		as.Equal(&routeParam{name: "2", isParam: false}, app.routes[2].routeParams[0])
		as.Nil(app.routes[2].handlerFuncNext)
		as.Nil(app.routes[2].middleware)

		e.GET("/no-sub").Expect().Status(http.StatusOK).Text().Equal("")
		e.GET("/main/1").Expect().Status(http.StatusOK).Text().Equal("")
		e.GET("/main/2").Expect().Status(http.StatusOK).Text().Equal("")
	}
	{
		// app use and router + params
		app, ts, e, as := newTestServer(t)
		defer ts.Close()
		router := NewRouter()
		app.Get("/no-sub/:name0", func(req *Req, res *Res) { res.JSON(req.Params) })
		router.Get("/1/:name1", func(req *Req, res *Res) { res.JSON(req.Params) })
		router.Get("/2/:name2", func(req *Req, res *Res) { res.JSON(req.Params) })
		app.Use("/main", router)

		as.Len(app.routes, 3)
		as.Equal("GET", app.routes[0].method)
		as.Equal("no-sub", app.routes[0].prepath)
		//as.Equal("/no-sub", app.routes[0].parentIndex)
		as.Len(app.routes[0].routeParams, 1)
		as.Equal(&routeParam{name: "name0", isParam: true}, app.routes[0].routeParams[0])
		as.Nil(app.routes[0].handlerFuncNext)
		as.Nil(app.routes[0].middleware)

		as.Equal("GET", app.routes[1].method)
		as.Equal("main", app.routes[1].prepath)
		//as.Equal("/main/1", app.routes[1].parentIndex)
		as.Len(app.routes[1].routeParams, 2)
		fmt.Printf("app.routes[1].routeParams %s\n", app.routes[1].routeParams)
		as.Equal(&routeParam{name: "1", isParam: false}, app.routes[1].routeParams[0])
		as.Equal(&routeParam{name: "name1", isParam: true}, app.routes[1].routeParams[1])
		as.Nil(app.routes[1].handlerFuncNext)
		as.Nil(app.routes[1].middleware)

		as.Equal("GET", app.routes[2].method)
		as.Equal("main", app.routes[2].prepath)
		//as.Equal("/main/2", app.routes[2].parentIndex)
		as.Len(app.routes[2].routeParams, 2)
		as.Equal(&routeParam{name: "2", isParam: false}, app.routes[2].routeParams[0])
		as.Equal(&routeParam{name: "name2", isParam: true}, app.routes[2].routeParams[1])
		as.Nil(app.routes[2].handlerFuncNext)
		as.Nil(app.routes[2].middleware)

		e.GET("/no-sub/nnn").Expect().Status(http.StatusOK).JSON().Equal(map[string]string{"name0": "nnn"})
		e.GET("/main/1/mmm").Expect().Status(http.StatusOK).JSON().Equal(map[string]string{"name1": "mmm"})
		e.GET("/main/2/jjj").Expect().Status(http.StatusOK).JSON().Equal(map[string]string{"name2": "jjj"})
	}
}
