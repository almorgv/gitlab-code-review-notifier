package scheduler

import "time"

type Config struct {
	RepeatInterval uint64
	FixedTimes     []string
	TimeZone       *time.Location
	WorkdayStartAt int
	WorkdayEndAt   int
}
