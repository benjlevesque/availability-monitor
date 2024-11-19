package web

import (
	"log"

	"github.com/EtincelleCoworking/availability-monitor/pkg/sensor"
)

type Reporter interface {
	Report(id string, reportType sensor.ReportType) error
}

type WebApi struct {
	logger   *log.Logger
	reporter Reporter
}

func NewWebApi(reporter Reporter, logger *log.Logger) *WebApi {
	return &WebApi{reporter: reporter, logger: logger}
}
