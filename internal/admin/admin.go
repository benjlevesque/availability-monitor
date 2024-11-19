package admin

import (
	"log"

	"github.com/EtincelleCoworking/availability-monitor/pkg/sensor"
)

type SensorManager interface {
	List() []sensor.Sensor
	ListPending() []sensor.Sensor
	SetLocation(mac string, location string) error
	Delete(mac string) error
}

type Credentials interface {
	Validate(username, password string) (bool, error)
}

// AdminApi struct
type AdminApi struct {
	logger        *log.Logger
	sensorManager SensorManager
	credentials   Credentials
}

func NewAdminApi(sensorManager SensorManager, credentials Credentials, logger *log.Logger) *AdminApi {
	return &AdminApi{logger: logger, credentials: credentials, sensorManager: sensorManager}
}
