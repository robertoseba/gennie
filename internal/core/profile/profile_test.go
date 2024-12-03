package profile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultProfile(t *testing.T) {
	t.Run("Default profile", func(t *testing.T) {
		d := DefaultProfile()

		assert.Equal(t, "Default assistant", d.Name)
		assert.Equal(t, "gennie", d.Author)
		assert.Equal(t, "default", d.Slug)
		assert.Equal(t, "You are a helpful cli assistant. Try to answer in a concise way providing the most relevant information. And examples when necessary.", d.Data)
	})
}
