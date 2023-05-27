package scraper_test

import (
	"context"
	"encoding/base64"
	"testing"

	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/domain/repository"
	scraperService "anilibrary-scraper/internal/domain/service/scraper"
	"anilibrary-scraper/internal/scraper"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type ScraperServiceSuite struct {
	suite.Suite

	repositoryMock *repository.MockAnimeRepositoryMockRecorder
	scraperMock    *scraper.MockContractMockRecorder
	service        scraperService.Service
}

func TestScraperServiceSuite(t *testing.T) {
	suite.Run(t, new(ScraperServiceSuite))
}

func (suite *ScraperServiceSuite) SetupSuite() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	var (
		repositoryMock = repository.NewMockAnimeRepository(ctrl)
		scraperMock    = scraper.NewMockContract(ctrl)
	)

	suite.repositoryMock = repositoryMock.EXPECT()
	suite.scraperMock = scraperMock.EXPECT()
	suite.service = scraperService.NewScraperService(repositoryMock, scraperMock)
}

func (suite *ScraperServiceSuite) TestProcess() {
	var (
		t       = suite.T()
		require = suite.Require()
	)

	t.Run("Errors", func(t *testing.T) {
		testCases := []string{"", "https://google.com"}

		for _, testCase := range testCases {
			suite.repositoryMock.FindByURL(gomock.Any(), gomock.Any()).Return(nil, nil)
			suite.repositoryMock.Create(gomock.Any(), testCase, gomock.Any()).Return(nil)

			suite.scraperMock.Scrape(gomock.Any(), testCase).Return(nil, scraper.ErrUnsupportedScraper)

			result, err := suite.service.Process(context.Background(), testCase)

			require.Error(err)
			require.Nil(result)
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

			suite.repositoryMock.FindByURL(gomock.Any(), url).Return(anime, nil)

			cached, err := suite.service.Process(context.Background(), url)

			require.NotNil(cached)
			require.NoError(err)
			require.Equal(anime, cached)
		})

		t.Run("Without cache", func(t *testing.T) {
			cases := []struct {
				name     string
				url      string
				expected *entity.Anime
			}{
				{
					name: "AnimeGo",
					url:  "https://animego.org/anime/naruto-102",
					expected: &entity.Anime{
						Image: base64.StdEncoding.EncodeToString(
							[]byte("https://animego.org/upload/anime/images/5a3ff73e8bd5b.jpg"),
						),
						Title:       "Наруто: Ураганные хроники",
						Status:      "Вышел",
						Episodes:    "500",
						Genres:      []string{"Боевые искусства", "Комедия", "Сёнэн"},
						VoiceActing: []string{"AniDUB", "AniLibria", "SHIZA Project", "2x2"},
						Synonyms:    []string{"Naruto: Shippuuden", "Naruto: Shippuden", "ナルト- 疾風伝"},
						Rating:      9.5,
					},
				},
				{
					name: "AnimeVostOrg",
					url:  "https://animevost.org/tip/tv/855-akame-ga-kill-ubiyca-akame.html",
					expected: &entity.Anime{
						Image: base64.StdEncoding.EncodeToString(
							[]byte("https://animevost.org/uploads/posts/2014-08/1409038345_1.jpg"),
						),
						Title:       "Убийца Акаме! / Akame ga Kill! ",
						Status:      "Вышел",
						Episodes:    "24",
						Genres:      []string{"приключения", "фэнтези"},
						VoiceActing: []string{"AnimeVost"},
						Synonyms:    []string{"Akame ga Kill!"},
						Rating:      10,
					},
				},
			}

			for _, testCase := range cases {
				t.Run(testCase.name, func(t *testing.T) {
					t.Parallel()

					suite.repositoryMock.FindByURL(gomock.Any(), gomock.Any()).Return(nil, nil)
					suite.repositoryMock.Create(gomock.Any(), testCase.url, gomock.Any()).Return(nil)

					suite.scraperMock.Scrape(gomock.Any(), testCase.url).Return(testCase.expected, nil)

					anime, err := suite.service.Process(context.Background(), testCase.url)

					require.NoError(err)
					require.NotNil(anime)
					require.Equal(testCase.expected, anime)

					_, err = base64.StdEncoding.DecodeString(anime.Image)
					require.NoError(err)
				})
			}
		})
	})
}
