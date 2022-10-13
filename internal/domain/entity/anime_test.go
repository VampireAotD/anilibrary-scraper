package entity

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnime(t *testing.T) {
	t.Run("ToJSON", func(t *testing.T) {
		t.Run("With error", func(t *testing.T) {
			anime := Anime{
				Rating: float32(math.Inf(1)),
			}

			parsed, err := anime.ToJSON()
			require.Nil(t, parsed)
			require.Error(t, err)
		})

		t.Run("Without error", func(t *testing.T) {
			anime := Anime{
				Title:  "random",
				Rating: 8.9,
			}

			parsed, err := anime.ToJSON()
			require.NotNil(t, parsed)
			require.NoError(t, err)
		})
	})

	t.Run("FromJSON", func(t *testing.T) {
		var anime Anime

		t.Run("With error", func(t *testing.T) {
			data := `{"rating":8,9}`

			parsed, err := anime.FromJSON([]byte(data))
			require.NotNil(t, parsed)
			require.Error(t, err)
		})

		t.Run("Without error", func(t *testing.T) {
			data := `{"title":"random","rating":8.9}`

			parsed, err := anime.FromJSON([]byte(data))
			require.NotNil(t, parsed)
			require.NoError(t, err)
		})
	})
}
