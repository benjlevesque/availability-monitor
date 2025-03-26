package server

import (
	"embed"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type httpServer struct {
	*echo.Echo
}

//go:embed views/*
var views embed.FS

func NewServer() *httpServer {
	e := echo.New()
	e.HideBanner = true
	e.Pre(middleware.AddTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	templates := template.Must(template.ParseFS(views, "views/*.html"))

	e.Renderer = &TemplateRenderer{
		templates: templates,
	}

	return &httpServer{
		Echo: e,
	}

}
