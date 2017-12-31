package gor

import (
	"net/http"
	"testing"
)

func TestProtocolSecure(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Get("/1", func(req *Req, res *Res) { res.Send(req.Protocol) })
	e.GET("/1").Expect().Status(http.StatusOK).Text().Equal("http")

	app.Get("/2", func(req *Req, res *Res) { res.Send(req.Secure) })
	e.GET("/2").Expect().Status(http.StatusOK).Text().Equal("false")
}

func TestHostname(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Get("/", func(req *Req, res *Res) { res.Send(req.Hostname) })
	e.GET("/").Expect().Status(http.StatusOK).Text().Equal("127.0.0.1")
}

func TestBaseURLOriginalURL(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Get("/1", func(req *Req, res *Res) { res.Send(req.BaseURL) })
	e.GET("/1").Expect().Status(http.StatusOK).Text().Equal("/1")
	e.GET("/1?a=1").Expect().Status(http.StatusOK).Text().Equal("/1")

	app.Get("/2", func(req *Req, res *Res) { res.Send(req.OriginalURL) })
	e.GET("/2").Expect().Status(http.StatusOK).Text().Equal("/2")
	e.GET("/2?a=1").Expect().Status(http.StatusOK).Text().Equal("/2?a=1")
}

func TestQuery(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Get("/query", func(req *Req, res *Res) { res.JSON(req.Query) })
	e.GET("/query?a=1&c=2&c=3").Expect().Status(http.StatusOK).JSON().Equal(map[string][]string{"a": {"1"}, "c": {"2", "3"}})
	e.GET("/query?q=a+b").Expect().Status(http.StatusOK).JSON().Equal(map[string][]string{"q": {"a b"}})
	// todo
	// e.GET("/query?order=desc&shoe[color]=blue&shoe[type]=converse").Expect().Status(http.StatusOK).JSON().Equal(map[string][]string{"q": {"a b"}})
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

func TestBody(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Post("/", func(req *Req, res *Res) {
		res.JSON(map[string]interface{}{"json": req.Body.JSON, "form-data": req.Body.FormData, "form-url-encoded": req.Body.FormURLEncoded})
	})
	e.POST("/").WithJSON(map[string]string{"s": "d"}).Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"form-data": nil, "form-url-encoded": nil, "json": map[string]string{"s": "d"}})
	e.POST("/").WithForm(map[string]string{"s": "d"}).Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"form-data": nil, "form-url-encoded": map[string][]string{"s": {"d"}}, "json": nil})
	e.POST("/").WithMultipart().WithForm(map[string]string{"s": "d"}).Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"form-data": map[string][]string{"s": {"d"}}, "form-url-encoded": nil, "json": nil})
}

func TestBind(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Post("/", func(req *Req, res *Res) {
		var t struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		if err := req.BindJSON(&t); err != nil {
			res.Error(err.Error())
			return
		}
		res.JSON(t)
	})
	e.POST("/").WithJSON(map[string]interface{}{"name": "1", "age": 24}).Expect().Status(http.StatusOK).JSON().Equal(map[string]interface{}{"name": "1", "age": 24})
}
