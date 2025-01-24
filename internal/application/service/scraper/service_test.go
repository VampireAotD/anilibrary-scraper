package scraper

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/VampireAotD/anilibrary-scraper/internal/domain/entity"
	"github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/scraper"

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

func (s *ScraperServiceSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	var (
		repositoryMock = NewMockAnimeRepository(ctrl)
		scraperMock    = NewMockScraper(ctrl)
	)

	s.scraperMock = scraperMock.EXPECT()
	s.repositoryMock = repositoryMock.EXPECT()
	s.service = NewScraperService(scraperMock, repositoryMock)
}

func (s *ScraperServiceSuite) TestProcess() {
	var (
		t       = s.T()
		require = s.Require()
	)

	t.Run("Unsupported sites", func(_ *testing.T) {
		testCases := []string{
			"",
			"https://google.com",
		}

		for _, testCase := range testCases {
			s.repositoryMock.FindByURL(gomock.Any(), testCase).Return(entity.Anime{}, entity.ErrAnimeNotFound)
			s.scraperMock.ScrapeAnime(gomock.Any(), testCase).Return(entity.Anime{}, scraper.ErrSiteNotSupported)

			anime, err := s.service.Process(context.Background(), testCase)
			require.Empty(anime)
			require.ErrorIs(err, scraper.ErrSiteNotSupported)
		}
	})

	t.Run("Supported sites", func(t *testing.T) {
		t.Run("Cache hit", func(_ *testing.T) {
			const url string = "https://animego.org/anime/blich-tysyacheletnyaya-krovavaya-voyna-2129"
			anime := entity.Anime{
				Image:    base64.StdEncoding.EncodeToString([]byte("image")),
				Title:    "Блич: Тысячелетняя кровавая война",
				Status:   entity.Ready,
				Episodes: 13,
				Rating:   9.7,
			}

			s.repositoryMock.FindByURL(gomock.Any(), url).Return(anime, nil)

			result, err := s.service.Process(context.Background(), url)

			require.NoError(err)
			require.NotEmpty(result)
			require.Equal(anime, result)
		})

		t.Run("Cache miss", func(t *testing.T) {
			testCases := []struct {
				name     string
				url      string
				expected entity.Anime
			}{
				{
					name: "AnimeGo",
					url:  "https://animego.org/anime/blich-tysyacheletnyaya-krovavaya-voyna-2129",
					expected: entity.Anime{
						Image:       base64.StdEncoding.EncodeToString([]byte("data:image/png;base64,image")),
						Title:       "Наруто: Ураганные хроники",
						Status:      entity.Ready,
						Episodes:    500,
						Genres:      []string{"Боевые искусства", "Комедия", "Сёнэн"},
						VoiceActing: []string{"AniDUB", "AniLibria", "SHIZA Project", "2x2"},
						Synonyms:    []string{"Naruto: Shippuuden", "Naruto: Shippuden", "ナルト- 疾風伝"},
						Rating:      9.5,
					},
				},
				{
					name: "AnimeVost",
					url:  "https://animevost.org/tip/tv/855-akame-ga-kill-ubiyca-akame.html",
					expected: entity.Anime{
						Image:       base64.StdEncoding.EncodeToString([]byte("data:image/png;base64,image")),
						Title:       "Убийца Акаме! / Akame ga Kill!",
						Status:      entity.Ready,
						Episodes:    24,
						Genres:      []string{"приключения", "фэнтези"},
						VoiceActing: []string{"AnimeVost"},
						Synonyms:    []string{"Akame ga Kill!"},
						Rating:      10,
					},
				},
			}

			for i := range testCases {
				t.Run(testCases[i].name, func(t *testing.T) {
					t.Parallel()

					s.repositoryMock.FindByURL(gomock.Any(), testCases[i].url).Return(entity.Anime{}, entity.ErrAnimeNotFound)
					s.scraperMock.ScrapeAnime(gomock.Any(), testCases[i].url).Return(testCases[i].expected, nil)
					s.repositoryMock.Create(gomock.Any(), gomock.Any()).Return(nil)

					anime, err := s.service.Process(context.Background(), testCases[i].url)

					require.NoError(err)
					require.NotEmpty(anime)
					require.Equal(testCases[i].expected, anime)
				})
			}
		})
	})
}
