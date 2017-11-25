package gor

import (
	"net/http"
	"testing"
)

func TestQuery(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	// todo path
	app.Get("/params?a=1&c=2&c=3", func(req *Req, res Res) { res.JSON(req.Query) })
	e.GET("/params?a=1&c=2&c=3").Expect().Status(http.StatusOK).JSON().Equal(map[string][]string{"a": {"1"}, "c": {"2", "3"}})
}

func TestHostname(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Get("/", func(req *Req, res Res) { res.Send(req.Hostname) })
	e.GET("/").Expect().Status(http.StatusOK).Text().Equal("127.0.0.1")
}
