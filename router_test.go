package gor

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"
)

func newTestServer(t *testing.T, handler http.Handler) (*httptest.Server, *httpexpect.Expect, *assert.Assertions) {
	ts := httptest.NewServer(handler)
	e := httpexpect.New(t, ts.URL)
	as := assert.New(t)

	return ts, e, as
}

func TestRouter(t *testing.T) {
	app := NewGor()

	app.Get("/", func(req *Req, res Res) {
		res.Send("Hello World")
	})

	app.Post("/", func(req *Req, res Res) {
		res.Send("Hello World2")
	})

	ts, e, _ := newTestServer(t, app)
	defer ts.Close()

	e.GET("/").Expect().Status(http.StatusOK).Text().Equal("Hello World")
	e.POST("/").Expect().Status(http.StatusOK).Text().Equal("Hello World2")
}
