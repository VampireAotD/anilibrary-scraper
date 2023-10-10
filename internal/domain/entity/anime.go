package entity

type Anime struct {
	Image       string   `json:"image"`
	Title       string   `json:"title"`
	Status      string   `json:"status"`
	Episodes    string   `json:"episodes"`
	Genres      []string `json:"genres"`
	VoiceActing []string `json:"voiceActing"`
	Synonyms    []string `json:"synonyms"`
	Rating      float32  `json:"rating"`
}
