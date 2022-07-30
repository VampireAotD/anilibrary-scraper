package scraper

import (
	"anilibrary-request-parser/app/internal/infrastructure/client"
)

type Scrapper struct {
	Url    string
	Client *client.Client
}

func New(url string, client *client.Client) *Scrapper {
	return &Scrapper{Url: url, Client: client}
}
