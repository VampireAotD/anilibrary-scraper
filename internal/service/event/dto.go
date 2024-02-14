package event

import "time"

type DTO struct {
	URL       string
	Time      time.Time
	IP        string
	UserAgent string
}
