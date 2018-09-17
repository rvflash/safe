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

// RenderHandler ...
type RenderHandler interface {
	PathHandler
	render.Render
}

// Router ...
type Router struct {
	app  *app.Safe
	port int
}

// NewRouter ...
func NewRouter(port int, app *app.Safe) *Router {
	return &Router{app: app, port: port}
}

// Addr ...
func (r *Router) Addr() string {
	return ":" + strconv.Itoa(r.port)
}

/*
// CertFile ...
func (r *Router) CertFile() string {
	return r.app.Root().Join("./testdata/server.pem")
}

// KeyFile ...
func (r *Router) KeyFile() string {
	return r.app.Root().Join("./testdata/server.key")
}
*/

// Handler ...
func (r *Router) Handler() http.Handler {
	// Gin with Logger and Recovery middleware
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// HTML Templates
	tmpl := render.New()

	// Landing page
	landing := Landing(r.app)
	tmpl.Add(landing)
	router.GET(landing.Path(), landing.Handle)
	router.POST(landing.Path(), landing.Handle)

	// Homepage
	home := Home(r.app)
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
	router.POST(vaults.Path(), r.secure(vaults.Handle))
	router.PUT(vaults.Path(), r.secure(vaults.Handle))

	// API: logout
	logout := Logout(r.app)
	router.GET(logout.Path(), r.secure(logout.Handle))

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
