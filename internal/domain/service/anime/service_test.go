package anime_test

import (
	"testing"

	"anilibrary-scraper/internal/domain/dto"
	"anilibrary-scraper/internal/domain/repository/mocks"
	"anilibrary-scraper/internal/domain/service/anime"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ScraperServiceSuite struct {
	suite.Suite

	repositoryMock *mocks.MockAnimeRepository
	service        anime.ScraperService
}

func TestScraperServiceSuite(t *testing.T) {
	suite.Run(t, new(ScraperServiceSuite))
}

func (suite *ScraperServiceSuite) SetupSuite() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	suite.repositoryMock = mocks.NewMockAnimeRepository(ctrl)
	suite.service = anime.NewScraperService(suite.repositoryMock)
}

func (suite *ScraperServiceSuite) composeDto(testCase string) dto.RequestDTO {
	return dto.RequestDTO{
		Url:       testCase,
		FromCache: true,
	}
}

func (suite *ScraperServiceSuite) TestProcess() {
	t := suite.T()
	cases := []struct {
		name         string
		url          string
		requireError bool
	}{
		{
			name:         "Random url",
			url:          "https://google.com",
			requireError: true,
		},
		{
			name:         "AnimeGo",
			url:          "https://animego.org/anime/naruto-102",
			requireError: false,
		},
		{
			name:         "AnimeVostOrg",
			url:          "https://animevost.org/tip/tv/5-naruto-shippuuden12.html",
			requireError: false,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			suite.repositoryMock.EXPECT().
				FindByUrl(gomock.Any(), gomock.Any()).
				Return(nil, nil)

			suite.repositoryMock.EXPECT().
				Create(gomock.Any(), testCase.url, gomock.Any()).
				Return(nil)

			result, err := suite.service.Process(suite.composeDto(testCase.url))

			if testCase.requireError {
				require.Error(t, err, testCase.name)
				require.Nil(t, result)
			} else {
				require.NoError(t, err, testCase.name)
				require.NotNil(t, result)
			}
		})
	}
}
