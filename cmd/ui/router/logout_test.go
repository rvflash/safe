// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package router_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rvflash/safe/app"
	"github.com/rvflash/safe/cmd/ui/router"
)

var tmpRedirect = `<a href="/">Temporary Redirect</a>`

func TestLogout(t *testing.T) {
	var (
		dt = []struct {
			app  *app.Safe
			code int
			body string
		}{
			{app: notConnected, code: http.StatusTemporaryRedirect, body: tmpRedirect},
			{app: newApp(true), code: http.StatusTemporaryRedirect, body: tmpRedirect},
		}
	)
	var addr string
	for i, tt := range dt {
		// Creates the test environment.
		r := router.NewRouter(80, tt.app, true)
		if addr = r.Addr(); addr != ":80" {
			t.Fatalf("unexpected server address: %s", addr)
		}
		h := router.Logout(tt.app)

		// Builds the request.
		req, err := http.NewRequest("GET", h.Path(), nil)
		if err != nil {
			t.Fatal(err)
		}
		// Listens the response.
		w := httptest.NewRecorder()
		r.Handler().ServeHTTP(w, req)
		if w.Code != tt.code {
			t.Errorf("%d. mismatch status code: got=%d exp=%d\n", i, w.Code, tt.code)
		}
		if !strings.HasPrefix(w.Body.String(), tt.body) {
			t.Errorf("%d. mismatch content:\ngot=%q\nexp=%q\n", i, w.Body, tt.body)
		}
		if err = tt.app.Logged(); err == nil {
			t.Errorf("%d. not logout: got=%q\n", i, err)
		}
	}
}
