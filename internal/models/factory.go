package models

import (
	"github.com/robertoseba/gennie/internal/common"
	"github.com/robertoseba/gennie/internal/httpclient"
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

func NewModel(modelType ModelEnum, client httpclient.IHttpClient, config common.Config) IModel {
	switch modelType {
	case OpenAI:
		return NewBaseModel(client, openai.NewProvider(string(modelType), config.OpenAiApiKey))
	case OpenAIMini:
		return NewBaseModel(client, openai.NewProvider(string(modelType), config.OpenAiApiKey))
	case ClaudeSonnet:
		return NewBaseModel(client, anthropic.NewProvider(string(modelType), config.AnthropicApiKey))
	case Maritaca:
		return NewBaseModel(client, maritaca.NewProvider(string(modelType), config.MaritacaApiKey))
	case Groq:
		return NewBaseModel(client, groq.NewProvider(string(modelType), config.GroqApiKey))
	case Ollama:
		return NewBaseModel(client, ollama.NewProvider(string(modelType), config.OllamaHost, config.OllamaModel))
	default:
		return NewBaseModel(client, openai.NewProvider(string(DefaultModel), config.OpenAiApiKey))
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
