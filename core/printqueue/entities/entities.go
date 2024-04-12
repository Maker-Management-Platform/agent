package entities

import "time"

type JobStatusUpdate struct {
	UUID     string        `json:"uuid"`
	Duration time.Duration `json:"duration"`
	StartAt  time.Time     `json:"startAt"`
	EndAt    time.Time     `json:"endAt"`
}
