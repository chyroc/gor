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

func TestSend(t *testing.T) {
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

func TestJson(t *testing.T) {
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
			app := NewGor()

			app.Get("/", func(req *Req, res Res) {
				res.Json(v)
			})

			ts, e, _ := newTestServer(t, app)
			defer ts.Close()

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
			app := NewGor()

			app.Get("/", func(req *Req, res Res) {
				res.Json(v)
			})

			ts, e, _ := newTestServer(t, app)
			defer ts.Close()

			e.GET("/").Expect().Status(http.StatusInternalServerError).Text().Equal(msg)
		}
	}
}
