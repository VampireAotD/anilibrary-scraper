package enum

type Status string

const (
	Ongoing  Status = "Онгоинг"
	Announce Status = "Анонс"
	Ready    Status = "Вышел"
)
