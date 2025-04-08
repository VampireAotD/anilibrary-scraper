package event

import "time"

type DTO struct {
	URL       string
	CreatedAt time.Time
	IP        string
	UserAgent string
}
