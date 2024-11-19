package web

import (
	"net/http"

	"github.com/EtincelleCoworking/availability-monitor/pkg/sensor"
	"github.com/labstack/echo/v4"
)

// Map Server Handlers
func (s *WebApi) MapHandlers(e *echo.Echo, prefix string) error {
	api := e.Group(prefix)
	api.POST("/alive/", func(c echo.Context) error {
		mac := c.QueryParam("mac")
		if mac == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Missing MAC"})
		}
		s.reporter.Report(mac, sensor.Alive)
		return c.JSON(http.StatusOK, echo.Map{})
	})
	api.POST("/motion/", func(c echo.Context) error {
		mac := c.QueryParam("mac")
		if mac == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Missing MAC"})
		}
		s.reporter.Report(mac, sensor.Motion)
		return c.JSON(http.StatusOK, echo.Map{})
	})

	return nil
}
