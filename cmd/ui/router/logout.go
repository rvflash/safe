// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rvflash/safe/app"
)

type logoutPage struct {
	app *app.Safe
}

// Logout ...
func Logout(app *app.Safe) PathHandler {
	return &logoutPage{app: app}
}

// Handle implements the Handler interface.
func (h *logoutPage) Handle(c *gin.Context) {
	h.app.LogOut()
	c.Redirect(http.StatusTemporaryRedirect, Landing(h.app, false).Path())
}

// Path implements the PathHandler interface.
func (h *logoutPage) Path() string {
	return "/logout"
}
