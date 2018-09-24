// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package router

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rvflash/elapsed"
	"github.com/rvflash/safe"
	"github.com/rvflash/safe/app"
)

type vaults struct {
	app *app.Safe
}

type vault struct {
	Name       string `json:"name"`
	Tag        string `json:"tag"`
	LastUpdate string `json:"last_upd"`
	Username   string `json:"user"`
	Password   string `json:"pass"`
	Strength   int    `json:"strength"`
	Safe       string `json:"safe,omitempty"`
	URL        string `json:"url,omitempty"`
	Note       string `json:"note,omitempty"`
}

func newVault(v *safe.Vault) *vault {
	if v == nil || v.Tag() == nil || v.Login() == nil {
		return nil
	}
	// Do not extract the plain password.
	d := &vault{
		Name:       v.Name(),
		Tag:        v.Tag().Name(),
		LastUpdate: elapsed.Time(v.Login().LastUpdate),
		Username:   v.Login().Name,
		Note:       v.Login().Note,
		Strength:   v.Login().Strength(),
	}
	if v.Login().URL != nil {
		d.URL = v.Login().URL.String()
	}
	if ok, err := v.Login().Safe(); !ok {
		d.Safe = err.Error()
	}
	return d
}

// Vaults ...
func Vaults(app *app.Safe) OptionalPathHandler {
	return &vaults{app: app}
}

func (h *vaults) create(c *gin.Context) {
	var l app.Login
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, toErr(err))
		return
	}
	v, err := h.app.CreateVault(c.Param("key"), c.Param("tag"), l)
	if err != nil {
		c.JSON(http.StatusForbidden, toErr(err))
		return
	}
	c.JSON(http.StatusOK, newVault(v))
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
	if _, ok := c.GetQuery("pwd"); ok {
		// base64 encoding password
		c.JSON(http.StatusOK, base64.StdEncoding.EncodeToString([]byte(v.Login().Password)))
	} else {
		c.JSON(http.StatusOK, newVault(v))
	}
}

type vaultHead struct {
	Name     string `json:"name"`
	Username string `json:"user"`
}

func newVaultHead(v *safe.Vault) *vaultHead {
	if v == nil || v.Login() == nil {
		return nil
	}
	d := &vaultHead{
		Name:     v.Name(),
		Username: v.Login().Name,
	}
	return d
}

func (h *vaults) readAll(c *gin.Context) {
	d, err := h.app.ListVaultByNames(c.Param("tag"), c.Query("prefix"))
	if err != nil {
		c.JSON(http.StatusNotFound, toErr(err))
		return
	}
	vh := make([]*vaultHead, len(d))
	for i, v := range d {
		vh[i] = newVaultHead(v)
	}
	c.JSON(http.StatusOK, vh)
}

func (h *vaults) update(c *gin.Context) {
	var l app.Login
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, toErr(err))
		return
	}
	v, err := h.app.UpdateVault(c.Param("key"), c.Param("tag"), l)
	if err != nil {
		c.JSON(http.StatusNotFound, toErr(err))
		return
	}
	c.JSON(http.StatusOK, newVault(v))
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

// OptionalPath implements the OptionalPathHandler interface.
func (h *vaults) OptionalPath() string {
	return "/tags/:tag/vaults/"
}
