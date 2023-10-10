package model

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnime(t *testing.T) {
	t.Run("Validate", func(t *testing.T) {
		t.Run("Without image", func(t *testing.T) {
			anime := Anime{Title: "random"}
			require.ErrorIs(t, anime.Validate(), ErrNotEnoughData)
		})

		t.Run("Without title", func(t *testing.T) {
			anime := Anime{Image: base64.StdEncoding.EncodeToString([]byte("random"))}
			require.ErrorIs(t, anime.Validate(), ErrNotEnoughData)
		})
	})

	t.Run("MapToDomainEntity", func(t *testing.T) {
		anime := Anime{Image: base64.StdEncoding.EncodeToString([]byte("random")), Title: "random"}
		require.NotNil(t, anime.MapToDomainEntity())
	})
}
