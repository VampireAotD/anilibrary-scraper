package model

type Event struct {
	URL       string `json:"url"`
	Timestamp int64  `json:"date"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}
