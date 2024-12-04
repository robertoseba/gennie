package models

import (
	"errors"

	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/robertoseba/gennie/internal/core/models/anthropic"
	"github.com/robertoseba/gennie/internal/core/models/groq"
	"github.com/robertoseba/gennie/internal/core/models/maritaca"
	"github.com/robertoseba/gennie/internal/core/models/ollama"
	"github.com/robertoseba/gennie/internal/core/models/openai"
)

var ErrModelNotFound = errors.New("model not found")

func NewModel(modelSlug string, client IApiClient, config config.Config) (*BaseModel, error) {
	model, ok := ParseFrom(modelSlug)
	if !ok {
		return nil, ErrModelNotFound
	}

	provider := providerFactory(model, &config)

	return newBaseModel(model, client, provider), nil
}

func providerFactory(m ModelEnum, config *config.Config) iModelProvider {
	switch m {
	case OpenAI:
		return openai.NewProvider(m.Slug(), config.APIKeys.OpenAiApiKey)
	case OpenAIMini:
		return openai.NewProvider(m.Slug(), config.APIKeys.OpenAiApiKey)
	case ClaudeSonnet:
		return anthropic.NewProvider(m.Slug(), config.APIKeys.AnthropicApiKey)
	case Maritaca:
		return maritaca.NewProvider(m.Slug(), config.APIKeys.MaritacaApiKey)
	case Groq:
		return groq.NewProvider(m.Slug(), config.APIKeys.GroqApiKey)
	case Ollama:
		return ollama.NewProvider(m.Slug(), config.Ollama.Host, config.Ollama.Model)
	default:
		return nil
	}
}
