package model

import (
	"encoding/base64"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestAnime(t *testing.T) {
	t.Run("Validate", func(t *testing.T) {
		validate := validator.New()

		t.Run("Without image", func(t *testing.T) {
			anime := Anime{Title: "random"}
			require.ErrorIs(t, anime.Validate(validate), ErrNotEnoughData)
		})

		t.Run("Without title", func(t *testing.T) {
			anime := Anime{Image: base64.StdEncoding.EncodeToString([]byte("random"))}
			require.ErrorIs(t, anime.Validate(validate), ErrNotEnoughData)
		})
	})

	t.Run("MapToDomainEntity", func(t *testing.T) {
		anime := Anime{Image: base64.StdEncoding.EncodeToString([]byte("random")), Title: "random"}
		mapped := anime.MapToDomainEntity()

		require.Equal(t, anime.Image, mapped.Image)
		require.Equal(t, anime.Title, mapped.Title)
	})
}
