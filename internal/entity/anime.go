package entity

type Anime struct {
	Image       string
	Title       string
	Status      string
	Episodes    string
	Genres      []string
	VoiceActing []string
	Synonyms    []string
	Rating      float32
}

func (a Anime) Acceptable() bool {
	return a.Image != "" && a.Title != ""
}
