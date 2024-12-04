package profile

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultProfile(t *testing.T) {
	t.Run("Default profile", func(t *testing.T) {
		d := DefaultProfile()

		require.Equal(t, "Default assistant", d.Name)
		require.Equal(t, "gennie", d.Author)
		require.Equal(t, "default", d.Slug)
		require.Equal(t, "You are a helpful cli assistant. Try to answer in a concise way providing the most relevant information. And examples when necessary.", d.Data)
	})
}
