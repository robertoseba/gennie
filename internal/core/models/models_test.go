package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFrom(t *testing.T) {
	t.Run("Loads correct enum", func(t *testing.T) {
		model, ok := ParseFrom("gpt-4o-mini")

		assert.True(t, ok)
		assert.Equal(t, OpenAIMini, model)
	})

	t.Run("Returns false if cant find enum", func(t *testing.T) {
		model, ok := ParseFrom("gpt-4o-mini-2")

		assert.False(t, ok)
		assert.Equal(t, DefaultModel, model)
	})
}

func TestListModels(t *testing.T) {
	models := ListModels()
	assert.Equal(t, 6, len(models))
	assert.Contains(t, models, OpenAI)
	assert.Contains(t, models, OpenAIMini)
	assert.Contains(t, models, ClaudeSonnet)
	assert.Contains(t, models, Maritaca)
	assert.Contains(t, models, Groq)
	assert.Contains(t, models, Ollama)
}
