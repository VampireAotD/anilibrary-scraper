// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package providers

import (
	anime2 "anilibrary-request-parser/internal/controller/http/api/anime"
	redis2 "anilibrary-request-parser/internal/domain/repository/redis"
	"anilibrary-request-parser/internal/domain/service/anime"
	"anilibrary-request-parser/pkg/logger"
	"github.com/go-redis/redis/v9"
)

// Injectors from wire.go:

func WireScraperService(client *redis.Client) (*anime.ScraperService, error) {
	animeRepository := redis2.NewAnimeRepository(client)
	scraperService := anime.NewScraperService(animeRepository)
	return scraperService, nil
}

func WireAnimeController(client *redis.Client, logger2 logger.Logger) (anime2.Controller, error) {
	scraperService, err := WireScraperService(client)
	if err != nil {
		return anime2.Controller{}, err
	}
	controller := anime2.NewController(logger2, scraperService)
	return controller, nil
}