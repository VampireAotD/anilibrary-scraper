package providers

import (
	"anilibrary-scraper/internal/handler/http/v1/anime"
	"github.com/google/wire"
)

// Handlers

var HTTPAnimeHandlerSet = wire.NewSet(
	redisAnimeRepositoryBinding,
	scraperBinding,
	scraperServiceBinding,
	anime.NewController,
)
