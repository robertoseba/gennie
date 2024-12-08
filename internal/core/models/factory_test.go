package models

import (
	"testing"

	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/robertoseba/gennie/internal/core/models/anthropic"
	"github.com/robertoseba/gennie/internal/core/models/groq"
	"github.com/robertoseba/gennie/internal/core/models/maritaca"
	"github.com/robertoseba/gennie/internal/core/models/ollama"
	"github.com/robertoseba/gennie/internal/core/models/openai"
	"github.com/stretchr/testify/require"
)

func TestNewModel(t *testing.T) {
	t.Run("OpenAI", func(t *testing.T) {
		m, err := NewModel("gpt-4o", nil, *config.NewConfig())
		require.NoError(t, err)
		require.NotNil(t, m)
		require.Equal(t, "GPT-4o (OPENAI)", m.Model().String())
		require.IsType(t, &openai.OpenAIModel{}, m.modelProvider)
	})

	t.Run("OpenAIMini", func(t *testing.T) {
		m, err := NewModel("gpt-4o-mini", nil, *config.NewConfig())
		require.NoError(t, err)
		require.NotNil(t, m)
		require.Equal(t, "GPT-4o-mini (OPENAI)", m.Model().String())
		require.IsType(t, &openai.OpenAIModel{}, m.modelProvider)
	})

	t.Run("ClaudeSonnet", func(t *testing.T) {
		m, err := NewModel("sonnet", nil, *config.NewConfig())
		require.NoError(t, err)
		require.NotNil(t, m)
		require.Equal(t, "Claude Sonnet 3.5 (ANTHROPIC)", m.Model().String())
		require.IsType(t, &anthropic.AnthropicModel{}, m.modelProvider)
	})

	t.Run("Maritaca", func(t *testing.T) {
		m, err := NewModel("maritaca", nil, *config.NewConfig())
		require.NoError(t, err)
		require.NotNil(t, m)
		require.Equal(t, "Maritaca (BR)", m.Model().String())
		require.IsType(t, &maritaca.MaritacaModel{}, m.modelProvider)
	})

	t.Run("Groq", func(t *testing.T) {
		m, err := NewModel("groq", nil, *config.NewConfig())
		require.NoError(t, err)
		require.NotNil(t, m)
		require.Equal(t, "Groq (LLAMA-3.2-90B)", m.Model().String())
		require.IsType(t, &groq.GroqModel{}, m.modelProvider)
	})

	t.Run("Ollama", func(t *testing.T) {
		m, err := NewModel("ollama", nil, *config.NewConfig())
		require.NoError(t, err)
		require.NotNil(t, m)
		require.Equal(t, "Ollama", m.Model().String())
		require.IsType(t, &ollama.OllamaAIModel{}, m.modelProvider)
	})

	t.Run("Invalid", func(t *testing.T) {
		m, err := NewModel("invalid", nil, *config.NewConfig())
		require.Nil(t, m)
		require.Error(t, err)
		require.Equal(t, "model not found", err.Error())
	})
}
