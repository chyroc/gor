package gor

import (
	"net/http"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
)

func TestStatus(t *testing.T) {
	{
		app, ts, e, _ := newTestServer(t)
		defer ts.Close()

		app.Get("/", func(req *Req, res *Res) { res.Status(http.StatusAccepted).Send("Hello World") })
		e.GET("/").Expect().Status(http.StatusAccepted).Text().Equal("Hello World")

	}
	{
		app, ts, e, _ := newTestServer(t)
		defer ts.Close()

		app.Get("/", func(req *Req, res *Res) { res.Status(-1).Send("Hello World") })
		e.GET("/").Expect().Status(http.StatusInternalServerError).Text().Equal("http status code is invalid")
	}
}

func TestSendStatus(t *testing.T) {
	{
		app, ts, e, _ := newTestServer(t)
		defer ts.Close()

		app.Get("/", func(req *Req, res *Res) { res.SendStatus(http.StatusAccepted) })
		e.GET("/").Expect().Status(http.StatusAccepted).Text().Equal("Accepted")
	}
	{
		app, ts, e, _ := newTestServer(t)
		defer ts.Close()

		app.Get("/", func(req *Req, res *Res) { res.SendStatus(-1) })
		e.GET("/").Expect().Status(http.StatusInternalServerError).Text().Equal("http status code is invalid")
	}
}

func TestSend(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Get("/", func(req *Req, res *Res) { res.Send("Hello World") })
	e.GET("/").Expect().Status(http.StatusOK).Text().Equal("Hello World")
}

func TestJSON(t *testing.T) {
	{
		for _, v := range []interface{}{
			struct {
				// struct
				Name     string `json:"name"`
				unExport string
			}{Name: "chyroc", unExport: "24"},
			map[string]string{"1": "2"}, // map
			[]string{"a", "b"},          // slice
			[1]int{1},                   // array
			nil,                         //nil
		} {
			app, ts, e, _ := newTestServer(t)
			defer ts.Close()

			app.Get("/", func(req *Req, res *Res) { res.JSON(v) })
			e.GET("/").Expect().Status(http.StatusOK).JSON().Equal(v)
		}
	}
	{
		// todo :Uintptr Func Chan Interface Ptr
		for msg, v := range map[string]interface{}{
			"[response type unsupported] [int] 1":             1,             // int
			"[response type unsupported] [int8] 1":            int8(1),       // int8
			"[response type unsupported] [int16] 1":           int16(1),      // int16
			"[response type unsupported] [int32] 1":           int32(1),      // int32
			"[response type unsupported] [int64] 1":           int64(1),      // int64
			"[response type unsupported] [uint] 1":            uint(1),       // uint
			"[response type unsupported] [uint8] 1":           uint8(1),      // uint8
			"[response type unsupported] [uint16] 1":          uint16(1),     // uint16
			"[response type unsupported] [uint32] 1":          uint32(1),     // uint32
			"[response type unsupported] [uint64] 1":          uint64(1),     // uint64
			"[response type unsupported] [float32] 1.1":       float32(1.1),  // float32
			"[response type unsupported] [float64] 1.1":       float64(1.1),  // float64
			"[response type unsupported] [complex64] (1+0i)":  complex64(1),  // complex64
			"[response type unsupported] [complex128] (1+0i)": complex128(1), // complex128
			"[response type unsupported] [string] string":     "string",      // string
			"[response type unsupported] [bool] false":        false,         // bool
		} {
			app, ts, e, _ := newTestServer(t)
			defer ts.Close()

			app.Get("/", func(req *Req, res *Res) { res.JSON(v) })
			e.GET("/").Expect().Status(http.StatusInternalServerError).Text().Equal(msg)
		}
	}
}

func TestRedirect(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Get("/", func(req *Req, res *Res) { res.Redirect("/b") })
	app.Get("/b", func(req *Req, res *Res) { res.Send("b") })
	e.GET("/").Expect().Status(http.StatusOK).Text().Equal("b")
	httpexpect.WithConfig(httpexpect.Config{
		Client:   &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }},
		BaseURL:  ts.URL,
		Reporter: httpexpect.NewAssertReporter(t),
	}).GET("/").Expect().Status(http.StatusFound).Text().Equal("Found. Redirecting to /b")
}

func TestAddHeader(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Get("/", func(req *Req, res *Res) {
		res.AddHeader("h", "h1")
		res.AddHeader("h", "h2")
	})
	e.GET("/").Expect().Status(http.StatusOK).Headers().ValueEqual("H", []string{"h1", "h2"})
}

func TestSetCookie(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Get("/1", func(req *Req, res *Res) { res.SetCookie("c", "c1", Cookie{}, Cookie{}) })
	app.Get("/2", func(req *Req, res *Res) { res.SetCookie("c", "c1", Cookie{}) })
	ti := time.Now().Add(time.Minute)
	app.Get("/3", func(req *Req, res *Res) {
		res.SetCookie("c", "c1", Cookie{
			Path:    "/",
			Expires: time.Unix(int64(ti.Second()), 0),
		})
	})

	e.GET("/1").Expect().Status(http.StatusInternalServerError).Text().Equal("only support one cookie option")
	e.GET("/2").Expect().Status(http.StatusOK).Cookie("c").Value().Equal("c1")
	e.GET("/3").Expect().Status(http.StatusOK).Cookie("c").Path().Equal("/")
	e.GET("/3").Expect().Status(http.StatusOK).Cookie("c").Expires().Equal(time.Unix(int64(ti.Second()), 0))
}

func TestEnd(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Get("/", func(req *Req, res *Res) { res.End() })
	e.GET("/").Expect().Status(http.StatusOK).Text().Equal("")
}
