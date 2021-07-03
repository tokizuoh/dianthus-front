package main

import (
	"io"
	"log"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("template/index.html")),
	}
	e.Renderer = t
	h := HogeTemplate{
		Fuga: "aa",
	}
	e.GET("/", func(c echo.Context) error {
		if err := c.Render(http.StatusOK, "index.html", h); err != nil {
			log.Println(err)
			return err
		}
		return nil
	})

	e.Logger.Fatal(e.Start(":8000"))
}

type HogeTemplate struct {
	Fuga string
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
