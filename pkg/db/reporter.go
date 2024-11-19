package db

import (
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/EtincelleCoworking/availability-monitor/pkg/sensor"
)

type DbReporter struct {
	db     *sqlx.DB
	logger *log.Logger
}

func NewDbReporter(db *sqlx.DB, logger *log.Logger) *DbReporter {
	reporter := DbReporter{db: db, logger: logger}
	return &reporter
}

func (r *DbReporter) Report(mac string, reportType sensor.ReportType) error {
	r.logger.Printf("%s is alive", mac)
	switch reportType {
	case sensor.Alive:
		r.updateAlive(mac)
	case sensor.Motion:
		r.updateMotion(mac)
	}
	r.updateAlive(mac)
	return nil
}

func (r *DbReporter) updateAlive(mac string) error {
	r.db.MustExec("INSERT into sensors (mac,last_seen_at) values ($1, NOW()) ON CONFLICT(mac) DO UPDATE SET last_seen_at=NOW()", mac)
	return nil
}

func (r *DbReporter) updateMotion(mac string) error {
	r.db.MustExec("UPDATE sensors set last_motion_at=NOW() where mac=$1", mac)
	return nil
}
