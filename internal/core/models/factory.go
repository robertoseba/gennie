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

func ListModels() []ModelEnum {
	return allModels()
}

func NewModel(modelType ModelEnum, client IApiClient, config config.Config) (*BaseModel, error) {
	provider := providerFactory(modelType, &config)

	if provider == nil {
		return nil, errors.New("model not found")
	}

	return newBaseModel(modelType, client, provider), nil
}

func providerFactory(modelType ModelEnum, config *config.Config) iModelProvider {
	switch modelType {
	case OpenAI:
		return openai.NewProvider(string(modelType), config.APIKeys.OpenAiApiKey)
	case OpenAIMini:
		return openai.NewProvider(string(modelType), config.APIKeys.OpenAiApiKey)
	case ClaudeSonnet:
		return anthropic.NewProvider(string(modelType), config.APIKeys.AnthropicApiKey)
	case Maritaca:
		return maritaca.NewProvider(string(modelType), config.APIKeys.MaritacaApiKey)
	case Groq:
		return groq.NewProvider(string(modelType), config.APIKeys.GroqApiKey)
	case Ollama:
		return ollama.NewProvider(string(modelType), config.Ollama.Host, config.Ollama.Model)
	default:
		return nil
	}
}
