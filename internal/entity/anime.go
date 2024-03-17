package entity

type Status string
type Type string

const (
	Ongoing Status = "Онгоинг"
	Ready   Status = "Вышел"

	Show  Type = "ТВ Сериал"
	Movie Type = "Фильм"
)

type Anime struct {
	Image       string
	Title       string
	Status      Status
	Type        Type
	Episodes    string
	Genres      []string
	VoiceActing []string
	Synonyms    []string
	Rating      float32
	Year        int
}

func (a Anime) Acceptable() bool {
	return a.Image != "" && a.Title != ""
}
