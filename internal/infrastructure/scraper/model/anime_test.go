package model

import (
	"encoding/base64"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestAnime(t *testing.T) {
	t.Run("Validation errors", func(t *testing.T) {
		t.Run("Required data", func(t *testing.T) {
			validate := validator.New()

			t.Run("Without image", func(t *testing.T) {
				anime := Anime{
					Title: "test",
					Type:  Movie,
					Year:  2020,
				}
				err := anime.Validate(validate)
				require.Error(t, err)
				require.ErrorContains(t, err, "Image")
			})

			t.Run("Without title", func(t *testing.T) {
				anime := Anime{
					Image: base64.StdEncoding.EncodeToString([]byte("data:image/gif;base64,image")),
					Type:  Show,
					Year:  2020,
				}
				err := anime.Validate(validate)
				require.Error(t, err)
				require.ErrorContains(t, err, "Title")
			})

			t.Run("Without type", func(t *testing.T) {
				anime := Anime{
					Image: base64.StdEncoding.EncodeToString([]byte("data:image/jpg;base64,image")),
					Title: "test",
					Year:  2020,
				}
				err := anime.Validate(validate)
				require.Error(t, err)
				require.ErrorContains(t, err, "Type")
			})

			t.Run("Without year", func(t *testing.T) {
				anime := Anime{
					Image: base64.StdEncoding.EncodeToString([]byte("data:image/webp;base64,image")),
					Title: "test",
					Type:  Movie,
				}
				err := anime.Validate(validate)
				require.Error(t, err)
				require.ErrorContains(t, err, "Year")
			})
		})
	})

	t.Run("Mapping", func(t *testing.T) {
		anime := Anime{
			Image: base64.StdEncoding.EncodeToString([]byte("data:image/jpg;base64,image")),
			Title: "test",
			Type:  Show,
			Year:  2020,
		}
		mapped := anime.MapToDomainEntity()
		require.Equal(t, anime.Image, mapped.Image)
		require.Equal(t, anime.Title, mapped.Title)
		require.Equal(t, anime.Type, mapped.Type)
		require.Equal(t, anime.Year, mapped.Year)
	})
}
