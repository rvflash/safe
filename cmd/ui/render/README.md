# render

Package render provides interface to render multiple templates with Gin.
 
## Example
 
 ```go
import (
	"html/template"
	
	"github.com/gin-gonic/gin"
	"github.com/rvflash/safe/cmd/ui/router/render"
)


func toHTML(s string) template.HTML {
	return template.HTML(s)
}

type mr struct{}

func (r mr) PageName() string {
	return "test"
}

func (r mr) TmplFiles() []string {
	return []string{"templates/tpl.html"}
}

func (r mr) FuncMap() template.FuncMap {
	return template.FuncMap{
        "toHTML": toHTML,
    }
}

// ...

gs := gin.New()
rs := templates.New()
rs.Add(&mr{})
gs.HTMLRender = rs.HTMLRender()
```
