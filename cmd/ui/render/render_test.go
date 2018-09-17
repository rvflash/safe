package render_test

import (
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/rvflash/safe/cmd/ui/render"
)

type mr struct{}

func (r mr) PageName() string {
	return "test"
}

func (r mr) TmplFiles() []string {
	return []string{"testdata/tpl.html"}
}

func (r mr) FuncMap() template.FuncMap {
	return nil
}

func ExampleNew() {
	gin.SetMode(gin.TestMode)
	gs := gin.New()
	rs := render.New()
	rs.Add(&mr{})
	gs.HTMLRender = rs.HTMLRender()
	// Output:
}
