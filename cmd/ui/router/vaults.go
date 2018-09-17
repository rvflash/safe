// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package router

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/rvflash/safe/app"
)

type vaults struct {
	app *app.Safe
}

// Vaults ...
func Vaults(app *app.Safe) PathHandler {
	return &vaults{app: app}
}

func (h *vaults) create(c *gin.Context) {
	d, err := h.data(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, toErr(err))
		return
	}
	v, err := h.app.CreateVault(c.Param("key"), c.Param("tag"), d)
	if err != nil {
		c.JSON(http.StatusForbidden, toErr(err))
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *vaults) delete(c *gin.Context) {
	err := h.app.DeleteVault(c.Param("key"), c.Param("tag"))
	if err != nil {
		c.JSON(http.StatusNotFound, toErr(err))
		return
	}
	c.Status(http.StatusOK)
}

func (h *vaults) readOne(c *gin.Context) {
	v, err := h.app.Vault(c.Param("key"), c.Param("tag"))
	if err != nil {
		c.JSON(http.StatusNotFound, toErr(err))
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *vaults) readAll(c *gin.Context) {
	v, err := h.app.ListVaultByNames(c.Param("tag"), c.Query("prefix"))
	if err != nil {
		c.JSON(http.StatusNotFound, toErr(err))
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *vaults) update(c *gin.Context) {
	d, err := h.data(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, toErr(err))
		return
	}
	v, err := h.app.UpdateVault(c.Param("key"), c.Param("tag"), d)
	if err != nil {
		c.JSON(http.StatusNotFound, toErr(err))
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *vaults) data(c *gin.Context) (url.Values, error) {
	type vault struct {
		Name     string `json:"name"`
		Password string `json:"pass"`
		URL      string `json:"url"`
		Note     string `json:"note"`
	}
	var obj vault
	if err := c.BindJSON(&obj); err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set(app.FormUser, obj.Name)
	v.Set(app.FormPass, obj.Password)
	v.Set(app.FormURL, obj.URL)
	v.Set(app.FormNote, obj.Note)

	return v, nil
}

// Handle implements the Handler interface.
func (h *vaults) Handle(c *gin.Context) {
	switch c.Request.Method {
	case "DELETE":
		h.delete(c)
	case "GET":
		if c.Param("key") == "" {
			h.readAll(c)
		} else {
			h.readOne(c)
		}
	case "POST":
		h.create(c)
	case "PUT":
		h.update(c)
	default:
		c.JSON(http.StatusMethodNotAllowed, toErr(errors.New(http.StatusText(http.StatusMethodNotAllowed))))
	}
}

// Path implements the PathHandler interface.
func (h *vaults) Path() string {
	return "/tags/:tag/vaults/:key"
}
