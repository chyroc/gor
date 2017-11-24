package gor

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestStatus(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Get("/", func(req *Req, res Res) { res.Status(http.StatusAccepted).Send("Hello World") })
	e.GET("/").Expect().Status(http.StatusAccepted).Text().Equal("Hello World")
}

func TestSend(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Get("/", func(req *Req, res Res) { res.Send("Hello World") })
	e.GET("/").Expect().Status(http.StatusOK).Text().Equal("Hello World")
}

func TestJSON(t *testing.T) {
	{
		for _, v := range []interface{}{
			struct {
				Name     string `json:"name"`
				unExport string
			}{Name: "chyroc", unExport: "24"}, // struct
			map[string]string{"1": "2"}, // map
			[]string{"a", "b"},          // slice
			[1]int{1},                   // array
		} {
			app, ts, e, _ := newTestServer(t)
			defer ts.Close()

			app.Get("/", func(req *Req, res Res) { res.JSON(v) })
			e.GET("/").Expect().Status(http.StatusOK).JSON().Equal(v)
		}
	}
	{
		// todo :Uintptr Func Chan Interface Ptr
		for msg, v := range map[string]interface{}{
			"[response type unsupported] [nil] <nil>\n":         nil,           // nil
			"[response type unsupported] [int] 1\n":             1,             // int
			"[response type unsupported] [int8] 1\n":            int8(1),       // int8
			"[response type unsupported] [int16] 1\n":           int16(1),      // int16
			"[response type unsupported] [int32] 1\n":           int32(1),      // int32
			"[response type unsupported] [int64] 1\n":           int64(1),      // int64
			"[response type unsupported] [uint] 1\n":            uint(1),       // uint
			"[response type unsupported] [uint8] 1\n":           uint8(1),      // uint8
			"[response type unsupported] [uint16] 1\n":          uint16(1),     // uint16
			"[response type unsupported] [uint32] 1\n":          uint32(1),     // uint32
			"[response type unsupported] [uint64] 1\n":          uint64(1),     // uint64
			"[response type unsupported] [float32] 1.1\n":       float32(1.1),  // float32
			"[response type unsupported] [float64] 1.1\n":       float64(1.1),  // float64
			"[response type unsupported] [complex64] (1+0i)\n":  complex64(1),  // complex64
			"[response type unsupported] [complex128] (1+0i)\n": complex128(1), // complex128
			"[response type unsupported] [string] string\n":     "string",      // string
			"[response type unsupported] [bool] false\n":        false,         // bool
		} {
			app, ts, e, _ := newTestServer(t)
			defer ts.Close()

			app.Get("/", func(req *Req, res Res) { res.JSON(v) })
			e.GET("/").Expect().Status(http.StatusInternalServerError).Text().Equal(msg)
		}
	}
}

func TestRedirect(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Get("/", func(req *Req, res Res) { res.Redirect("/b") })
	app.Get("/b", func(req *Req, res Res) { res.Send("b") })
	e.GET("/").Expect().Status(http.StatusOK).Text().Equal("b")
	httpexpect.WithConfig(httpexpect.Config{
		Client:   &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }},
		BaseURL:  ts.URL,
		Reporter: httpexpect.NewAssertReporter(t),
	}).GET("/").Expect().Status(http.StatusFound).Text().Equal("Found. Redirecting to /b")
}

func TestEnd(t *testing.T) {
	app, ts, e, _ := newTestServer(t)
	defer ts.Close()

	app.Get("/", func(req *Req, res Res) { res.End() })
	e.GET("/").Expect().Status(http.StatusOK).Text().Equal("")
}
