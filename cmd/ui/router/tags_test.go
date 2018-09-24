// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package router_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rvflash/safe/app"
	"github.com/rvflash/safe/cmd/ui/router"
)

func TestTags(t *testing.T) {
	var (
		dt = []struct {
			// Handler
			app *app.Safe
			// Request
			method string
			json   io.Reader
			// Response
			code int
			body string
		}{
			{
				app:    notConnected,
				code:   http.StatusUnauthorized,
				body:   `{"error":"Unauthorized"}`,
				method: "GET",
			},
			{
				app:    connected,
				body:   `["Job","Private"]`,
				code:   http.StatusOK,
				method: "GET",
			},
			{
				app:    notConnected,
				method: "POST",
				code:   http.StatusUnauthorized,
				body:   `{"error":"Unauthorized"}`,
			},
			{
				app:    connected,
				method: "POST",
				code:   http.StatusBadRequest,
				body:   `{"error":"missing form body"}`,
			},
			{
				app:    connected,
				method: "POST",
				code:   http.StatusOK,
				json:   strings.NewReader(`{"name":"Test"}`),
				body:   `"Test"`,
			},
			/* @todo test with custom handler
			{
				app:    connected,
				method: "PUT",
				code:   http.StatusMethodNotAllowed,
				body:   `{"error":"Method Not Allowed"}`,
			},
			*/
		}
	)
	var addr string
	for i, tt := range dt {
		// Creates the test environment.
		r := router.NewRouter(80, tt.app, true)
		if addr = r.Addr(); addr != ":80" {
			t.Fatalf("unexpected server address: %s", addr)
		}
		h := router.Tags(tt.app)

		// Builds the request.
		req, err := http.NewRequest(tt.method, h.Path(), tt.json)
		if err != nil {
			t.Fatal(err)
		}
		if tt.method == "POST" && tt.json != nil {
			req.Header.Set("Content-type", "application/json")
		}
		// Listens the response.
		w := httptest.NewRecorder()
		r.Handler().ServeHTTP(w, req)
		if w.Code != tt.code {
			t.Errorf("%d. mismatch status code: got=%d exp=%d\n", i, w.Code, tt.code)
		}
		if w.Body.String() != tt.body {
			t.Errorf("%d. mismatch content: got=%s exp=%s\n", i, w.Body, tt.body)
		}
	}
}
