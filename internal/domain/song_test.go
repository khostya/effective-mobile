package domain

import (
	"github.com/khostya/effective-mobile/internal/dto"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestText_GetVerse(t *testing.T) {
	verse1 := "ля-ля\nля"
	verse2 := "тля-тл\nя-тля"
	text := Text(strings.Join([]string{verse1, verse2}, "\n\n"))

	t.Run("one line", func(t *testing.T) {
		text := Text("1")
		verses, err := text.GetVerse(dto.Page{
			Page: 1,
			Size: 10,
		})
		require.NoError(t, err)
		require.Equal(t, []string{"1"}, verses)
	})

	t.Run("error out of range", func(t *testing.T) {
		text := Text("")
		verses, err := text.GetVerse(dto.Page{
			Page: 4,
			Size: 10,
		})
		require.ErrorIs(t, err, ErrOutOfRange)
		require.Nil(t, verses)
	})

	t.Run("verse1", func(t *testing.T) {
		verses, err := text.GetVerse(dto.Page{
			Page: 1,
			Size: 1,
		})
		require.Equal(t, []string{verse1}, verses)
		require.NoError(t, err)
	})

	t.Run("verse2", func(t *testing.T) {
		verses, err := text.GetVerse(dto.Page{
			Page: 2,
			Size: 1,
		})
		require.Equal(t, []string{verse2}, verses)
		require.NoError(t, err)
	})

	t.Run("error page invalid", func(t *testing.T) {
		verses, err := text.GetVerse(dto.Page{
			Page: 0,
			Size: 1,
		})
		require.Error(t, err)
		require.Nil(t, verses)
	})
}
