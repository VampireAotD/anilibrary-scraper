package model

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToEntity(t *testing.T) {
	testCase := &Anime{
		Image:       base64.StdEncoding.EncodeToString([]byte("random")),
		Title:       "random",
		Status:      Ongoing,
		Episodes:    "12 / 100",
		Genres:      []string{"random", "genre"},
		VoiceActing: []string{"random", "voice acting"},
		Rating:      4.5,
	}

	entity := testCase.ToEntity()
	require.NotNil(t, entity)
	require.Equal(t, testCase.Image, entity.Image)
	_, err := base64.StdEncoding.DecodeString(entity.Image)
	require.NoError(t, err)
	require.Equal(t, testCase.Title, entity.Title)
	require.Equal(t, testCase.Status, Status(entity.Status))
	require.Equal(t, testCase.Episodes, entity.Episodes)
	require.Equal(t, testCase.Genres, entity.Genres)
	require.Equal(t, testCase.VoiceActing, entity.VoiceActing)
	require.Equal(t, testCase.Rating, entity.Rating)
}
