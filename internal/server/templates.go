package server

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func (t *TemplateRenderer) ParseTemplates() error {
	t.templates = template.New("")
	err := filepath.Walk("./pkg/views", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			_, err = t.templates.ParseFiles(path)
			if err != nil {
				log.Println(err)
			}
		}

		return err
	})

	return err
}
