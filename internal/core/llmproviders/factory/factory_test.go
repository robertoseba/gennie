package factory

import (
	"testing"

	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/robertoseba/gennie/internal/core/llmproviders/groq"
	"github.com/robertoseba/gennie/internal/core/llmproviders/maritaca"
	"github.com/robertoseba/gennie/internal/core/llmproviders/ollama"
	"github.com/robertoseba/gennie/internal/core/llmproviders/openai"
	"github.com/stretchr/testify/require"
)

func TestNewModel(t *testing.T) {
	t.Run("OpenAI", func(t *testing.T) {
		m := NewProvider("gpt-4o", nil, *config.NewConfig())
		require.NotNil(t, m)
		require.IsType(t, &openai.OpenAIModel{}, m)
	})

	t.Run("OpenAIMini", func(t *testing.T) {
		m := NewProvider("gpt-4o-mini", nil, *config.NewConfig())
		require.NotNil(t, m)
		require.IsType(t, &openai.OpenAIModel{}, m)
	})

	t.Run("ClaudeSonnet", func(t *testing.T) {
		m, err := NewProvider("sonnet", nil, *config.NewConfig())
		require.NoError(t, err)
		require.NotNil(t, m)
		require.Equal(t, "Claude Sonnet 3.5 (ANTHROPIC)", m.Model().String())
		require.IsType(t, &anthropic.AnthropicModel{}, m.modelProvider)
	})

	t.Run("Maritaca", func(t *testing.T) {
		m, err := NewProvider("maritaca", nil, *config.NewConfig())
		require.NoError(t, err)
		require.NotNil(t, m)
		require.Equal(t, "Maritaca (BR)", m.Model().String())
		require.IsType(t, &maritaca.MaritacaModel{}, m.modelProvider)
	})

	t.Run("Groq", func(t *testing.T) {
		m, err := NewProvider("groq", nil, *config.NewConfig())
		require.NoError(t, err)
		require.NotNil(t, m)
		require.Equal(t, "Groq (DeepSeek-R1-Distill-Llama-70B)", m.Model().String())
		require.IsType(t, &groq.GroqModel{}, m.modelProvider)
	})

	t.Run("Ollama", func(t *testing.T) {
		m, err := NewProvider("ollama", nil, *config.NewConfig())
		require.NoError(t, err)
		require.NotNil(t, m)
		require.Equal(t, "Ollama", m.Model().String())
		require.IsType(t, &ollama.OllamaAIModel{}, m.modelProvider)
	})

	t.Run("Invalid", func(t *testing.T) {
		m, err := NewProvider("invalid", nil, *config.NewConfig())
		require.Nil(t, m)
		require.Error(t, err)
		require.Equal(t, "model not found", err.Error())
	})
}
