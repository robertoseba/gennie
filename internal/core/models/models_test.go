package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseFrom(t *testing.T) {
	t.Run("Loads correct enum", func(t *testing.T) {
		model, ok := ParseFrom("gpt-4o-mini")

		require.True(t, ok)
		require.Equal(t, OpenAIMini, model)
	})

	t.Run("Returns false if cant find enum", func(t *testing.T) {
		model, ok := ParseFrom("gpt-4o-mini-2")

		require.False(t, ok)
		require.Equal(t, DefaultModel, model)
	})
}

func TestListModels(t *testing.T) {
	models := ListModels()
	require.Len(t, models, 6)
	require.Contains(t, models, OpenAI)
	require.Contains(t, models, OpenAIMini)
	require.Contains(t, models, ClaudeSonnet)
	require.Contains(t, models, Maritaca)
	require.Contains(t, models, Groq)
	require.Contains(t, models, Ollama)
}
