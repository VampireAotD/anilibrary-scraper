package entity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnime(t *testing.T) {
	t.Run("HasRequiredData", func(t *testing.T) {
		t.Run("Empty entity", func(t *testing.T) {
			anime := new(Anime)
			require.ErrorIs(t, anime.HasRequiredData(), ErrNotEnoughData)
		})

		t.Run("Partially filled", func(t *testing.T) {
			anime := &Anime{
				Title:  "random",
				Rating: 9.1,
			}
			require.ErrorIs(t, anime.HasRequiredData(), ErrNotEnoughData)
		})
	})
}
