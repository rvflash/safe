// Package render provides interface to render multiple templates with Gin
package render

import (
	"html/template"

	"github.com/gin-contrib/multitemplate"
)

// Render must be implemented by any handler to render multiple templates.
type Render interface {
	// PageName returns the name of the page.
	PageName() string
	// TmplFiles returns the paths to the template file(s) of this page.
	TmplFiles() []string
	// FuncMap returns the custom template functions to use.
	FuncMap() template.FuncMap
}

// R is the registry of renders.
type R struct {
	multi multitemplate.Render
}

// New returns an instance a new registry of template.
// It will index all HTML templates to render.
func New() *R {
	return &R{multi: multitemplate.New()}
}

// Add adds a render into the registry.
func (r *R) Add(h Render) *R {
	if len(h.FuncMap()) == 0 {
		r.multi.AddFromFiles(h.PageName(), h.TmplFiles()...)
	} else {
		r.multi.AddFromFilesFuncs(h.PageName(), h.FuncMap(), h.TmplFiles()...)
	}
	return r
}

// HTMLRender returns a HTML render as expected by Gin.
func (r *R) HTMLRender() multitemplate.Render {
	return r.multi
}
