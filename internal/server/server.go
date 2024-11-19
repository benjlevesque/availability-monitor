package server

import (
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type httpServer struct {
	*echo.Echo
}

func NewServer() *httpServer {
	e := echo.New()
	e.HideBanner = true
	e.Pre(middleware.AddTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	templates := template.Must(template.ParseGlob("pkg/views/*.html"))

	e.Renderer = &TemplateRenderer{
		templates: templates,
	}

	return &httpServer{
		Echo: e,
	}

}
