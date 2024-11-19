package db

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/EtincelleCoworking/availability-monitor/pkg/sensor"
)

type DbSensorManager struct {
	db     *sqlx.DB
	logger *log.Logger
}

func NewDbSensorManager(db *sqlx.DB, logger *log.Logger) *DbSensorManager {
	mgr := DbSensorManager{db: db, logger: logger}
	return &mgr
}

func (l *DbSensorManager) List() []sensor.Sensor {
	dbSensors := []*Sensor{}
	l.db.Select(&dbSensors, "SELECT * FROM sensors where location IS NOT NULL")

	models := make([]sensor.Sensor, len(dbSensors))
	for i, v := range dbSensors {
		models[i] = SensorDisplay{v}
	}
	return models

}

func (l *DbSensorManager) ListPending() []sensor.Sensor {
	dbSensors := []*Sensor{}
	l.db.Select(&dbSensors, "SELECT * FROM sensors where location IS NULL")

	models := make([]sensor.Sensor, len(dbSensors))
	for i, v := range dbSensors {
		models[i] = SensorDisplay{v}
	}
	return models

}

func (l *DbSensorManager) SetLocation(mac string, name string) error {
	r, err := l.db.Exec("UPDATE sensors SET location=$2 where mac=$1", mac, name)
	if err != nil {
		return err
	}

	nb, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if nb <= 0 {
		return fmt.Errorf("Sensor not found for mac %s", mac)
	}
	return nil
}

func (l *DbSensorManager) Delete(mac string) error {
	r, err := l.db.Exec("DELETE FROM sensors where mac=$1", mac)
	if err != nil {
		return err
	}

	nb, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if nb <= 0 {
		return fmt.Errorf("Sensor not found for mac %s", mac)
	}
	return nil

}

type SensorDisplay struct {
	*Sensor
}

func (s SensorDisplay) Location() string {
	return s.Sensor.Location.String
}
func (s SensorDisplay) MAC() string {
	return s.Sensor.Mac
}
func (s SensorDisplay) CreatedAt() time.Time {
	return s.Sensor.CreatedAt
}
func (s SensorDisplay) LastMotionAt() time.Time {
	return s.Sensor.LastMotionAt.Time
}
func (s SensorDisplay) LastSeenAt() time.Time {
	return s.Sensor.LastSeenAt
}
