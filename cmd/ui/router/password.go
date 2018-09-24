// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package router

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rvflash/safe/app"
)

type password struct{}

// Password ...
func Password() PathHandler {
	return &password{}
}

// Handle implements the Handler interface.
func (h *password) Handle(c *gin.Context) {
	switch c.Request.Method {
	case "POST":
		h.create(c)
	default:
		c.JSON(http.StatusMethodNotAllowed, toErr(errors.New(http.StatusText(http.StatusMethodNotAllowed))))
	}
}

func (h *password) create(c *gin.Context) {
	type pwd struct {
		Length      int  `form:"length"`
		Digits      int  `form:"digits"`
		Symbols     int  `form:"symbols"`
		NoUpper     bool `form:"no_upper"`
		AllowRepeat bool `form:"repeat"`
	}
	var (
		data pwd
		err  error
	)
	if err = c.Bind(&data); err != nil {
		c.JSON(http.StatusBadRequest, toErr(err))
		return
	}
	p, err := app.GeneratePassword(data.Length, data.Digits, data.Symbols, data.NoUpper, data.AllowRepeat)
	if err != nil {
		c.JSON(http.StatusForbidden, toErr(err))
		return
	}
	c.JSON(http.StatusOK, base64.StdEncoding.EncodeToString([]byte(p)))
}

// Path implements the PathHandler interface.
func (h *password) Path() string {
	return "/pwd"
}
