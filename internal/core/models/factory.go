package models

import (
	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/robertoseba/gennie/internal/core/models/anthropic"
	"github.com/robertoseba/gennie/internal/core/models/groq"
	"github.com/robertoseba/gennie/internal/core/models/maritaca"
	"github.com/robertoseba/gennie/internal/core/models/ollama"
	"github.com/robertoseba/gennie/internal/core/models/openai"
)

type ModelEnum string

const (
	OpenAIMini   ModelEnum = "gpt-4o-mini"
	OpenAI       ModelEnum = "gpt-4o"
	ClaudeSonnet ModelEnum = "sonnet"
	Maritaca     ModelEnum = "maritaca"
	Groq         ModelEnum = "groq"
	Ollama       ModelEnum = "ollama"
)

const DefaultModel = OpenAIMini

func (m ModelEnum) String() string {
	switch m {
	case OpenAIMini:
		return "GPT-4o-mini (OPENAI)"
	case OpenAI:
		return "GPT-4o (OPENAI)"
	case ClaudeSonnet:
		return "Claude Sonnet 3.5 (ANTHROPIC)"
	case Maritaca:
		return "Maritaca (BR)"
	case Groq:
		return "Groq (LLAMA-3.2-90B)"
	case Ollama:
		return "Ollama"
	default:
		return DefaultModel.String()
	}
}

func ListModels() []ModelEnum {
	return []ModelEnum{OpenAI, OpenAIMini, ClaudeSonnet, Maritaca, Groq, Ollama}
}

func ListModelsSlug() []string {
	models := ListModels()
	modelsSlug := make([]string, len(models))

	for i, model := range models {
		modelsSlug[i] = string(model)
	}

	return modelsSlug
}

func NewModel(modelType ModelEnum, client IApiClient, config config.Config) IModel {
	return newBaseModel(client, providerFactory(modelType, config))
}

func providerFactory(modelType ModelEnum, config config.Config) iModelProvider {
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
		return openai.NewProvider(string(DefaultModel), config.APIKeys.OpenAiApiKey)
	}
}
