package scraper

import (
	"anilibrary-request-parser/app/internal/infrastructure/client"
	"anilibrary-request-parser/app/pkg/logger"
)

type Scrapper struct {
	Url    string
	Client *client.Client
	Logger logger.Logger
}

func New(url string, client *client.Client, logger logger.Logger) *Scrapper {
	return &Scrapper{Url: url, Client: client, Logger: logger}
}
