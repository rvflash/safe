// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package router

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rvflash/safe/app"
)

type tags struct {
	app *app.Safe
}

// Tags ...
func Tags(app *app.Safe) PathHandler {
	return &tags{app: app}
}

// Handle implements the Handler interface.
func (h *tags) Handle(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		h.read(c)
	case "POST":
		h.create(c)
	default:
		c.JSON(http.StatusMethodNotAllowed, toErr(errors.New(http.StatusText(http.StatusMethodNotAllowed))))
	}
}

func (h *tags) create(c *gin.Context) {
	// Creates a new one (POST with JSON data).
	type tag struct {
		Name string `form:"name"`
	}
	var (
		data tag
		err  error
	)
	if err = c.Bind(&data); err != nil {
		c.JSON(http.StatusBadRequest, toErr(err))
		return
	}
	t, err := h.app.CreateTag(data.Name)
	if err != nil {
		c.JSON(http.StatusForbidden, toErr(err))
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *tags) read(c *gin.Context) {
	t, err := h.app.ListTagByNames()
	if err != nil {
		c.JSON(http.StatusNotFound, toErr(err))
		return
	}
	c.JSON(http.StatusOK, t)
}

// Path implements the PathHandler interface.
func (h *tags) Path() string {
	return "/tag"
}
