package service

import (
	"context"

	"anilibrary-scraper/internal/domain/dto"
	"anilibrary-scraper/internal/domain/entity"
)

//go:generate mockgen -source=contract.go -destination=./mocks/service_mocks.go -package=mocks

type (
	ScraperService interface {
		Process(ctx context.Context, dto dto.RequestDTO) (*entity.Anime, error)
	}
)
