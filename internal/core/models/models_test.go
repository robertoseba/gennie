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

func TestSlugAndString(t *testing.T) {
	testTable := []struct {
		model          ModelEnum
		expectedSlug   string
		expectedString string
	}{
		{OpenAI, "gpt-4o", "GPT-4o (OPENAI)"},
		{OpenAIMini, "gpt-4o-mini", "GPT-4o-mini (OPENAI)"},
		{ClaudeSonnet, "sonnet", "Claude Sonnet 3.5 (ANTHROPIC)"},
		{Maritaca, "maritaca", "Maritaca (BR)"},
		{Groq, "groq", "Groq (LLAMA-3.2-90B)"},
		{Ollama, "ollama", "Ollama"},
	}

	for _, tt := range testTable {
		t.Run(tt.expectedString, func(t *testing.T) {
			require.Equal(t, tt.expectedSlug, tt.model.Slug())
			require.Equal(t, tt.expectedString, tt.model.String())
		})
	}
}
