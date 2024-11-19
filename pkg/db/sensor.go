package db

import (
	"database/sql"
	"time"
)

type Sensor struct {
	Location     sql.NullString `db:"location"`
	Mac          string         `db:"mac"`
	CreatedAt    time.Time      `db:"created_at"`
	LastSeenAt   time.Time      `db:"last_seen_at"`
	LastMotionAt sql.NullTime   `db:"last_motion_at"`
}
