package factory

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/robertoseba/gennie/internal/core/llmproviders/anthropic"
	"github.com/robertoseba/gennie/internal/core/llmproviders/gemini"
)

func TestParseFrom(t *testing.T) {
	t.Run("Loads correct enum", func(t *testing.T) {
		model, ok := ParseFrom("gpt-4.1-mini")

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
	require.Len(t, models, 7)
	require.Contains(t, models, OpenAI)
	require.Contains(t, models, OpenAIMini)
	require.Contains(t, models, ClaudeSonnet)
	require.Contains(t, models, Haiku)
	require.Contains(t, models, Gemini)
	require.Contains(t, models, Groq)
	require.Contains(t, models, Ollama)
}

func TestSlugAndString(t *testing.T) {
	testTable := []struct {
		model          ModelEnum
		expectedSlug   string
		expectedString string
	}{
		{OpenAI, "gpt-4.1", "GPT-4.1 (OPENAI)"},
		{OpenAIMini, "gpt-4.1-mini", "GPT-4.1-mini (OPENAI)"},
		{ClaudeSonnet, anthropic.ExportedSonnetSlug, anthropic.ExportedSonnetDescription},
		{Haiku, anthropic.ExportedHaikuSlug, anthropic.ExportedHaikuDescription},
		{Gemini, gemini.ExportedModelSlug, gemini.ExportedModelDescription},
		{Groq, "groq", "Groq (Qwen-QWQ-32B)"},
		{Ollama, "ollama", "Ollama"},
	}

	for _, tt := range testTable {
		require.Equal(t, tt.expectedSlug, tt.model.Slug())
		require.Equal(t, tt.expectedString, tt.model.String())
	}
}
