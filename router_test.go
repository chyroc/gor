package gor

import (
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
	{
		app, ts, _, as := newTestServer(t)
		defer ts.Close()

		// `/main/:sub` router with app+router
		router := NewRouter()
		router.Get("/:sub", func(req *Req, res *Res) { res.End() })
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
	}
}
