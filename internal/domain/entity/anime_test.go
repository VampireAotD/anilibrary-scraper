package entity

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnime(t *testing.T) {
	t.Run("FromJSON", func(t *testing.T) {
		testCase := `{"title":"random","status":"Онгоинг","episodes":"123","genres":["random"],"voiceActing":["random"],"rating":9.5}`

		var anime Anime
		decoded, err := anime.FromJSON([]byte(testCase))
		require.NoError(t, err)
		require.NotZero(t, decoded.Title)
		require.NotZero(t, decoded.Status)
		require.NotZero(t, decoded.Episodes)
		require.NotZero(t, decoded.Genres)
		require.NotZero(t, decoded.VoiceActing)
		require.NotZero(t, decoded.Rating)
	})

	t.Run("IsValid", func(t *testing.T) {
		t.Run("Empty entity", func(t *testing.T) {
			anime := new(Anime)
			require.ErrorIs(t, anime.IsValid(), ErrInvalidData)
		})

		t.Run("Partially filled", func(t *testing.T) {
			anime := &Anime{
				Title:  "random",
				Rating: 9.1,
			}
			require.ErrorIs(t, anime.IsValid(), ErrInvalidData)
		})
	})

	t.Run("ToJSON", func(t *testing.T) {
		testCase := &Anime{
			Image:       "",
			Title:       "random",
			Status:      "Онгоинг",
			Episodes:    "123",
			Genres:      []string{"random"},
			VoiceActing: []string{"random"},
			Rating:      9.5,
		}

		encoded, err := testCase.ToJSON()
		require.NoError(t, err)

		var anime *Anime
		require.NoError(t, json.Unmarshal(encoded, &anime))
		require.Equal(t, testCase, anime)
	})
}
