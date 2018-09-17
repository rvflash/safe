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

type landingPage struct {
	app *app.Safe
}

// Landing ...
func Landing(app *app.Safe) RenderHandler {
	return &landingPage{app: app}
}

// Handle implements the Handler interface.
func (h *landingPage) Handle(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		if logged := h.app.Logged(); logged == nil {
			c.Redirect(http.StatusMovedPermanently, Home(h.app).Path())
		} else {
			c.HTML(http.StatusOK, h.PageName(), gin.H{"SignUp": logged == app.ErrNotFound})
		}
	case "POST":
		type pass struct {
			Phrase string `form:"phrase"`
		}
		var (
			data pass
			err  error
		)
		if err = c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, toErr(err))
			return
		}
		if err = h.app.Login(data.Phrase); err != nil {
			c.JSON(http.StatusUnauthorized, toErr(err))
			return
		}
		c.JSON(http.StatusOK, gin.H{"goto": Home(h.app).Path()})
	}
}

// PageName implements the Handler interface.
func (h *landingPage) PageName() string {
	return "landing"
}

// Path implements the PathHandler interface.
func (h *landingPage) Path() string {
	return "/"
}

// TmplFiles implements the Handler interface.
func (h *landingPage) TmplFiles() []string {
	return []string{
		h.app.Root().Join("template/landing.tmpl"),
		h.app.Root().Join("template/common/foot.tmpl"),
		h.app.Root().Join("template/common/head.tmpl"),
	}
}

// FuncMap implements the Handler interface.
func (h *landingPage) FuncMap() template.FuncMap {
	return nil
}
