package entity

type Status string

type Anime struct {
	Image       string
	Title       string
	Status      Status
	Episodes    string
	Genres      []string
	VoiceActing []string
	Rating      float32
}
