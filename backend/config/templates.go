package config

import (
	"io"

	"github.com/dannyvankooten/extemplate"
	"github.com/labstack/echo/v4"
)

func InitTemplates(e *echo.Echo) error {
	xt := extemplate.New()
	if err := xt.ParseDir("views/", []string{".html"}); err != nil {
		return err
	}
	t := &Template{
		worker: xt,
	}
	e.Renderer = t
	return nil
}

type Template struct {
	worker *extemplate.Extemplate
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	d := make(map[string]interface{}, 0)
	d["data"] = data
	return t.worker.ExecuteTemplate(w, name, d)
}
