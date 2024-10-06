package models

import (
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/models/anthropic"
	"github.com/robertoseba/gennie/internal/models/openai"
)

type ModelEnum string

const (
	OpenAIMini   ModelEnum = "gpt-4o-mini"
	OpenAI       ModelEnum = "gpt-4o"
	ClaudeSonnet ModelEnum = "sonnet"
	Maritaca     ModelEnum = "maritaca"
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
		return DefaultModel.String()
	}
}

func NewModel(modelType ModelEnum, client httpclient.IHttpClient) IModel {
	switch modelType {
	case OpenAI:
		return NewBaseModel(client, openai.NewProvider(string(modelType)))
	case OpenAIMini:
		return NewBaseModel(client, openai.NewProvider(string(modelType)))
	case ClaudeSonnet:
		return NewBaseModel(client, anthropic.NewProvider(string(modelType)))
	case Maritaca:
		panic("Not implemented yet")
	default:
		return NewBaseModel(client, openai.NewProvider(string(modelType)))
	}

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
