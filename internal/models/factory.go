package models

import (
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/models/anthropic"
	"github.com/robertoseba/gennie/internal/models/openai"
)

type ModelEnum string

const (
	OpenAIMini   ModelEnum = "gpt-4o-mini"
	OpenAI                 = "gpt-4o"
	ClaudeSonnet           = "claude-3-5-sonnet-20240620"
	Maritaca               = "maritaca"
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
		return "Maritaca (USP-BR)"
	default:
		panic("Invalid model")
	}
}

func NewModel(modelType ModelEnum, client httpclient.IHttpClient) IModel {
	model := map[ModelEnum]func(httpclient.IHttpClient) IModel{
		OpenAIMini: func(httpclient.IHttpClient) IModel {
			return openai.NewModel(client, string(modelType))
		},
		OpenAI: func(httpclient.IHttpClient) IModel {
			return openai.NewModel(client, string(modelType))
		},
		ClaudeSonnet: func(httpclient.IHttpClient) IModel {
			return anthropic.NewModel(client, string(modelType))
		},
	}

	activeModel, ok := model[modelType]

	if !ok {
		return model[DefaultModel](client)
	}

	return activeModel(client)
}

func ListModels() []ModelEnum {
	return []ModelEnum{OpenAI, OpenAIMini, ClaudeSonnet, Maritaca}
}

func ListModelsSlug() []string {
	models := ListModels()
	modelsSlug := make([]string, len(models))

	for i, model := range models {
		modelsSlug[i] = string(model)
	}

	return modelsSlug
}
