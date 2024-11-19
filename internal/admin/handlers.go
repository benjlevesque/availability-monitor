package admin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Map Server Handlers
func (s *AdminApi) MapHandlers(e *echo.Echo, prefix string) error {
	admin := e.Group(prefix)
	admin.Use(middleware.BasicAuth(func(username, password string, ctx echo.Context) (bool, error) {
		return s.credentials.Validate(username, password)
	}))

	admin.GET("/sensors/", func(c echo.Context) error {
		type sensorItem struct {
			MAC       string    `json:"mac"`
			CreatedAt time.Time `json:"createdAt"`
			Location  string    `json:"location"`
		}
		sensors := s.sensorManager.List()

		response := []sensorItem{}
		for _, sensor := range sensors {
			response = append(response, sensorItem{
				MAC:       sensor.MAC(),
				CreatedAt: sensor.CreatedAt(),
				Location:  sensor.Location(),
			})
		}
		return c.Render(http.StatusOK, "admin_sensors.html", response)
	})

	admin.GET("/sensors/pending/", func(c echo.Context) error {
		type pendingSensor struct {
			MAC       string    `json:"mac"`
			CreatedAt time.Time `json:"createdAt"`
		}
		sensors := s.sensorManager.ListPending()

		response := []pendingSensor{}
		for _, sensor := range sensors {
			response = append(response, pendingSensor{
				MAC:       sensor.MAC(),
				CreatedAt: sensor.CreatedAt(),
			})
		}

		return c.Render(http.StatusOK, "admin_pending.html", response)
	})

	admin.POST("/sensors/:mac/", func(c echo.Context) error {
		mac := c.Param("mac")
		location := c.Request().FormValue("location")
		s.sensorManager.SetLocation(mac, location)
		return c.Redirect(http.StatusFound, fmt.Sprintf("/%s/sensors/pending", prefix))
	})

	admin.DELETE("/sensors/:mac/", func(c echo.Context) error {
		mac := c.Param("mac")
		err := s.sensorManager.Delete(mac)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusAccepted, echo.Map{})
	})

	return nil
}
