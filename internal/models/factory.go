package models

import (
	"github.com/robertoseba/gennie/internal/apiclient"
	"github.com/robertoseba/gennie/internal/common"
	"github.com/robertoseba/gennie/internal/models/anthropic"
	"github.com/robertoseba/gennie/internal/models/groq"
	"github.com/robertoseba/gennie/internal/models/maritaca"
	"github.com/robertoseba/gennie/internal/models/ollama"
	"github.com/robertoseba/gennie/internal/models/openai"
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

func NewModel(modelType ModelEnum, client apiclient.IApiClient, config common.Config) IModel {
	return newBaseModel(client, providerFactory(modelType, config))
}

func providerFactory(modelType ModelEnum, config common.Config) iModelProvider {
	switch modelType {
	case OpenAI:
		return openai.NewProvider(string(modelType), config.OpenAiApiKey)
	case OpenAIMini:
		return openai.NewProvider(string(modelType), config.OpenAiApiKey)
	case ClaudeSonnet:
		return anthropic.NewProvider(string(modelType), config.AnthropicApiKey)
	case Maritaca:
		return maritaca.NewProvider(string(modelType), config.MaritacaApiKey)
	case Groq:
		return groq.NewProvider(string(modelType), config.GroqApiKey)
	case Ollama:
		return ollama.NewProvider(string(modelType), config.OllamaHost, config.OllamaModel)
	default:
		return openai.NewProvider(string(DefaultModel), config.OpenAiApiKey)
	}
}
