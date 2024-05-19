package model

type Event struct {
	URL       string `json:"url"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Timestamp int64  `json:"date"`
}
