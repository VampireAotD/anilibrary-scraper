package scraper_test

import (
	"context"
	"testing"

	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/domain/repository/mocks"
	"anilibrary-scraper/internal/domain/service/scraper"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ScraperServiceSuite struct {
	suite.Suite

	repositoryMock *mocks.MockAnimeRepositoryMockRecorder
	service        scraper.Service
}

func TestScraperServiceSuite(t *testing.T) {
	suite.Run(t, new(ScraperServiceSuite))
}

func (suite *ScraperServiceSuite) SetupSuite() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	repository := mocks.NewMockAnimeRepository(ctrl)

	suite.repositoryMock = repository.EXPECT()
	suite.service = scraper.NewScraperService(repository)
}

func (suite *ScraperServiceSuite) TestProcess() {
	t := suite.T()

	t.Run("Errors", func(t *testing.T) {
		testCases := []string{"", "https://google.com"}

		for _, testCase := range testCases {
			suite.repositoryMock.FindByUrl(gomock.Any(), gomock.Any()).Return(nil, nil)
			suite.repositoryMock.Create(gomock.Any(), testCase, gomock.Any()).Return(nil)

			result, err := suite.service.Process(context.Background(), testCase)

			require.Error(t, err)
			require.Nil(t, result)
		}
	})

	t.Run("Supported urls", func(t *testing.T) {
		t.Run("Retrieve from cache", func(t *testing.T) {
			const url string = "https://animego.org/anime/blich-tysyacheletnyaya-krovavaya-voyna-2129"
			anime := &entity.Anime{
				Title:    "Блич: Тысячелетняя кровавая война",
				Status:   "Онгоинг",
				Episodes: "1 / ?",
				Rating:   9.7,
			}

			suite.repositoryMock.FindByUrl(gomock.Any(), url).Return(anime, nil)

			cached, err := suite.service.Process(context.Background(), url)

			require.NotNil(t, cached)
			require.NoError(t, err)
			require.Equal(t, anime, cached)
		})

		t.Run("Without cache", func(t *testing.T) {
			cases := []struct {
				name string
				url  string
			}{
				{
					name: "AnimeGo",
					url:  "https://animego.org/anime/naruto-102",
				},
				{
					name: "AnimeVostOrg",
					url:  "https://animevost.org/tip/tv/5-naruto-shippuuden12.html",
				},
			}

			for _, testCase := range cases {
				t.Run(testCase.name, func(t *testing.T) {
					suite.repositoryMock.FindByUrl(gomock.Any(), gomock.Any()).Return(nil, nil)
					suite.repositoryMock.Create(gomock.Any(), testCase.url, gomock.Any()).Return(nil)

					result, err := suite.service.Process(context.Background(), testCase.url)

					require.NoError(t, err)
					require.NotNil(t, result)
				})
			}
		})
	})
}
