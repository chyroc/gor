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

func renderParamQuery(req *Req, res *Res) {
	res.JSON(map[string]interface{}{
		"params": req.Params,
		"query":  req.Query,
	})
}
func TestUse(t *testing.T) {
	// app
	//{
	//	app, ts, e, as := newTestServer(t)
	//	defer ts.Close()
	//
	//	// `/` `/2` `/3/3`
	//	app.Get("/", func(req *Req, res *Res) { renderParamQuery(req, res) })
	//	app.Get("/2", func(req *Req, res *Res) { renderParamQuery(req, res) })
	//	app.Get("/3/3", func(req *Req, res *Res) { renderParamQuery(req, res) })
	//	as.Len(app.routes, 3)
	//	assertRoute(as, "GET", "/", fullMatch, handlerFunc, app.routes[0])
	//	assertRoute(as, "GET", "/2", fullMatch, handlerFunc, app.routes[1])
	//	assertRoute(as, "GET", "/3/3", fullMatch, handlerFunc, app.routes[2])
	//
	//	e.GET("/?a=b").Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"params": map[string]string{}, "query": map[string][]string{"a": {"b"}}})
	//	e.GET("/2?a=b").Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"params": map[string]string{}, "query": map[string][]string{"a": {"b"}}})
	//	e.GET("/3/3?a=b").Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"params": map[string]string{}, "query": map[string][]string{"a": {"b"}}})
	//}
	//{
	//	app, ts, e, as := newTestServer(t)
	//	defer ts.Close()
	//
	//	// `/:user` `/1/:user` `/2/:user/not-param/:name`
	//	app.Get("/1/:user", func(req *Req, res *Res) { renderParamQuery(req, res) })
	//	app.Get("/:user", func(req *Req, res *Res) { renderParamQuery(req, res) })
	//	app.Get("/2/:user/not-param/:name", func(req *Req, res *Res) { renderParamQuery(req, res) })
	//	as.Len(app.routes, 3)
	//	assertRoute(as, "GET", "/1/:user", fullMatch, handlerFunc, app.routes[0])
	//	assertRoute(as, "GET", "/:user", fullMatch, handlerFunc, app.routes[1])
	//	assertRoute(as, "GET", "/2/:user/not-param/:name", fullMatch, handlerFunc, app.routes[2])
	//
	//	e.GET("/1/user?a=b").Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"params": map[string]string{"user": "user"}, "query": map[string][]string{"a": {"b"}}})
	//	e.GET("/user?a=b").Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"params": map[string]string{"user": "user"}, "query": map[string][]string{"a": {"b"}}})
	//	e.GET("/2/user/not-param/name?a=b").Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"params": map[string]string{"user": "user", "name": "name"}, "query": map[string][]string{"a": {"b"}}})
	//}
	{
		// only app use
		app, ts, e, as := newTestServer(t)
		defer ts.Close()

		app.Use("/1", func(req *Req, res *Res) { renderParamQuery(req, res) })
		app.Use("/2", func(req *Req, res *Res, next Next) { renderParamQuery(req, res) })

		as.Len(app.routes, 2)
		assertRoute(as, "ALL", "/1", preMatch, handlerFunc, app.routes[0])
		assertRoute(as, "ALL", "/2", preMatch, handlerFuncNext, app.routes[1])

		//e.GET("/1?a=b").Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"params": map[string]string{}, "query": map[string][]string{"a": {"b"}}})
		//e.GET("/1/2/3?a=b").Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"params": map[string]string{}, "query": map[string][]string{"a": {"b"}}})
		e.GET("/2?a=b").Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"params": map[string]string{}, "query": map[string][]string{"a": {"b"}}})
		//e.GET("/2/2/3?a=b").Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"params": map[string]string{}, "query": map[string][]string{"a": {"b"}}})
	}
	{
		// app use and router
		app, ts, e, as := newTestServer(t)
		defer ts.Close()
		router := NewRouter()

		app.Get("/no-sub", func(req *Req, res *Res) { renderParamQuery(req, res) })
		router.Get("/1", func(req *Req, res *Res) { renderParamQuery(req, res) })
		router.Get("/2", func(req *Req, res *Res) { renderParamQuery(req, res) })
		app.Use("/main", router)

		as.Len(router.routes, 2)
		assertRoute(as, "GET", "/1", fullMatch, handlerFunc, router.routes[0])
		assertRoute(as, "GET", "/2", fullMatch, handlerFunc, router.routes[1])

		as.Len(app.routes, 2)
		assertRoute(as, "GET", "/no-sub", fullMatch, handlerFunc, app.routes[0])
		assertRoute(as, "ALL", "/main", preMatch, unkonwFunc, app.routes[1])
		as.Len(app.routes[1].children, 2)
		assertRoute(as, "GET", "/1", fullMatch, handlerFunc, app.routes[1].children[0])
		assertRoute(as, "GET", "/2", fullMatch, handlerFunc, app.routes[1].children[1])

		e.GET("/no-sub?a=b").Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"params": map[string]string{}, "query": map[string][]string{"a": {"b"}}})
		e.GET("/main/1?a=b").Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"params": map[string]string{}, "query": map[string][]string{"a": {"b"}}})
		e.GET("/main/2?a=b").Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"params": map[string]string{}, "query": map[string][]string{"a": {"b"}}})
	}
	{
		// app use and router + params
		app, ts, e, as := newTestServer(t)
		defer ts.Close()
		router := NewRouter()

		app.Get("/no-sub/:name0", func(req *Req, res *Res) { renderParamQuery(req, res) })
		router.Get("/1/:name1", func(req *Req, res *Res) { renderParamQuery(req, res) })
		router.Get("/2/:name2", func(req *Req, res *Res) { renderParamQuery(req, res) })
		app.Use("/main", router)

		as.Len(router.routes, 2)
		assertRoute(as, "GET", "/1/:name1", fullMatch, handlerFunc, router.routes[0])
		assertRoute(as, "GET", "/2/:name2", fullMatch, handlerFunc, router.routes[1])

		as.Len(app.routes, 2)
		assertRoute(as, "GET", "/no-sub/:name0", fullMatch, handlerFunc, app.routes[0])
		assertRoute(as, "ALL", "/main", preMatch, unkonwFunc, app.routes[1])
		as.Len(app.routes[1].children, 2)
		assertRoute(as, "GET", "/1/:name1", fullMatch, handlerFunc, app.routes[1].children[0])
		assertRoute(as, "GET", "/2/:name2", fullMatch, handlerFunc, app.routes[1].children[1])

		e.GET("/no-sub/name0?a=b").Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"params": map[string]string{"name0": "name0"}, "query": map[string][]string{"a": {"b"}}})
		e.GET("/main/1/name1?a=b").Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"params": map[string]string{"name1": "name1"}, "query": map[string][]string{"a": {"b"}}})
		e.GET("/main/2/name2?a=b").Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"params": map[string]string{"name2": "name2"}, "query": map[string][]string{"a": {"b"}}})
	}
}
