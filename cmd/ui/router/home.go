// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package router

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rvflash/safe/app"
)

type homePage struct {
	app   *app.Safe
	debug bool
}

// Home ...
func Home(app *app.Safe, debug bool) RenderHandler {
	return &homePage{app: app, debug: debug}
}

// Handle implements the Handler interface.
func (h *homePage) Handle(c *gin.Context) {
	c.HTML(http.StatusOK, h.PageName(), gin.H{"IsDebug": h.debug})
}

// PageName implements the Handler interface.
func (h *homePage) PageName() string {
	return "home"
}

// Path implements the PathHandler interface.
func (h *homePage) Path() string {
	return "/home"
}

// TmplFiles implements the Handler interface.
func (h *homePage) TmplFiles() []string {
	return []string{
		h.app.Root().Join("template/home.tmpl"),
		h.app.Root().Join("template/common/foot.tmpl"),
		h.app.Root().Join("template/common/head.tmpl"),
	}
}

// PathName implements the render.Render interface.
func (h *homePage) FuncMap() template.FuncMap {
	return nil
}
