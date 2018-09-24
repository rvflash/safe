// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package router

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rvflash/safe/app"
	"github.com/rvflash/safe/cmd/ui/render"
)

// Handler ...
type Handler interface {
	Handle(c *gin.Context)
}

// PathHandler ...
type PathHandler interface {
	Handler
	Path() string
}

// OptionalPathHandler ...
type OptionalPathHandler interface {
	PathHandler
	OptionalPath() string
}

// RenderHandler ...
type RenderHandler interface {
	PathHandler
	render.Render
}

// Router ...
type Router struct {
	app   *app.Safe
	port  int
	debug bool
}

// NewRouter ...
func NewRouter(port int, app *app.Safe, test bool) *Router {
	return &Router{app: app, port: port, debug: test}
}

// Addr ...
func (r *Router) Addr() string {
	return ":" + strconv.Itoa(r.port)
}

// Handler ...
func (r *Router) Handler() http.Handler {
	if r.debug {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	if !r.debug {
		router.Use(gin.Logger(), gin.Recovery())
	}

	// HTML Templates
	tmpl := render.New()

	// Landing page
	landing := Landing(r.app, r.debug)
	tmpl.Add(landing)
	router.GET(landing.Path(), landing.Handle)
	router.POST(landing.Path(), landing.Handle)

	// Homepage
	home := Home(r.app, r.debug)
	tmpl.Add(home)
	router.GET(home.Path(), r.secure(home.Handle))

	// API: tags
	tags := Tags(r.app)
	router.GET(tags.Path(), r.secure(tags.Handle))
	router.POST(tags.Path(), r.secure(tags.Handle))

	// API: vaults
	vaults := Vaults(r.app)
	router.DELETE(vaults.Path(), r.secure(vaults.Handle))
	router.GET(vaults.Path(), r.secure(vaults.Handle))
	router.GET(vaults.OptionalPath(), r.secure(vaults.Handle))
	router.POST(vaults.Path(), r.secure(vaults.Handle))
	router.PUT(vaults.Path(), r.secure(vaults.Handle))

	// API: password (generate a new one)
	pwd := Password()
	router.POST(pwd.Path(), pwd.Handle)

	// API: logout
	logout := Logout(r.app)
	router.GET(logout.Path(), logout.Handle)

	// Static media.
	router.Static("/css/", r.app.Root().Join("./static/css/"))
	router.Static("/js/", r.app.Root().Join("./static/js/"))

	// Not found
	router.NoRoute(NotFound().Handle)

	// Overloads the HTML render.
	router.HTMLRender = tmpl.HTMLRender()

	return router
}

func (r *Router) secure(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if r.app.Logged() != nil {
			c.JSON(http.StatusUnauthorized, toErr(errors.New(http.StatusText(http.StatusUnauthorized))))
			return
		}
		next(c)
	}
}
