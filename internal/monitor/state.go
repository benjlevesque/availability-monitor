package monitor

import (
	"time"

	"github.com/EtincelleCoworking/availability-monitor/pkg/sensor"
)

func (m *MonitorService) getState(s sensor.Sensor) string {
	now := time.Now().UTC()

	if now.After(s.LastSeenAt().Add(m.unreachableTimeout)) {
		return "unreachable"
	}
	if now.After(s.LastMotionAt().Add(m.freeTimeout)) {
		return "free"
	}
	if now.After(s.LastMotionAt().Add(m.maybeFreeTimeout)) {
		return "maybe"
	}
	return "busy"
}
