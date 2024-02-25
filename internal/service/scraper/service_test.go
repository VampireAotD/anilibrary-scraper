package scraper

import (
	"context"
	"encoding/base64"
	"testing"

	"anilibrary-scraper/internal/entity"
	"anilibrary-scraper/internal/scraper"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ScraperServiceSuite struct {
	suite.Suite

	repositoryMock *MockAnimeRepositoryMockRecorder
	scraperMock    *MockScraperMockRecorder
	service        Service
}

func TestScraperServiceSuite(t *testing.T) {
	suite.Run(t, new(ScraperServiceSuite))
}

func (ss *ScraperServiceSuite) SetupSuite() {
	ctrl := gomock.NewController(ss.T())
	defer ctrl.Finish()

	var (
		repositoryMock = NewMockAnimeRepository(ctrl)
		scraperMock    = NewMockScraper(ctrl)
	)

	ss.repositoryMock = repositoryMock.EXPECT()
	ss.scraperMock = scraperMock.EXPECT()
	ss.service = NewScraperService(repositoryMock, scraperMock)
}

func (ss *ScraperServiceSuite) TestProcess() {
	var (
		t       = ss.T()
		require = ss.Require()
	)

	t.Run("Errors", func(_ *testing.T) {
		testCases := []string{"", "https://google.com"}

		for _, testCase := range testCases {
			ss.repositoryMock.FindByURL(gomock.Any(), gomock.Any()).Return(entity.Anime{}, nil)
			ss.repositoryMock.Create(gomock.Any(), gomock.Any()).Return(nil)

			ss.scraperMock.ScrapeAnime(gomock.Any(), testCase).Return(entity.Anime{}, scraper.ErrUnsupportedScraper)

			result, err := ss.service.Process(context.Background(), testCase)

			require.Error(err)
			require.Empty(result)
		}
	})

	t.Run("Supported urls", func(t *testing.T) {
		t.Run("Retrieve from cache", func(_ *testing.T) {
			const url string = "https://animego.org/anime/blich-tysyacheletnyaya-krovavaya-voyna-2129"
			anime := entity.Anime{
				Image:    base64.StdEncoding.EncodeToString([]byte("random")),
				Title:    "Блич: Тысячелетняя кровавая война",
				Status:   "Онгоинг",
				Episodes: "1 / ?",
				Rating:   9.7,
			}

			ss.repositoryMock.FindByURL(gomock.Any(), url).Return(anime, nil)

			cached, err := ss.service.Process(context.Background(), url)

			require.NotEmpty(cached)
			require.NoError(err)
			require.Equal(anime, cached)
		})

		t.Run("Without cache", func(t *testing.T) {
			cases := []struct {
				name     string
				url      string
				expected entity.Anime
			}{
				{
					name: "AnimeGo",
					url:  "https://animego.org/anime/naruto-102",
					expected: entity.Anime{
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
					expected: entity.Anime{
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

			for i := range cases {
				t.Run(cases[i].name, func(t *testing.T) {
					t.Parallel()

					ss.repositoryMock.FindByURL(gomock.Any(), gomock.Any()).Return(entity.Anime{}, nil)
					ss.repositoryMock.Create(gomock.Any(), gomock.Any()).Return(nil)

					ss.scraperMock.ScrapeAnime(gomock.Any(), cases[i].url).Return(cases[i].expected, nil)

					anime, err := ss.service.Process(context.Background(), cases[i].url)

					require.NoError(err)
					require.NotEmpty(anime)
					require.Equal(cases[i].expected, anime)

					_, err = base64.StdEncoding.DecodeString(anime.Image)
					require.NoError(err)
				})
			}
		})
	})
}
