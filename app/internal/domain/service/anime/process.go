package anime

import "anilibrary-request-parser/app/internal/domain/entity"

func (s *ScraperService) Process() (*entity.Anime, error) {
	anime, err := s.scraper.Process()

	if err != nil {
		s.logger.Error(err.Error())
		return anime, err
	}

	return anime, nil
}
