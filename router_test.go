package gor

import (
	"net/http"
	"testing"
)

func TestRouterParams(t *testing.T) {
	{
		app, ts, _, as := newTestServer(t)
		defer ts.Close()

		// 1 `/` router
		app.Get("/", func(req *Req, res *Res) { res.End() })
		as.Len(app.routes, 1)
		as.Equal("", app.routes[0].prepath)
		as.Len(app.routes[0].routeParams, 0)

		// 2 `/2` router
		app.Get("/2", func(req *Req, res *Res) { res.End() })
		as.Len(app.routes, 2)
		as.Equal("", app.routes[0].prepath)
		as.Len(app.routes[0].routeParams, 0)
		as.Equal("2", app.routes[1].prepath)
		as.Len(app.routes[1].routeParams, 0)

		// 3 `/3/3` router
		app.Get("/3/3", func(req *Req, res *Res) { res.End() })
		as.Len(app.routes, 3)
		as.Equal("", app.routes[0].prepath)
		as.Len(app.routes[0].routeParams, 0)
		as.Equal("2", app.routes[1].prepath)
		as.Len(app.routes[1].routeParams, 0)
		as.Equal("3", app.routes[2].prepath)
		as.Len(app.routes[2].routeParams, 1)
		as.Equal(&routeParam{name: "3", isParam: false}, app.routes[2].routeParams[0])
	}
	{
		app, ts, _, as := newTestServer(t)
		defer ts.Close()

		// 1 `/1/:user` router
		app.Get("/:user", func(req *Req, res *Res) { res.End() })
		as.Len(app.routes, 1)
		as.Equal("", app.routes[0].prepath)
		as.Len(app.routes[0].routeParams, 1)
		as.Equal(&routeParam{name: "user", isParam: true}, app.routes[0].routeParams[0])
	}
	{
		app, ts, _, as := newTestServer(t)
		defer ts.Close()

		// 1 `/1/:user` router
		app.Get("/1/:user", func(req *Req, res *Res) { res.End() })
		as.Len(app.routes, 1)
		as.Equal("1", app.routes[0].prepath)
		as.Len(app.routes[0].routeParams, 1)
		as.Equal(&routeParam{name: "user", isParam: true}, app.routes[0].routeParams[0])
	}
	{
		app, ts, _, as := newTestServer(t)
		defer ts.Close()

		// 1 `/1/:user` router
		app.Get("/1/:user/not-param/:name", func(req *Req, res *Res) { res.End() })
		as.Len(app.routes, 1)
		as.Equal("1", app.routes[0].prepath)
		as.Len(app.routes[0].routeParams, 3)
		as.Equal(&routeParam{name: "user", isParam: true}, app.routes[0].routeParams[0])
		as.Equal(&routeParam{name: "not-param", isParam: false}, app.routes[0].routeParams[1])
		as.Equal(&routeParam{name: "name", isParam: true}, app.routes[0].routeParams[2])
	}
	{
		_, ts, _, as := newTestServer(t)
		defer ts.Close()

		// 1 `/sub` router with router
		router := NewRouter()
		router.Get("/sub", func(req *Req, res *Res) { res.End() })

		// router
		as.Len(router.routes, 1)
		as.Equal("sub", router.routes[0].prepath)
		as.Len(router.routes[0].routeParams, 0)
	}
}

func TestRouterUse(t *testing.T) {
	// next 2: main router not exist
	{
		app, ts, e, as := newTestServer(t)
		defer ts.Close()

		// `/main/:sub` router with app+router
		router := NewRouter()
		router.Get("/sub", func(req *Req, res *Res) { res.Send("x") })
		app.UseN("/main", router)

		// router
		as.Len(router.routes, 1)
		as.Equal("sub", router.routes[0].prepath)
		as.Len(router.routes[0].routeParams, 0)

		// app
		as.Len(app.routes, 1)
		as.Equal("main", app.routes[0].prepath)
		as.Len(app.routes[0].routeParams, 1)
		as.Equal(&routeParam{name: "sub", isParam: false}, app.routes[0].routeParams[0])

		e.GET("/main/sub?a=1&b=2&b=3").Expect().Status(http.StatusOK).Text().Equal("x")
	}
	{
		app, ts, e, as := newTestServer(t)
		defer ts.Close()

		// `/main/:sub` router with app+router
		router := NewRouter()
		router.Get("/:sub", func(req *Req, res *Res) { res.JSON(req.Params) })
		app.UseN("/main", router)

		// router
		as.Len(router.routes, 1)
		as.Equal("", router.routes[0].prepath)
		as.Len(router.routes[0].routeParams, 1)
		as.Equal(&routeParam{name: "sub", isParam: true}, router.routes[0].routeParams[0])

		// app
		as.Len(app.routes, 1)
		as.Equal("main", app.routes[0].prepath)
		as.Len(app.routes[0].routeParams, 1)
		as.Equal(&routeParam{name: "sub", isParam: true}, app.routes[0].routeParams[0])

		e.GET("/main/sub").Expect().Status(http.StatusOK).JSON().Equal(map[string]string{"sub": "sub"})
	}
	// next 2: main router exist
	{
		app, ts, e, as := newTestServer(t)
		defer ts.Close()

		app.Get("/main/not-sub", func(req *Req, res *Res) { res.Send("no-sub") })
		router := NewRouter()
		router.Get("/sub", func(req *Req, res *Res) { res.Send("sub") })
		app.UseN("/main", router)

		// router
		as.Len(router.routes, 1)
		as.Equal("sub", router.routes[0].prepath)
		as.Len(router.routes[0].routeParams, 0)

		// app
		as.Len(app.routes, 2)
		as.Equal("main", app.routes[0].prepath)
		as.Len(app.routes[0].routeParams, 1)
		as.Equal(&routeParam{name: "not-sub", isParam: false}, app.routes[0].routeParams[0])
		as.Equal("main", app.routes[1].prepath)
		as.Len(app.routes[1].routeParams, 1)
		as.Equal(&routeParam{name: "sub", isParam: false}, app.routes[1].routeParams[0])

		e.GET("/main/not-sub").Expect().Status(http.StatusOK).Text().Equal("no-sub")
		e.GET("/main/sub").Expect().Status(http.StatusOK).Text().Equal("sub")

	}
	{
		app, ts, e, as := newTestServer(t)
		defer ts.Close()

		app.Get("/main", func(req *Req, res *Res) { res.Send("main") })
		router := NewRouter()
		router.Get("/:sub", func(req *Req, res *Res) { res.JSON(req.Params) })
		app.UseN("/main", router)

		// router
		as.Len(router.routes, 1)
		as.Equal("", router.routes[0].prepath)
		as.Len(router.routes[0].routeParams, 1)
		as.Equal(&routeParam{name: "sub", isParam: true}, router.routes[0].routeParams[0])

		// app
		as.Len(app.routes, 2)
		as.Equal("main", app.routes[0].prepath)
		as.Len(app.routes[0].routeParams, 0)
		as.Equal("main", app.routes[1].prepath)
		as.Len(app.routes[1].routeParams, 1)
		as.Equal(&routeParam{name: "sub", isParam: true}, app.routes[1].routeParams[0])

		// todo main
		e.GET("/main").Expect().Status(http.StatusOK).Text().Equal("main")
		e.GET("/main/sub").Expect().Status(http.StatusOK).JSON().Equal(map[string]string{"sub": "sub"})
	}
}

func TestRouterUseWithParamsAndQuery(t *testing.T) {
	// next 2: main router not exist
	{
		app, ts, e, _ := newTestServer(t)
		defer ts.Close()

		router := NewRouter()
		router.Get("/sub", func(req *Req, res *Res) { res.JSON(map[string]interface{}{"query": req.Query, "params": req.Params}) })
		app.UseN("/main", router)

		e.GET("/main/sub?a=1&b=2&b=3").Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"query": map[string][]string{"a": {"1"}, "b": {"2", "3"}}, "params": map[string]string{}})
	}
	{
		app, ts, e, _ := newTestServer(t)
		defer ts.Close()

		router := NewRouter()
		router.Get("/:sub", func(req *Req, res *Res) { res.JSON(map[string]interface{}{"query": req.Query, "params": req.Params}) })
		app.UseN("/main", router)

		e.GET("/main/sub?a=1&b=2&b=3").Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"query": map[string][]string{"a": {"1"}, "b": {"2", "3"}}, "params": map[string]string{"sub": "sub"}})
	}
	{
		app, ts, e, _ := newTestServer(t)
		defer ts.Close()

		router := NewRouter()
		router.Get("/sub/:sub", func(req *Req, res *Res) { res.JSON(map[string]interface{}{"query": req.Query, "params": req.Params}) })
		app.UseN("/main", router)

		e.GET("/main/sub/sub2?a=1&b=2&b=3").Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"query": map[string][]string{"a": {"1"}, "b": {"2", "3"}}, "params": map[string]string{"sub": "sub2"}})
	}
}
