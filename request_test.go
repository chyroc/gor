package gor

import (
	"net/http"
	"testing"
)

func TestQuery(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	// todo path
	app.Get("/query", func(req *Req, res *Res) { res.JSON(req.Query) })
	e.GET("/query?a=1&c=2&c=3").Expect().Status(http.StatusOK).JSON().Equal(map[string][]string{"a": {"1"}, "c": {"2", "3"}})
}

func TestHostname(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Get("/", func(req *Req, res *Res) { res.Send(req.Hostname) })
	e.GET("/").Expect().Status(http.StatusOK).Text().Equal("127.0.0.1")
}

func TestBaseURL(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Get("/a", func(req *Req, res *Res) { res.Send(req.BaseURL) })
	e.GET("/a?a=2").Expect().Status(http.StatusOK).Text().Equal("/a")
}

func TestParams(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Get("/a/:user", func(req *Req, res *Res) { res.JSON(req.Params) })
	e.GET("/a").Expect().Status(http.StatusNotFound).Text().Equal("Not Found")
	e.GET("/a/user").Expect().Status(http.StatusOK).JSON().Equal(map[string]string{"user": "user"})
	e.GET("/a/user/user").Expect().Status(http.StatusNotFound).Text().Equal("Not Found")

	app.Get("/b/:user/:name", func(req *Req, res *Res) { res.JSON(req.Params) })
	e.GET("/b").Expect().Status(http.StatusNotFound).Text().Equal("Not Found")
	e.GET("/b/user").Expect().Status(http.StatusNotFound).Text().Equal("Not Found")
	e.GET("/b/user/name").Expect().Status(http.StatusOK).JSON().Equal(map[string]string{"user": "user", "name": "name"})
	e.GET("/b/user/name/name").Expect().Status(http.StatusNotFound).Text().Equal("Not Found")

	app.Get("/c/:user/noparam/:name", func(req *Req, res *Res) { res.JSON(req.Params) })
	e.GET("/c").Expect().Status(http.StatusNotFound).Text().Equal("Not Found")
	e.GET("/c/user").Expect().Status(http.StatusNotFound).Text().Equal("Not Found")
	e.GET("/c/user/noparam").Expect().Status(http.StatusNotFound).Text().Equal("Not Found")
	e.GET("/c/user/noparam/name").Expect().Status(http.StatusOK).JSON().Equal(map[string]string{"user": "user", "name": "name"})
	e.GET("/c/user/no-match-param/name").Expect().Status(http.StatusNotFound).Text().Equal("Not Found")
	e.GET("/c/user/noparam/name/name").Expect().Status(http.StatusNotFound).Text().Equal("Not Found")
}
