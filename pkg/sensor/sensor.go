package sensor

import "time"

type Sensor interface {
	MAC() string
	Location() string
	CreatedAt() time.Time
	LastSeenAt() time.Time
	LastMotionAt() time.Time
}
